package domain

import (
	"time"

	"github.com/google/uuid"
)

type DayOfWeek string

const (
	DayEveryday  DayOfWeek = "everyday"
	DayMonday    DayOfWeek = "monday"
	DayTuesday   DayOfWeek = "tuesday"
	DayWednesday DayOfWeek = "wednesday"
	DayThursday  DayOfWeek = "thursday"
	DayFriday    DayOfWeek = "friday"
	DaySaturday  DayOfWeek = "saturday"
	DaySunday    DayOfWeek = "sunday"
)

type ReadingReminder struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	ReminderTime string    `json:"reminder_time"`
	DayOfWeek    DayOfWeek `json:"day_of_week"`
	IsEnabled    bool      `json:"is_enabled"`
	Message      string    `json:"message"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ReminderWithUser struct {
	Reminder *ReadingReminder
	FCMToken string
	Timezone string
	UserID   uuid.UUID
}

type CreateReminderRequest struct {
	ReminderTime string    `json:"reminder_time"`
	DayOfWeek    DayOfWeek `json:"day_of_week"`
	Message      string    `json:"message"`
}

type UpdateReminderRequest struct {
	ReminderTime string    `json:"reminder_time"`
	DayOfWeek    DayOfWeek `json:"day_of_week"`
	Message      string    `json:"message"`
}

type ReadingReminderRepository interface {
	Create(userID uuid.UUID, reminder *ReadingReminder) (*ReadingReminder, error)
	GetByID(id uuid.UUID) (*ReadingReminder, error)
	GetByUserID(userID uuid.UUID) ([]*ReadingReminder, error)
	Update(reminder *ReadingReminder) error
	Delete(id uuid.UUID) error
	GetDueReminders(currentTimeUTC time.Time) ([]*ReminderWithUser, error)
}

type ReadingReminderUseCase interface {
	CreateReminder(userID uuid.UUID, req *CreateReminderRequest) (*ReadingReminder, error)
	GetUserReminders(userID uuid.UUID) ([]*ReadingReminder, error)
	GetReminderByID(id uuid.UUID) (*ReadingReminder, error)
	UpdateReminder(id uuid.UUID, userID uuid.UUID, req *UpdateReminderRequest) (*ReadingReminder, error)
	ToggleReminder(id uuid.UUID, userID uuid.UUID) (*ReadingReminder, error)
	DeleteReminder(id uuid.UUID, userID uuid.UUID) error
}
