package scheduler

import (
	"context"
	"time"

	"github.com/dev-hyunsang/my-own-library-backend/internal/domain"
	"github.com/dev-hyunsang/my-own-library-backend/internal/infrastructure/fcm"
	"github.com/dev-hyunsang/my-own-library-backend/logger"
	"github.com/go-co-op/gocron/v2"
)

type ReminderScheduler struct {
	scheduler    gocron.Scheduler
	reminderRepo domain.ReadingReminderRepository
	userRepo     domain.UserRepository
	fcmService   *fcm.FCMService
}

func NewReminderScheduler(reminderRepo domain.ReadingReminderRepository, userRepo domain.UserRepository, fcmService *fcm.FCMService) (*ReminderScheduler, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	return &ReminderScheduler{
		scheduler:    s,
		reminderRepo: reminderRepo,
		userRepo:     userRepo,
		fcmService:   fcmService,
	}, nil
}

func (rs *ReminderScheduler) Start() error {
	// 개인 리마인더 체크 (매분)
	_, err := rs.scheduler.NewJob(
		gocron.CronJob("* * * * *", false),
		gocron.NewTask(rs.checkReminders),
	)
	if err != nil {
		return err
	}

	// 매일 10시 고정 알림 (KST 기준)
	_, err = rs.scheduler.NewJob(
		gocron.CronJob("0 10 * * *", false),
		gocron.NewTask(rs.sendDailyReadingReminder),
	)
	if err != nil {
		return err
	}

	// 매일 20시 고정 알림 (KST 기준)
	_, err = rs.scheduler.NewJob(
		gocron.CronJob("0 20 * * *", false),
		gocron.NewTask(rs.sendDailyReadingReminder),
	)
	if err != nil {
		return err
	}

	rs.scheduler.Start()
	logger.Sugar().Info("Reading reminder scheduler started (with daily 10:00, 20:00 notifications)")
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
			logger.Sugar().Warnf("Failed to load timezone %s: %v", tz, err)
			continue
		}

		localTime := time.Now().In(loc)
		rs.processRemindersForTimezone(ctx, localTime, tz)
	}
}

func (rs *ReminderScheduler) processRemindersForTimezone(ctx context.Context, localTime time.Time, timezone string) {
	reminders, err := rs.reminderRepo.GetDueReminders(localTime)
	if err != nil {
		logger.Sugar().Errorf("Failed to get due reminders: %v", err)
		return
	}

	for _, rw := range reminders {
		if rw.Timezone != timezone {
			continue
		}

		if rw.FCMToken == "" {
			logger.Sugar().Warnf("User %s has no FCM token, skipping reminder", rw.UserID.String())
			continue
		}

		if rs.fcmService == nil {
			logger.Sugar().Warn("FCM service is not initialized, skipping reminder")
			continue
		}

		err := rs.fcmService.SendPush(ctx, rw.FCMToken, "Reading Reminder", rw.Reminder.Message)
		if err != nil {
			logger.Sugar().Errorf("Failed to send push notification for user %s: %v", rw.UserID.String(), err)
			continue
		}

		logger.Sugar().Infof("Sent reading reminder to user %s: %s", rw.UserID.String(), rw.Reminder.Message)
	}
}

func (rs *ReminderScheduler) sendDailyReadingReminder() {
	ctx := context.Background()

	if rs.fcmService == nil {
		logger.Sugar().Warn("FCM service is not initialized, skipping daily reading reminder")
		return
	}

	if rs.userRepo == nil {
		logger.Sugar().Warn("User repository is not initialized, skipping daily reading reminder")
		return
	}

	users, err := rs.userRepo.GetAllUsersWithFCM()
	if err != nil {
		logger.Sugar().Errorf("Failed to get users with FCM token: %v", err)
		return
	}

	title := "나만의 서재"
	body := "오늘의 독서는 하셨나요? 저희랑 함께 독서 해요!"

	successCount := 0
	for _, user := range users {
		if user.FCMToken == "" {
			continue
		}

		err := rs.fcmService.SendPush(ctx, user.FCMToken, title, body)
		if err != nil {
			logger.Sugar().Errorf("Failed to send daily reminder to user %s: %v", user.ID.String(), err)
			continue
		}
		successCount++
	}

	logger.Sugar().Infof("Daily reading reminder sent to %d/%d users", successCount, len(users))
}
