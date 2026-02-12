package mysql

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dev-hyunsang/my-own-library-backend/internal/domain"
	"github.com/dev-hyunsang/my-own-library-backend/lib/ent"
	"github.com/dev-hyunsang/my-own-library-backend/lib/ent/readingreminder"
	"github.com/dev-hyunsang/my-own-library-backend/lib/ent/user"
	"github.com/dev-hyunsang/my-own-library-backend/logger"
	"github.com/google/uuid"
)

type ReadingReminderRepository struct {
	client *ent.Client
}

func NewReadingReminderRepository(client *ent.Client) *ReadingReminderRepository {
	return &ReadingReminderRepository{
		client: client,
	}
}

func (r *ReadingReminderRepository) Create(userID uuid.UUID, reminder *domain.ReadingReminder) (*domain.ReadingReminder, error) {
	rr, err := r.client.ReadingReminder.Create().
		SetID(reminder.ID).
		SetReminderTime(reminder.ReminderTime).
		SetDayOfWeek(readingreminder.DayOfWeek(reminder.DayOfWeek)).
		SetIsEnabled(reminder.IsEnabled).
		SetMessage(reminder.Message).
		SetOwnerID(userID).
		Save(context.Background())

	if err != nil {
		logger.Sugar().Errorf("알림 생성 중 오류 발생: %v", err)
		return nil, fmt.Errorf("알림 생성 중 오류가 발생했습니다: %w", err)
	}

	logger.Sugar().Infof("새로운 알림을 생성했습니다. 알림ID: %s, 사용자ID: %s", rr.ID.String(), userID.String())

	return &domain.ReadingReminder{
		ID:           rr.ID,
		UserID:       userID,
		ReminderTime: rr.ReminderTime,
		DayOfWeek:    domain.DayOfWeek(rr.DayOfWeek),
		IsEnabled:    rr.IsEnabled,
		Message:      rr.Message,
		CreatedAt:    rr.CreatedAt,
		UpdatedAt:    rr.UpdatedAt,
	}, nil
}

func (r *ReadingReminderRepository) GetByID(id uuid.UUID) (*domain.ReadingReminder, error) {
	rr, err := r.client.ReadingReminder.Query().
		Where(readingreminder.ID(id)).
		WithOwner().
		Only(context.Background())

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.ErrReminderNotFound
		}
		return nil, fmt.Errorf("알림 조회 중 오류가 발생했습니다: %w", err)
	}

	var userID uuid.UUID
	if rr.Edges.Owner != nil {
		userID = rr.Edges.Owner.ID
	}

	return &domain.ReadingReminder{
		ID:           rr.ID,
		UserID:       userID,
		ReminderTime: rr.ReminderTime,
		DayOfWeek:    domain.DayOfWeek(rr.DayOfWeek),
		IsEnabled:    rr.IsEnabled,
		Message:      rr.Message,
		CreatedAt:    rr.CreatedAt,
		UpdatedAt:    rr.UpdatedAt,
	}, nil
}

func (r *ReadingReminderRepository) GetByUserID(userID uuid.UUID) ([]*domain.ReadingReminder, error) {
	reminders, err := r.client.ReadingReminder.Query().
		Where(readingreminder.HasOwnerWith(user.ID(userID))).
		Order(ent.Asc(readingreminder.FieldReminderTime)).
		All(context.Background())

	if err != nil {
		return nil, fmt.Errorf("사용자 알림 목록 조회 중 오류가 발생했습니다: %w", err)
	}

	result := make([]*domain.ReadingReminder, len(reminders))
	for i, rr := range reminders {
		result[i] = &domain.ReadingReminder{
			ID:           rr.ID,
			UserID:       userID,
			ReminderTime: rr.ReminderTime,
			DayOfWeek:    domain.DayOfWeek(rr.DayOfWeek),
			IsEnabled:    rr.IsEnabled,
			Message:      rr.Message,
			CreatedAt:    rr.CreatedAt,
			UpdatedAt:    rr.UpdatedAt,
		}
	}

	return result, nil
}

func (r *ReadingReminderRepository) Update(reminder *domain.ReadingReminder) error {
	err := r.client.ReadingReminder.UpdateOneID(reminder.ID).
		SetReminderTime(reminder.ReminderTime).
		SetDayOfWeek(readingreminder.DayOfWeek(reminder.DayOfWeek)).
		SetIsEnabled(reminder.IsEnabled).
		SetMessage(reminder.Message).
		Exec(context.Background())

	if err != nil {
		if ent.IsNotFound(err) {
			return domain.ErrReminderNotFound
		}
		return fmt.Errorf("알림 수정 중 오류가 발생했습니다: %w", err)
	}

	logger.Sugar().Infof("알림을 수정했습니다. 알림ID: %s", reminder.ID.String())
	return nil
}

func (r *ReadingReminderRepository) Delete(id uuid.UUID) error {
	err := r.client.ReadingReminder.DeleteOneID(id).Exec(context.Background())
	if err != nil {
		if ent.IsNotFound(err) {
			return domain.ErrReminderNotFound
		}
		return fmt.Errorf("알림 삭제 중 오류가 발생했습니다: %w", err)
	}

	logger.Sugar().Infof("알림을 삭제했습니다. 알림ID: %s", id.String())
	return nil
}

func (r *ReadingReminderRepository) GetDueReminders(currentTimeUTC time.Time) ([]*domain.ReminderWithUser, error) {
	currentHHMM := currentTimeUTC.Format("15:04")
	currentWeekday := strings.ToLower(currentTimeUTC.Weekday().String())

	reminders, err := r.client.ReadingReminder.Query().
		Where(
			readingreminder.IsEnabled(true),
			readingreminder.ReminderTime(currentHHMM),
			readingreminder.Or(
				readingreminder.DayOfWeekEQ(readingreminder.DayOfWeekEveryday),
				readingreminder.DayOfWeekEQ(readingreminder.DayOfWeek(currentWeekday)),
			),
		).
		WithOwner().
		All(context.Background())

	if err != nil {
		return nil, fmt.Errorf("알림 조회 중 오류가 발생했습니다: %w", err)
	}

	result := make([]*domain.ReminderWithUser, 0, len(reminders))
	for _, rr := range reminders {
		if rr.Edges.Owner == nil {
			continue
		}

		owner := rr.Edges.Owner
		if owner.FcmToken == "" {
			continue
		}

		result = append(result, &domain.ReminderWithUser{
			Reminder: &domain.ReadingReminder{
				ID:           rr.ID,
				UserID:       owner.ID,
				ReminderTime: rr.ReminderTime,
				DayOfWeek:    domain.DayOfWeek(rr.DayOfWeek),
				IsEnabled:    rr.IsEnabled,
				Message:      rr.Message,
				CreatedAt:    rr.CreatedAt,
				UpdatedAt:    rr.UpdatedAt,
			},
			FCMToken: owner.FcmToken,
			Timezone: owner.Timezone,
			UserID:   owner.ID,
		})
	}

	return result, nil
}
