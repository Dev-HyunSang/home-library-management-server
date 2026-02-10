# API Docs

## Users

### POST `/api/users/signin`

#### Request

```json
{
  "email": "me@hyunsang.dev",
  "password": "1q2w3e4r!"
}
```

#### Response

```json
{
  "user": {
    "id": "dcb05d32-79c7-11f0-ad24-acde48001122",
    "nick_name": "hyunsang",
    "email": "me@hyunsang.dev",
    "is_published": false,
    "is_terms_agreed": true,
    "timezone": "Asia/Seoul",
    "created_at": "2025-08-15T11:06:16Z",
    "updated_at": "2025-08-15T11:06:16Z"
  },
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 3600
}
```

### POST `/api/users/signup`

#### Request

```json
{
  "nick_name": "example",
  "email": "example@example.com",
  "password": "example",
  "is_published": true,
  "is_terms_agreed": true
}
```

#### Response

```json
{
  "id": "803a215c-80e1-11f0-b1ad-acde48001122",
  "nick_name": "example",
  "email": "example@example.com",
  "is_published": true,
  "is_terms_agreed": true,
  "timezone": "Asia/Seoul",
  "created_at": "2025-08-24T20:57:25.636597+09:00",
  "updated_at": "2025-08-24T20:57:25.636597+09:00"
}
```

### GET `/api/users/check/nickname`

- 닉네임 사용 가능 여부 확인
- 인증 불필요

#### Request

Query parameter: `nickname` (소문자 영문 `a-z`, 숫자 `0-9`, `_`, `.` 허용)

```
GET /api/users/check/nickname?nickname=dev_hyunsang_0625
```

#### Response (사용 가능)

```json
{
  "is_success": true,
  "message": "사용 가능한 닉네임입니다."
}
```

#### Response (사용 불가 - 이미 사용 중)

HTTP 409 Conflict

```json
{
  "is_success": false,
  "message": "이미 사용 중인 닉네임입니다."
}
```

#### Response (사용 불가 - 잘못된 형식)

HTTP 400 Bad Request

```json
{
  "error_code": "INVALID_NICKNAME",
  "error_message": "닉네임 형식이 올바르지 않습니다."
}
```

### GET `/api/users/verify/email/:email`

- 이메일 인증번호 발송
- 인증 불필요
- 인증번호 유효시간: 5분

#### Request

```
GET /api/users/verify/email/example@example.com
```

#### Response (성공)

```json
{
  "is_success": true,
  "message": "해당 이메일로 인증번호를 발송했습니다."
}
```

#### Response (이미 가입된 이메일)

HTTP 409 Conflict

```json
{
  "is_success": false,
  "message": "동일한 메일 주소가 이미 사용중입니다."
}
```

### POST `/api/users/verify/code`

- 이메일 인증번호 확인
- 인증 불필요

#### Request

```json
{
  "email": "example@example.com",
  "code": "123456"
}
```

#### Response (성공)

```json
{
  "is_success": true,
  "message": "이메일 인증이 완료되었습니다."
}
```

#### Response (실패)

HTTP 400 Bad Request

```json
{
  "is_success": false,
  "message": "인증번호가 유효하지 않거나 만료되었습니다."
}
```

### POST `/api/users/forgot-password`

- 변경 시 변경된 비밀번호로 로그인 성공

#### Request

```json
{
  "email": "me@hyunsang.dev"
}
```

#### Response

```json
{
  "email": "me@hyunsang.dev",
  "message": "비밀번호 재설정 이메일이 발송되었습니다"
}
```

### PUT `/api/users/fcm-token`

- FCM 디바이스 토큰 업데이트 (푸시 알림용)
- Authorization: Bearer {token} 필요

#### Request

```json
{
  "fcm_token": "dGVzdC10b2tlbi1mb3ItZmNtLXB1c2gtbm90aWZpY2F0aW9u..."
}
```

#### Response

```json
{
  "message": "FCM 토큰이 성공적으로 업데이트되었습니다."
}
```

