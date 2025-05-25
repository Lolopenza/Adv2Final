package domain

import (
	"time"

	"github.com/google/uuid"
)

// BusinessSlot represents a time slot that can be booked
type BusinessSlot struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	BusinessID uuid.UUID `json:"business_id" gorm:"type:uuid;not null"`
	Date       time.Time `json:"date" gorm:"type:date;not null"`
	Time       time.Time `json:"time" gorm:"type:time;not null"`
	IsBooked   bool      `json:"is_booked" gorm:"default:false"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// Appointment represents a booked slot
type Appointment struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	SlotID    uuid.UUID `json:"slot_id" gorm:"type:uuid;not null;foreignKey:ID"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// CreateSlotRequest represents the request to create multiple slots
type CreateSlotRequest struct {
	BusinessID uuid.UUID `json:"business_id" binding:"required"`
	Date       string    `json:"date" binding:"required"`
	Slots      []string  `json:"slots" binding:"required"`
}

// BookAppointmentRequest represents the request to book a slot
type BookAppointmentRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
	SlotID uuid.UUID `json:"slot_id" binding:"required"`
} 