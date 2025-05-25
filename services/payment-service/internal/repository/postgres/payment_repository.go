package postgres

import (
	"context"
	"database/sql"
	"log"

	"payment-service/internal/domain"
)

type paymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) domain.PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(payment *domain.Payment) error {
	query := `
		INSERT INTO payments (id, amount, currency, status, customer_email, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(query,
		payment.ID,
		payment.Amount,
		payment.Currency,
		payment.Status,
		payment.CustomerEmail,
		payment.Description,
		payment.CreatedAt,
		payment.UpdatedAt,
	)
	return err
}

func (r *paymentRepository) GetByID(id string) (*domain.Payment, error) {
	query := `
		SELECT id, amount, currency, status, customer_email, description, created_at, updated_at
		FROM payments
		WHERE id = $1
	`
	payment := &domain.Payment{}
	err := r.db.QueryRow(query, id).Scan(
		&payment.ID,
		&payment.Amount,
		&payment.Currency,
		&payment.Status,
		&payment.CustomerEmail,
		&payment.Description,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (r *paymentRepository) Update(payment *domain.Payment) error {
	log.Printf("Repository Update method called for payment ID: %s", payment.ID)

	query := `
		UPDATE payments
		SET amount = $1, currency = $2, status = $3, customer_email = $4, description = $5, updated_at = $6
		WHERE id = $7
	`
	log.Printf("Update query: %s", query)
	log.Printf("Parameters: amount=%.2f, currency=%s, status=%s, email=%s, description=%s, updatedAt=%v, id=%s",
		payment.Amount, payment.Currency, payment.Status, payment.CustomerEmail, payment.Description, payment.UpdatedAt, payment.ID)

	result, err := r.db.Exec(query,
		payment.Amount,
		payment.Currency,
		payment.Status,
		payment.CustomerEmail,
		payment.Description,
		payment.UpdatedAt,
		payment.ID,
	)

	if err != nil {
		log.Printf("Error executing update query: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
	} else {
		log.Printf("Rows affected by update: %d", rowsAffected)
		if rowsAffected == 0 {
			log.Printf("Warning: No rows were updated. Payment ID may not exist: %s", payment.ID)
		}
	}

	return nil
}

func (r *paymentRepository) Delete(id string) error {
	query := `DELETE FROM payments WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *paymentRepository) List(customerEmail string, page, limit int32) ([]*domain.Payment, int32, error) {
	log.Printf("Repository List method called with: customerEmail=%q, page=%d, limit=%d", customerEmail, page, limit)

	// Calculate offset
	offset := (page - 1) * limit
	log.Printf("Calculated offset: %d", offset)

	// Build the query based on whether customerEmail is provided
	var countQuery, listQuery string
	var args []interface{}
	var countArgs []interface{}

	// Base queries
	countBase := "SELECT COUNT(*) FROM payments"
	listBase := `
		SELECT id, amount, currency, status, customer_email, description, created_at, updated_at
		FROM payments
	`

	// Check if we need to filter by customer email
	if customerEmail != "" {
		countQuery = countBase + " WHERE customer_email = $1"
		listQuery = listBase + " WHERE customer_email = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3"
		args = []interface{}{customerEmail, limit, offset}
		countArgs = []interface{}{customerEmail}
		log.Printf("Using customer email filter: %s", customerEmail)
	} else {
		countQuery = countBase
		listQuery = listBase + " ORDER BY created_at DESC LIMIT $1 OFFSET $2"
		args = []interface{}{limit, offset}
		countArgs = []interface{}{} // Empty slice for no args
		log.Printf("No customer email filter")
	}

	log.Printf("Count query: %s", countQuery)
	log.Printf("List query: %s", listQuery)

	// First get the total count
	var total int32
	var err error
	if customerEmail != "" {
		// Only use args with the count query if filtering by email
		err = r.db.QueryRow(countQuery, countArgs...).Scan(&total)
	} else {
		// No args needed for a simple COUNT query
		err = r.db.QueryRow(countQuery).Scan(&total)
	}

	if err != nil {
		log.Printf("Error executing count query: %v", err)
		return nil, 0, err
	}
	log.Printf("Total count: %d", total)

	// Now get the actual payments
	rows, err := r.db.Query(listQuery, args...)
	if err != nil {
		log.Printf("Error executing list query: %v", err)
		return nil, 0, err
	}
	defer rows.Close()

	var payments []*domain.Payment
	for rows.Next() {
		payment := &domain.Payment{}
		err := rows.Scan(
			&payment.ID,
			&payment.Amount,
			&payment.Currency,
			&payment.Status,
			&payment.CustomerEmail,
			&payment.Description,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, 0, err
		}
		payments = append(payments, payment)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error from rows: %v", err)
		return nil, 0, err
	}

	log.Printf("Successfully retrieved %d payments", len(payments))
	return payments, total, nil
}

func (r *paymentRepository) WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
