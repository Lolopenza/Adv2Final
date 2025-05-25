package repository

import (
	"context"
	"time"

	"appointment-service/internal/domain"

	"github.com/google/uuid"
)

// Repository defines the interface for data access
type Repository interface {
	// BusinessSlot operations
	CreateSlots(ctx context.Context, slots []domain.BusinessSlot) error
	GetAvailableSlots(ctx context.Context, businessID uuid.UUID, date time.Time) ([]domain.BusinessSlot, error)
	UpdateSlotStatus(ctx context.Context, slotID uuid.UUID, isBooked bool) error

	// Appointment operations
	CreateAppointment(ctx context.Context, appointment domain.Appointment) error
	GetAppointmentByID(ctx context.Context, id uuid.UUID) (*domain.Appointment, error)
	GetAppointmentsByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Appointment, error)
} 