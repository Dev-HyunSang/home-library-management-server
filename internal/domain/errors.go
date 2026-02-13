package domain

import "errors"

var (
	ErrNotFound              = errors.New("해당 정보를 찾을 수 없습니다.")
	ErrAlreadyExists         = errors.New("해당 정보는 이미 존재합니다.")
	ErrInvalidInput          = errors.New("유효하지 않은 입력입니다.")
	ErrInternal              = errors.New("내부 오류가 발생했습니다.")
	ErrUserNotLoggedIn       = errors.New("사용자가 로그인하지 않았습니다.")
	ErrInvalidCredentials    = errors.New("유효하지 않은 자격 증명입니다.")
	ErrPermissionDenied      = errors.New("권한이 거부되었습니다.")
	ErrPrivateAccount        = errors.New("개인 계정입니다.")
	ErrInvalidNickname       = errors.New("유효하지 않은 닉네임입니다.")
	ErrAlreadyNickname       = errors.New("이미 존재하는 닉네임입니다.")
	ErrInvalidCSRFToken      = errors.New("유효하지 않은 CSRF 토큰입니다.")
	ErrTooManyRequests       = errors.New("너무 많은 요청입니다. 잠시 후 다시 시도해주세요.")
	ErrTokenExpired          = errors.New("토큰이 만료되었습니다.")
	ErrInvalidToken          = errors.New("유효하지 않은 토큰입니다.")
	ErrTermsNotAgreed        = errors.New("이용약관에 동의해야 합니다.")
	ErrInvalidReminderTime   = errors.New("유효하지 않은 알림 시간입니다. HH:MM 형식이어야 합니다.")
	ErrInvalidDayOfWeek      = errors.New("유효하지 않은 요일입니다.")
	ErrInvalidTimezone       = errors.New("유효하지 않은 타임존입니다.")
	ErrReminderNotFound      = errors.New("해당 알림을 찾을 수 없습니다.")
	ErrReminderOwnerMismatch = errors.New("알림 소유자가 일치하지 않습니다.")
	ErrEmailNotVerified      = errors.New("이메일 인증이 완료되지 않았습니다.")
	ErrVerificationCodeSent  = errors.New("이미 인증 메일이 발송되었습니다. 5분 후 다시 시도해주세요.")
	ErrPasswordMismatch      = errors.New("새 비밀번호가 일치하지 않습니다.")
)
