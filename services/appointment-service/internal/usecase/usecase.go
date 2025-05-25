package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"appointment-service/internal/domain"
	"appointment-service/internal/repository"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type UseCase interface {
	CreateSlots(ctx context.Context, req domain.CreateSlotRequest) error
	GetAvailableSlots(ctx context.Context, businessID uuid.UUID, date time.Time) ([]domain.BusinessSlot, error)
	BookAppointment(ctx context.Context, req domain.BookAppointmentRequest) (*domain.Appointment, error)
}

type appointmentUseCase struct {
	repo  repository.Repository
	nats  *nats.Conn
	redis *redis.Client
}

func NewAppointmentUseCase(repo repository.Repository, natsURL string, redisURL string) (UseCase, error) {
	// Connect to NATS
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, err
	}

	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})

	return &appointmentUseCase{
		repo:  repo,
		nats:  nc,
		redis: rdb,
	}, nil
}

func (uc *appointmentUseCase) CreateSlots(ctx context.Context, req domain.CreateSlotRequest) error {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return err
	}

	var slots []domain.BusinessSlot
	for _, timeStr := range req.Slots {
		t, err := time.Parse("15:04", timeStr)
		if err != nil {
			return err
		}

		slot := domain.BusinessSlot{
			ID:         uuid.New(),
			BusinessID: req.BusinessID,
			Date:       date,
			Time:       t,
			IsBooked:   false,
		}
		slots = append(slots, slot)
	}

	if err := uc.repo.CreateSlots(ctx, slots); err != nil {
		return err
	}

	// Invalidate Redis cache
	key := fmt.Sprintf("slots:%s:%s", req.BusinessID, req.Date)
	uc.redis.Del(ctx, key)

	return nil
}

func (uc *appointmentUseCase) GetAvailableSlots(ctx context.Context, businessID uuid.UUID, date time.Time) ([]domain.BusinessSlot, error) {
	// Try to get from cache first
	key := fmt.Sprintf("slots:%s:%s", businessID, date.Format("2006-01-02"))
	cached, err := uc.redis.Get(ctx, key).Result()
	if err == nil {
		var slots []domain.BusinessSlot
		if err := json.Unmarshal([]byte(cached), &slots); err == nil {
			return slots, nil
		}
	}

	// If not in cache, get from database
	slots, err := uc.repo.GetAvailableSlots(ctx, businessID, date)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if data, err := json.Marshal(slots); err == nil {
		uc.redis.Set(ctx, key, data, 5*time.Minute)
	}

	return slots, nil
}

func (uc *appointmentUseCase) BookAppointment(ctx context.Context, req domain.BookAppointmentRequest) (*domain.Appointment, error) {
	appointment := domain.Appointment{
		ID:     uuid.New(),
		UserID: req.UserID,
		SlotID: req.SlotID,
	}

	if err := uc.repo.CreateAppointment(ctx, appointment); err != nil {
		return nil, err
	}

	// Publish event to NATS
	event := map[string]interface{}{
		"appointment_id": appointment.ID,
		"user_id":        appointment.UserID,
		"slot_id":        appointment.SlotID,
		"timestamp":      time.Now(),
	}

	if data, err := json.Marshal(event); err == nil {
		uc.nats.Publish("appointment.created", data)
	}

	// Invalidate Redis cache
	uc.redis.Del(ctx, fmt.Sprintf("slots:*"))

	return &appointment, nil
}