### PUT `/api/users/timezone`

- 사용자 타임존 업데이트 (알림 시간 계산용)
- Authorization: Bearer {token} 필요

#### Request

```json
{
  "timezone": "Asia/Seoul"
}
```

#### Response

```json
{
  "message": "타임존이 성공적으로 업데이트되었습니다.",
  "timezone": "Asia/Seoul"
}
```

---

## Reading Reminders

독서 리마인더 기능 - 사용자가 설정한 시간에 "책 읽을 시간이에요" 알림을 받을 수 있음

### POST `/api/reminders`

- 새 알림 생성
- Authorization: Bearer {token} 필요

#### Request

```json
{
  "reminder_time": "20:00",
  "day_of_week": "everyday",
  "message": "책 읽을 시간이에요!"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| reminder_time | string | Yes | HH:MM 형식 (예: "09:00", "20:30") |
| day_of_week | string | No | everyday, monday, tuesday, wednesday, thursday, friday, saturday, sunday (기본값: everyday) |
| message | string | No | 알림 메시지 (기본값: "책 읽을 시간이에요!") |

#### Response

```json
{
  "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "user_id": "dcb05d32-79c7-11f0-ad24-acde48001122",
  "reminder_time": "20:00",
  "day_of_week": "everyday",
  "is_enabled": true,
  "message": "책 읽을 시간이에요!",
  "created_at": "2025-01-24T10:00:00Z",
  "updated_at": "2025-01-24T10:00:00Z"
}
```

### GET `/api/reminders`

- 내 알림 목록 조회
- Authorization: Bearer {token} 필요

#### Response

```json
{
  "reminders": [
    {
      "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "user_id": "dcb05d32-79c7-11f0-ad24-acde48001122",
      "reminder_time": "20:00",
      "day_of_week": "everyday",
      "is_enabled": true,
      "message": "책 읽을 시간이에요!",
      "created_at": "2025-01-24T10:00:00Z",
      "updated_at": "2025-01-24T10:00:00Z"
    },
    {
      "id": "b2c3d4e5-f6a7-8901-bcde-f23456789012",
      "user_id": "dcb05d32-79c7-11f0-ad24-acde48001122",
      "reminder_time": "08:00",
      "day_of_week": "monday",
      "is_enabled": false,
      "message": "월요일 아침 독서!",
      "created_at": "2025-01-24T11:00:00Z",
      "updated_at": "2025-01-24T11:00:00Z"
    }
  ],
  "count": 2
}
```

### PUT `/api/reminders/:id`

- 알림 수정
- Authorization: Bearer {token} 필요

#### Request

```json
{
  "reminder_time": "21:00",
  "day_of_week": "friday",
  "message": "금요일 저녁 독서 시간!"
}
```

#### Response

```json
{
  "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "user_id": "dcb05d32-79c7-11f0-ad24-acde48001122",
  "reminder_time": "21:00",
  "day_of_week": "friday",
  "is_enabled": true,
  "message": "금요일 저녁 독서 시간!",
  "created_at": "2025-01-24T10:00:00Z",
  "updated_at": "2025-01-24T12:00:00Z"
}
```

### PATCH `/api/reminders/:id/toggle`

- 알림 활성화/비활성화 토글
- Authorization: Bearer {token} 필요

#### Response

```json
{
  "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "user_id": "dcb05d32-79c7-11f0-ad24-acde48001122",
  "reminder_time": "20:00",
  "day_of_week": "everyday",
  "is_enabled": false,
  "message": "책 읽을 시간이에요!",
  "created_at": "2025-01-24T10:00:00Z",
  "updated_at": "2025-01-24T13:00:00Z"
}
```

### DELETE `/api/reminders/:id`

- 알림 삭제
- Authorization: Bearer {token} 필요

#### Response

- 204 No Content

---

## Books

### GET `/api/books/get`

- Authorization: Bearer {token} 필요

#### Response

```json
{
  "data": [
    {
      "id": "66de89ea-7da9-11f0-a3cb-acde48001122",
      "owner_id": "dcb05d32-79c7-11f0-ad24-acde48001122",
      "title": "민법총칙: 민법강의 1 (민법강의 1)",
      "author": "곽윤직^김재형",
      "book_isbn": "9791130325194",
      "registered_at": "2025-08-20T09:38:18Z",
      "complated_at": "2025-08-20T09:38:18Z"
    },
    {
      "id": "96ac1bec-7da9-11f0-a3cb-acde48001122",
      "owner_id": "dcb05d32-79c7-11f0-ad24-acde48001122",
      "title": "민법총칙: 민법강의 1 (민법강의 1)",
      "author": "곽윤직^김재형",
      "book_isbn": "9791130325194",
      "registered_at": "2025-08-20T09:39:38Z",
      "complated_at": "2025-08-20T09:39:38Z"
    }
  ],
  "is_success": true,
  "responsed_at": "2025-08-24T21:00:21.726454+09:00"
}
```

### POST `/api/books/add`

- Authorization: Bearer {token} 필요

#### Request

```json
{
  "title": "결혼ㆍ여름",
  "author": "알베르 카뮈",
  "book_isbn": "9791198375308"
}
```

#### Response

```json
{
  "data": {
    "id": "8ab63926-80e2-11f0-a669-acde48001122",
    "owner_id": "dcb05d32-79c7-11f0-ad24-acde48001122",
    "title": "결혼ㆍ여름",
    "author": "알베르 카뮈",
    "book_isbn": "9791198375308",
    "registered_at": "2025-08-24T21:04:52.664916+09:00",
    "complated_at": "2025-08-24T21:04:52.664917+09:00"
  },
  "is_success": true,
  "responsed_at": "2025-08-24T21:04:52.670547+09:00"
}
```

### DELETE `/api/books/delete/:id`

- Ex) `/api/books/delete/ef6e7a96-7da8-11f0-9a1c-acde48001122`
- Authorization: Bearer {token} 필요

#### Response

- 204 No Content

### POST `/api/books/search`

- ISBN으로 책 정보 검색 (네이버 API)
- Authorization: Bearer {token} 필요

#### Request

```json
{
  "book_isbn": "9791198375308"
}
```

#### Response

```json
{
  "data": {
    "lastBuildDate": "Tue, 19 Aug 2025 17:22:16 +0900",
    "total": 1,
    "start": 1,
    "display": 1,
    "items": [
      {
        "title": "결혼ㆍ여름 (태양, 입맞춤, 압생트 향… 청년 카뮈의 찬란한 감성)",
        "link": "https://search.shopping.naver.com/book/catalog/41241642621",
        "image": "https://shopping-phinf.pstatic.net/main_4124164/41241642621.20241102071337.jpg",
        "author": "알베르 카뮈",
        "discount": "17820",
        "publisher": "녹색광선",
        "pubdate": "20230804",
        "isbn": "9791198375308",
        "description": "..."
      }
    ]
  }
}
```

---

## Reviews (ISBN 기반)

ISBN을 기반으로 책 리뷰를 작성하고 조회하는 API. 사용자당 ISBN별로 1개의 리뷰만 작성 가능.

### POST `/api/reviews/:isbn`

- 리뷰 작성
- Authorization: Bearer {token} 필요
- 사용자당 ISBN별 1개 리뷰만 작성 가능

#### Request

```
POST /api/reviews/9788960777330
```

```json
{
  "content": "정말 좋은 책입니다. 추천합니다!",
  "rating": 5,
  "is_public": true
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| content | string | Yes | 리뷰 내용 |
| rating | int | Yes | 별점 (1-5) |
| is_public | bool | No | 공개 여부 (기본값: false) |

#### Response (성공)

HTTP 201 Created

```json
{
  "is_success": true,
  "message": "리뷰가 생성되었습니다.",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "owner_id": "123e4567-e89b-12d3-a456-426614174000",
    "book_isbn": "9788960777330",
    "content": "정말 좋은 책입니다. 추천합니다!",
    "rating": 5,
    "is_public": true,
    "created_at": "2026-02-10T15:30:00Z",
    "updated_at": "2026-02-10T15:30:00Z"
  }
}
```

#### Response (중복 리뷰)

HTTP 400 Bad Request

```json
{
  "is_success": false,
  "message": "이미 해당 책에 대한 리뷰를 작성했습니다",
  "time": "2026-02-10 15:30:00"
}
```

### GET `/api/reviews/:isbn`

- 해당 ISBN의 공개 리뷰 목록 조회
- 인증 불필요

#### Request

```
GET /api/reviews/9788960777330
```

#### Response

```json
{
  "is_success": true,
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "owner_id": "123e4567-e89b-12d3-a456-426614174000",
      "owner_nickname": "dev_hyunsang",
      "book_isbn": "9788960777330",
      "content": "정말 좋은 책입니다!",
      "rating": 5,
      "is_public": true,
      "created_at": "2026-02-10T15:30:00Z",
      "updated_at": "2026-02-10T15:30:00Z"
    }
  ],
  "count": 1
}
```

### GET `/api/reviews/:isbn/:id`

- 특정 리뷰 조회
- 인증 불필요

#### Request

```
GET /api/reviews/9788960777330/550e8400-e29b-41d4-a716-446655440000
```

#### Response

```json
{
  "is_success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "owner_id": "123e4567-e89b-12d3-a456-426614174000",
    "book_isbn": "9788960777330",
    "content": "정말 좋은 책입니다!",
    "rating": 5,
    "is_public": true,
    "created_at": "2026-02-10T15:30:00Z",
    "updated_at": "2026-02-10T15:30:00Z"
  }
}
```

### PUT `/api/reviews/:isbn/:id`

- 리뷰 수정
- Authorization: Bearer {token} 필요
- 본인 리뷰만 수정 가능

#### Request

```
PUT /api/reviews/9788960777330/550e8400-e29b-41d4-a716-446655440000
```

```json
{
  "content": "수정된 리뷰 내용입니다.",
  "rating": 4,
  "is_public": false
}
```

#### Response

```json
{
  "is_success": true,
  "message": "리뷰가 수정되었습니다.",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "owner_id": "123e4567-e89b-12d3-a456-426614174000",
    "book_isbn": "9788960777330",
    "content": "수정된 리뷰 내용입니다.",
    "rating": 4,
    "is_public": false,
    "created_at": "2026-02-10T15:30:00Z",
    "updated_at": "2026-02-10T16:00:00Z"
  }
}
```

### DELETE `/api/reviews/:isbn/:id`

- 리뷰 삭제
- Authorization: Bearer {token} 필요
- 본인 리뷰만 삭제 가능

#### Request

```
DELETE /api/reviews/9788960777330/550e8400-e29b-41d4-a716-446655440000
```

#### Response

```json
{
  "is_success": true,
  "message": "리뷰가 삭제되었습니다."
}
```

### GET `/api/reviews/me`

- 내 리뷰 목록 조회
- Authorization: Bearer {token} 필요

#### Response

```json
{
  "is_success": true,
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "owner_id": "123e4567-e89b-12d3-a456-426614174000",
      "book_isbn": "9788960777330",
      "content": "정말 좋은 책입니다!",
      "rating": 5,
      "is_public": true,
      "created_at": "2026-02-10T15:30:00Z",
      "updated_at": "2026-02-10T15:30:00Z"
    }
  ],
  "count": 1
}
```

---

## Bookmarks

### POST `/api/books/bookmarks/add/:id`

- 북마크 추가
- Authorization: Bearer {token} 필요

### GET `/api/books/bookmarks/get`

- 내 북마크 목록 조회
- Authorization: Bearer {token} 필요

### DELETE `/api/books/bookmarks/delete/:id`

- 북마크 삭제
- Authorization: Bearer {token} 필요

---

## Auth

### POST `/api/auth/refresh`

- 토큰 갱신

### POST `/api/auth/revoke-all`

- 모든 토큰 무효화
- Authorization: Bearer {token} 필요

### GET `/api/auth/rate-limit`

- Rate limit 상태 확인
- Authorization: Bearer {token} 필요

---

## Admin

관리자 전용 API.

### 인증 방식

1. **Bootstrap Key** (최초 API Key 생성용): `X-Admin-Bootstrap-Key` 헤더
2. **API Key** (일반 관리자 API): `X-Admin-API-Key` 헤더 또는 `api_key` 쿼리 파라미터

### POST `/api/admin/bootstrap/api-keys`

- 최초 Admin API Key 생성 (Bootstrap Key 필요)
- X-Admin-Bootstrap-Key: {ADMIN_BOOTSTRAP_KEY} 필요

#### Request

```json
{
  "name": "Main Admin Key",
  "expires_at": "2026-01-01T00:00:00Z"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | Yes | API Key 이름/설명 |
| expires_at | string | No | 만료 시간 (ISO 8601, null이면 무제한) |

#### Response

```json
{
  "success": true,
  "message": "API Key가 생성되었습니다. 이 키는 다시 표시되지 않으니 안전하게 보관하세요.",
  "api_key": {
    "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "name": "Main Admin Key",
    "key_prefix": "Xk@9mZ#L",
    "raw_key": "Xk@9mZ#L4pQ!wE2rT6yU",
    "is_active": true,
    "expires_at": "2026-01-01T00:00:00Z",
    "created_at": "2025-01-24T10:00:00Z",
    "updated_at": "2025-01-24T10:00:00Z"
  }
}
```

### POST `/api/admin/notifications/broadcast`

- 전체 사용자에게 일괄 푸시 알림 발송
- X-Admin-API-Key: {API_KEY} 필요

#### Request

```json
{
  "title": "공지사항",
  "message": "새로운 기능이 추가되었습니다!"
}
```

#### Response

```json
{
  "success": true,
  "message": "알림 발송이 완료되었습니다.",
  "total_users": 150,
  "sent_count": 148,
  "failed_count": 2
}
```

### GET `/api/admin/api-keys`

- API Key 목록 조회
- X-Admin-API-Key: {API_KEY} 필요

#### Response

```json
{
  "success": true,
  "api_keys": [
    {
      "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "name": "Main Admin Key",
      "key_prefix": "Xk@9mZ#L",
      "is_active": true,
      "last_used_at": "2025-01-24T15:30:00Z",
      "expires_at": null,
      "created_at": "2025-01-24T10:00:00Z",
      "updated_at": "2025-01-24T15:30:00Z"
    }
  ],
  "count": 1
}
```

### POST `/api/admin/api-keys`

- 새 API Key 생성
- X-Admin-API-Key: {API_KEY} 필요

#### Request

```json
{
  "name": "Secondary Admin Key",
  "expires_at": null
}
```

#### Response

```json
{
  "success": true,
  "message": "API Key가 생성되었습니다. 이 키는 다시 표시되지 않으니 안전하게 보관하세요.",
  "api_key": {
    "id": "b2c3d4e5-f6a7-8901-bcde-f23456789012",
    "name": "Secondary Admin Key",
    "key_prefix": "7Yh&nM2@",
    "raw_key": "7Yh&nM2@kL!pWx5vC9zQ",
    "is_active": true,
    "created_at": "2025-01-24T12:00:00Z"
  }
}
```

### PATCH `/api/admin/api-keys/:id/deactivate`

- API Key 비활성화
- X-Admin-API-Key: {API_KEY} 필요

#### Response

```json
{
  "success": true,
  "message": "API Key가 비활성화되었습니다."
}
```

### DELETE `/api/admin/api-keys/:id`

- API Key 삭제
- X-Admin-API-Key: {API_KEY} 필요

#### Response

- 204 No Content
