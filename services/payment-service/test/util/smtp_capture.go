package util

import (
	"bytes"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

// EmailMessage represents a captured email message
type EmailMessage struct {
	From       string
	To         []string
	Subject    string
	Body       string
	RawContent []byte
	ReceivedAt time.Time
}

// SMTPCapture is a simple SMTP server that captures emails for testing
type SMTPCapture struct {
	addr     string
	listener net.Listener
	messages []*EmailMessage
	mutex    sync.Mutex
	running  bool
	stopChan chan struct{}
	wg       sync.WaitGroup
}

// NewSMTPCapture creates a new SMTP capture server on the specified port
func NewSMTPCapture(port int) (*SMTPCapture, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return nil, fmt.Errorf("failed to start SMTP capture server: %w", err)
	}

	return &SMTPCapture{
		addr:     fmt.Sprintf("127.0.0.1:%d", port),
		listener: listener,
		messages: make([]*EmailMessage, 0),
		stopChan: make(chan struct{}),
	}, nil
}

// Start begins accepting connections
func (s *SMTPCapture) Start() {
	s.running = true
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()

		for s.running {
			conn, err := s.listener.Accept()
			if err != nil {
				if s.running { // Only log if still running
					fmt.Printf("SMTP server error accepting connection: %v\n", err)
				}
				continue
			}

			s.wg.Add(1)
			go s.handleConnection(conn)
		}
	}()
}

// Stop shuts down the server
func (s *SMTPCapture) Stop() {
	s.running = false
	close(s.stopChan)
	s.listener.Close()
	s.wg.Wait()
}

// GetMessages returns all captured messages
func (s *SMTPCapture) GetMessages() []*EmailMessage {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Return a copy to avoid race conditions
	msgsCopy := make([]*EmailMessage, len(s.messages))
	copy(msgsCopy, s.messages)
	return msgsCopy
}

// ClearMessages clears all captured messages
func (s *SMTPCapture) ClearMessages() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.messages = make([]*EmailMessage, 0)
}

// GetConfig returns SMTP configuration for this server
func (s *SMTPCapture) GetConfig() (host string, port string) {
	parts := strings.Split(s.addr, ":")
	return parts[0], parts[1]
}

func (s *SMTPCapture) handleConnection(conn net.Conn) {
	defer s.wg.Done()
	defer conn.Close()

	// Send greeting
	conn.Write([]byte("220 SMTP Capture Server Ready\r\n"))

	// Buffer to collect the message
	var messageBuffer bytes.Buffer
	var fromAddr string
	var toAddrs []string

	scanner := NewCRLFScanner(conn)

	// Simple SMTP dialog
	for scanner.Scan() {
		line := scanner.Text()

		// Convert to uppercase for command matching
		ucLine := strings.ToUpper(line)

		// Handle SMTP commands
		switch {
		case strings.HasPrefix(ucLine, "HELO") || strings.HasPrefix(ucLine, "EHLO"):
			// Respond with capabilities including AUTH
			conn.Write([]byte("250-Hello\r\n"))
			conn.Write([]byte("250-SIZE 35882577\r\n"))
			conn.Write([]byte("250-8BITMIME\r\n"))
			conn.Write([]byte("250-AUTH PLAIN LOGIN\r\n"))
			conn.Write([]byte("250 OK\r\n"))

		case strings.HasPrefix(ucLine, "AUTH"):
			// Accept any authentication
			conn.Write([]byte("235 Authentication successful\r\n"))

		case strings.HasPrefix(ucLine, "MAIL FROM:"):
			fromAddr = extractEmail(line[10:])
			conn.Write([]byte("250 OK\r\n"))

		case strings.HasPrefix(ucLine, "RCPT TO:"):
			toAddrs = append(toAddrs, extractEmail(line[8:]))
			conn.Write([]byte("250 OK\r\n"))

		case ucLine == "DATA":
			conn.Write([]byte("354 Start mail input; end with <CRLF>.<CRLF>\r\n"))

			// Collect message data
			inBody := false
			var subject string

			for scanner.Scan() {
				msgLine := scanner.Text()

				// End of message
				if msgLine == "." {
					break
				}

				// Process headers
				if !inBody {
					if msgLine == "" {
						inBody = true
					} else if strings.HasPrefix(strings.ToLower(msgLine), "subject:") {
						subject = strings.TrimSpace(msgLine[8:])
					}
				}

				messageBuffer.WriteString(msgLine)
				messageBuffer.WriteString("\r\n")
			}

			// Store the message
			s.storeMessage(fromAddr, toAddrs, subject, messageBuffer.String(), messageBuffer.Bytes())
			messageBuffer.Reset()

			conn.Write([]byte("250 OK: Message accepted\r\n"))

		case ucLine == "QUIT":
			conn.Write([]byte("221 Bye\r\n"))
			return

		default:
			conn.Write([]byte("500 Command not recognized\r\n"))
		}
	}
}

func (s *SMTPCapture) storeMessage(from string, to []string, subject, body string, raw []byte) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	msg := &EmailMessage{
		From:       from,
		To:         to,
		Subject:    subject,
		Body:       body,
		RawContent: raw,
		ReceivedAt: time.Now(),
	}

	s.messages = append(s.messages, msg)
}

func extractEmail(str string) string {
	str = strings.TrimSpace(str)
	if strings.HasPrefix(str, "<") && strings.HasSuffix(str, ">") {
		str = str[1 : len(str)-1]
	}
	return str
}

// CRLFScanner helps scan SMTP messages which use CRLF line endings
type CRLFScanner struct {
	reader *bytes.Reader
	conn   net.Conn
	buffer []byte
}

// NewCRLFScanner creates a new scanner for reading SMTP messages
func NewCRLFScanner(conn net.Conn) *CRLFScanner {
	return &CRLFScanner{
		conn:   conn,
		buffer: make([]byte, 0, 1024),
	}
}

// Scan reads the next line
func (s *CRLFScanner) Scan() bool {
	buf := make([]byte, 1024)

	for {
		n, err := s.conn.Read(buf)
		if err != nil {
			return false
		}

		s.buffer = append(s.buffer, buf[:n]...)

		// Check if we have a complete line
		if idx := bytes.Index(s.buffer, []byte("\r\n")); idx >= 0 {
			s.reader = bytes.NewReader(s.buffer[:idx])
			s.buffer = s.buffer[idx+2:]
			return true
		}
	}
}

// Text returns the current line
func (s *CRLFScanner) Text() string {
	if s.reader != nil {
		bytes, _ := s.readAll(s.reader)
		return string(bytes)
	}
	return ""
}

func (s *CRLFScanner) readAll(r *bytes.Reader) ([]byte, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r)
	return buf.Bytes(), err
}

// ParseEmailContent parses an email string to extract subject and body
func ParseEmailContent(content string) (subject, body string) {
	lines := strings.Split(content, "\r\n")
	headerDone := false

	for _, line := range lines {
		if !headerDone {
			if line == "" {
				headerDone = true
				continue
			}

			if strings.HasPrefix(strings.ToLower(line), "subject:") {
				subject = strings.TrimSpace(line[8:])
			}
		} else {
			body += line + "\n"
		}
	}

	return subject, strings.TrimSpace(body)
}
