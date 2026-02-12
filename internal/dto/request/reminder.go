package request

import "github.com/dev-hyunsang/home-library-backend/internal/domain"

// CreateReminderRequest is the request body for creating a reading reminder
type CreateReminderRequest struct {
	ReminderTime string           `json:"reminder_time"`
	DayOfWeek    domain.DayOfWeek `json:"day_of_week"`
	Message      string           `json:"message"`
}

// UpdateReminderRequest is the request body for updating a reading reminder
type UpdateReminderRequest struct {
	ReminderTime string           `json:"reminder_time"`
	DayOfWeek    domain.DayOfWeek `json:"day_of_week"`
	Message      string           `json:"message"`
}
