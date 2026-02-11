package scheduler

import (
	"context"
	"time"

	"github.com/dev-hyunsang/home-library-backend/internal/domain"
	"github.com/dev-hyunsang/home-library-backend/internal/infrastructure/kafka"
	"github.com/dev-hyunsang/home-library-backend/logger"
	"github.com/go-co-op/gocron/v2"
)

type ReminderScheduler struct {
	scheduler     gocron.Scheduler
	reminderRepo  domain.ReadingReminderRepository
	kafkaProducer *kafka.Producer
}

func NewReminderScheduler(reminderRepo domain.ReadingReminderRepository, kafkaProducer *kafka.Producer) (*ReminderScheduler, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	return &ReminderScheduler{
		scheduler:     s,
		reminderRepo:  reminderRepo,
		kafkaProducer: kafkaProducer,
	}, nil
}

func (rs *ReminderScheduler) Start() error {
	_, err := rs.scheduler.NewJob(
		gocron.CronJob("* * * * *", false),
		gocron.NewTask(rs.checkReminders),
	)
	if err != nil {
		return err
	}

	rs.scheduler.Start()
	logger.Init().Sugar().Info("Reading reminder scheduler started")
	return nil
}

func (rs *ReminderScheduler) Stop() error {
	return rs.scheduler.Shutdown()
}

func (rs *ReminderScheduler) checkReminders() {
	ctx := context.Background()

	timezones := []string{
		"Asia/Seoul",
		"Asia/Tokyo",
		"America/New_York",
		"America/Los_Angeles",
		"Europe/London",
		"Europe/Paris",
		"UTC",
	}

	for _, tz := range timezones {
		loc, err := time.LoadLocation(tz)
		if err != nil {
			logger.Init().Sugar().Warnf("Failed to load timezone %s: %v", tz, err)
			continue
		}

		localTime := time.Now().In(loc)
		rs.processRemindersForTimezone(ctx, localTime, tz)
	}
}

func (rs *ReminderScheduler) processRemindersForTimezone(ctx context.Context, localTime time.Time, timezone string) {
	reminders, err := rs.reminderRepo.GetDueReminders(localTime)
	if err != nil {
		logger.Init().Sugar().Errorf("Failed to get due reminders: %v", err)
		return
	}

	for _, rw := range reminders {
		if rw.Timezone != timezone {
			continue
		}

		if rw.FCMToken == "" {
			logger.Init().Sugar().Warnf("User %s has no FCM token, skipping reminder", rw.UserID.String())
			continue
		}

		err := rs.kafkaProducer.ProduceNotification(
			ctx,
			rw.UserID.String(),
			"Reading Reminder",
			rw.Reminder.Message,
			"reading_reminder",
		)
		if err != nil {
			logger.Init().Sugar().Errorf("Failed to produce notification for user %s: %v", rw.UserID.String(), err)
			continue
		}

		logger.Init().Sugar().Infof("Sent reading reminder to user %s: %s", rw.UserID.String(), rw.Reminder.Message)
	}
}
