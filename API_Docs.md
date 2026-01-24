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

## Reviews

### POST `/api/books/reviews/`

- 책 리뷰 작성
- Authorization: Bearer {token} 필요

### GET `/api/books/reviews/get`

- 내 리뷰 목록 조회
- Authorization: Bearer {token} 필요

#### Response

```json
{
  "data": [
    {
      "id": "ea6fe0a1-1f6f-40cd-a568-b0895e76d139",
      "book_id": "db00f4e8-95f5-11f0-aa2d-420b6f780d98",
      "owner_id": "9681c23c-95de-11f0-8288-420b6f780d98",
      "content": "제목도, 저자도, 첫 문장도, 이미 유명한 '알베르 카뮈'의 <이방인>을 이제야 읽었다...",
      "rating": 4,
      "created_at": "2025-09-20T10:49:28Z",
      "updated_at": "2025-09-20T10:49:28Z"
    }
  ],
  "is_success": true,
  "responsed_at": "2025-09-20T20:39:25.965901+09:00"
}
```

### GET `/api/books/reviews/:book_id`

- 특정 책의 공개 리뷰 조회 (인증 불필요)

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
