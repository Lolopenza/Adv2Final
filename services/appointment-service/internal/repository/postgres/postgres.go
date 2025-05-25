package postgres

import (
	"context"
	"time"

	"appointment-service/internal/domain"
	"appointment-service/internal/repository"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgresRepository struct {
	db *gorm.DB
}

// NewPostgresRepository creates a new PostgreSQL repository instance
func NewPostgresRepository(dsn string) (repository.Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate the schema
	if err := db.AutoMigrate(&domain.BusinessSlot{}, &domain.Appointment{}); err != nil {
		return nil, err
	}

	return &postgresRepository{db: db}, nil
}

func (r *postgresRepository) CreateSlots(ctx context.Context, slots []domain.BusinessSlot) error {
	return r.db.WithContext(ctx).Create(&slots).Error
}

func (r *postgresRepository) GetAvailableSlots(ctx context.Context, businessID uuid.UUID, date time.Time) ([]domain.BusinessSlot, error) {
	var slots []domain.BusinessSlot
	err := r.db.WithContext(ctx).
		Where("business_id = ? AND date = ? AND is_booked = ?", businessID, date, false).
		Find(&slots).Error
	return slots, err
}

func (r *postgresRepository) UpdateSlotStatus(ctx context.Context, slotID uuid.UUID, isBooked bool) error {
	return r.db.WithContext(ctx).
		Model(&domain.BusinessSlot{}).
		Where("id = ?", slotID).
		Update("is_booked", isBooked).Error
}

func (r *postgresRepository) CreateAppointment(ctx context.Context, appointment domain.Appointment) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create appointment
		if err := tx.Create(&appointment).Error; err != nil {
			return err
		}

		// Update slot status
		if err := tx.Model(&domain.BusinessSlot{}).
			Where("id = ?", appointment.SlotID).
			Update("is_booked", true).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *postgresRepository) GetAppointmentByID(ctx context.Context, id uuid.UUID) (*domain.Appointment, error) {
	var appointment domain.Appointment
	err := r.db.WithContext(ctx).First(&appointment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &appointment, nil
}

func (r *postgresRepository) GetAppointmentsByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Appointment, error) {
	var appointments []domain.Appointment
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&appointments).Error
	return appointments, err
}
