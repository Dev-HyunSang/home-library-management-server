package usecase

import (
	"regexp"

	"github.com/dev-hyunsang/home-library/internal/domain"
	repository "github.com/dev-hyunsang/home-library/internal/repository/mysql"
	"github.com/google/uuid"
)

type readingReminderUseCase struct {
	reminderRepo domain.ReadingReminderRepository
}

func NewReadingReminderUseCase(reminderRepo *repository.ReadingReminderRepository) *readingReminderUseCase {
	return &readingReminderUseCase{reminderRepo: reminderRepo}
}

var timeRegex = regexp.MustCompile(`^([01]?[0-9]|2[0-3]):[0-5][0-9]$`)

func validateReminderTime(timeStr string) bool {
	return timeRegex.MatchString(timeStr)
}

func validateDayOfWeek(day domain.DayOfWeek) bool {
	validDays := map[domain.DayOfWeek]bool{
		domain.DayEveryday:  true,
		domain.DayMonday:    true,
		domain.DayTuesday:   true,
		domain.DayWednesday: true,
		domain.DayThursday:  true,
		domain.DayFriday:    true,
		domain.DaySaturday:  true,
		domain.DaySunday:    true,
	}
	return validDays[day]
}

func (uc *readingReminderUseCase) CreateReminder(userID uuid.UUID, req *domain.CreateReminderRequest) (*domain.ReadingReminder, error) {
	if !validateReminderTime(req.ReminderTime) {
		return nil, domain.ErrInvalidReminderTime
	}

	dayOfWeek := req.DayOfWeek
	if dayOfWeek == "" {
		dayOfWeek = domain.DayEveryday
	}

	if !validateDayOfWeek(dayOfWeek) {
		return nil, domain.ErrInvalidDayOfWeek
	}

	message := req.Message
	if message == "" {
		message = "책 읽을 시간이에요!"
	}

	reminder := &domain.ReadingReminder{
		ID:           uuid.New(),
		ReminderTime: req.ReminderTime,
		DayOfWeek:    dayOfWeek,
		IsEnabled:    true,
		Message:      message,
	}

	return uc.reminderRepo.Create(userID, reminder)
}

func (uc *readingReminderUseCase) GetUserReminders(userID uuid.UUID) ([]*domain.ReadingReminder, error) {
	return uc.reminderRepo.GetByUserID(userID)
}

func (uc *readingReminderUseCase) GetReminderByID(id uuid.UUID) (*domain.ReadingReminder, error) {
	return uc.reminderRepo.GetByID(id)
}

func (uc *readingReminderUseCase) UpdateReminder(id uuid.UUID, userID uuid.UUID, req *domain.UpdateReminderRequest) (*domain.ReadingReminder, error) {
	reminder, err := uc.reminderRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if reminder.UserID != userID {
		return nil, domain.ErrReminderOwnerMismatch
	}

	if req.ReminderTime != "" {
		if !validateReminderTime(req.ReminderTime) {
			return nil, domain.ErrInvalidReminderTime
		}
		reminder.ReminderTime = req.ReminderTime
	}

	if req.DayOfWeek != "" {
		if !validateDayOfWeek(req.DayOfWeek) {
			return nil, domain.ErrInvalidDayOfWeek
		}
		reminder.DayOfWeek = req.DayOfWeek
	}

	if req.Message != "" {
		reminder.Message = req.Message
	}

	err = uc.reminderRepo.Update(reminder)
	if err != nil {
		return nil, err
	}

	return reminder, nil
}

func (uc *readingReminderUseCase) ToggleReminder(id uuid.UUID, userID uuid.UUID) (*domain.ReadingReminder, error) {
	reminder, err := uc.reminderRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if reminder.UserID != userID {
		return nil, domain.ErrReminderOwnerMismatch
	}

	reminder.IsEnabled = !reminder.IsEnabled

	err = uc.reminderRepo.Update(reminder)
	if err != nil {
		return nil, err
	}

	return reminder, nil
}

func (uc *readingReminderUseCase) DeleteReminder(id uuid.UUID, userID uuid.UUID) error {
	reminder, err := uc.reminderRepo.GetByID(id)
	if err != nil {
		return err
	}

	if reminder.UserID != userID {
		return domain.ErrReminderOwnerMismatch
	}

	return uc.reminderRepo.Delete(id)
}
