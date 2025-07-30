# API Docs

## `/users`

### POST `/`

#### Request

```json
{
  "nick_name": "testuser",
  "email": "testuser@example.com",
  "password": "testpassword"
}
```

### Response

```json
{
  "id": "31313e84-69e3-11f0-b034-acde48001122",
  "nick_name": "testuser",
  "email": "testuser@example.com",
  "password": "$2a$10$r6dyWvUeAr22CQ6j2rtpj.CKSkeJEghrL5u8SpbgQgv4GduPUCYzW",
  "created_at": "",
  "updated_at": ""
}
```

### GET '/:id'

```json
HTTP/1.1 200 OK
Date: Sat, 26 Jul 2025 07:21:36 GMT
Content-Type: application/json
Content-Length: 245
Connection: close

{
  "id": "31313e84-69e3-11f0-b034-acde48001122",
  "nick_name": "testuser",
  "email": "testuser@example.com",
  "password": "$2a$10$r6dyWvUeAr22CQ6j2rtpj.CKSkeJEghrL5u8SpbgQgv4GduPUCYzW",
  "created_at": "2025-07-26T05:41:35Z",
  "updated_at": "2025-07-26T05:41:35Z"
}
```

### Not Found - Response

```json
HTTP/1.1 404 Not Found
Date: Sat, 26 Jul 2025 07:32:12 GMT
Content-Type: application/json
Content-Length: 47
Connection: close

{
  "error": "not found user: ent: user not found"
}
```
