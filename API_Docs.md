# API Docs

## POST `/api/users/login`

### Request

```json
{
  "email": "me@hyunsang.dev",
  "password": "1q2w3e4r!"
}
```

### Response

```json
{
  "id": "dcb05d32-79c7-11f0-ad24-acde48001122",
  "nick_name": "hyunsang",
  "email": "me@hyunsang.dev",
  "password": "$2a$10$2ChYRO3mBvMYr7.iLsk5XeMPYFRJW1q/ZV9olS.PJCW53k1VFzgk2",
  "is_published": false,
  "created_at": "2025-08-15T11:06:16Z",
  "updated_at": "2025-08-15T11:06:16Z"
}
```

## POST `/api/users/register`

### Requset

```json
{
  "nick_name": "example",
  "email": "example@example.com",
  "password": "example",
  "is_published": true
}
```

### Response

```json
{
  "id": "803a215c-80e1-11f0-b1ad-acde48001122",
  "nick_name": "example",
  "email": "example@example.com",
  "password": "$2a$10$ENbcBVNwy.rbV.rkUugr/uNsUmecD0IRAdTMmvS4K1G7Oc2JZKguG",
  "is_published": true,
  "created_at": "2025-08-24T20:57:25.636597+09:00",
  "updated_at": "2025-08-24T20:57:25.636597+09:00"
}
```

### GET `/api/books/`

### Response

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

## POST `/api/books/`

### Request

```json
{
  "title": "결혼ㆍ여름",
  "author": "알베르 카뮈",
  "book_isbn": "9791198375308"
}
```

### Response

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

## DELETE `/api/books/:id`

- Ex) `/books/ef6e7a96-7da8-11f0-9a1c-acde48001122`

### Response

- 204 No Content

## POST `/api/books/search`

### Request

```json
{
  "book_isbn": "9791198375308"
}
```

### Response

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
        "description": "영원한 청춘의 책, 알베르 카뮈의 『결혼 · 여름』이 교보문고 특별 리커버 판으로 새롭게 선보인다. 북 디자인 키워드는 ‘포레스트(Forest)’\n『결혼 · 여름』은 카뮈 사상의 핵심인 ‘부조리’와 ‘반항’의 출발 및 완성 과정이 육성으로 들리는 듯한 자전적 기록이다. 하지만 수많은 이들을 벅차 오르게 했던 『결혼 · 여름』의 가장 큰 매력은 감각적이며 관능적인 문체다. 드물게 시와 사상, 예술과 철학이 완벽하게 결합된 에세이가 우리에게 닿았다.\n\n이 에세이가 출간된 시기는 카뮈가 『이방인』으로 최고의 작가가 되기 전이다. 카뮈의 유년기부터 20대 초중반까지의 시간은 그야말로 좌절과 불확실함의 연속이었다. 학교에 다니는 것조차 사치였던 가난한 유년시절, 열일곱 살에 발병해 그를 죽음 근처로 몰아갔던 폐결핵, 스물한 살에 감행한 사랑하는 사람과의 이른 결혼과 파국, 폐결핵 병력으로 인한 교수 응시 자격의 박탈.\n\n그럼에도 불구하고 그는 다음과 같이 쓴다. 사는 것이 파멸을 향해 달려가는 것이라 해도, 이 세계 속에서 사랑과 욕망을 찾아 걸어 나가겠다고.\n\n『결혼 · 여름』의 오리지널 표지가 티파사의 바다 이미지를 표현한 것이라면, \n리커버:k 표지는 흰색과 녹색 컬러를 메인으로 작업하여 알제의 여름 이미지를 표현했다.\n\n어떤 글은 시간이 흘러도 전혀 나이를 먹지 않는다. 그는 불의의 사고로 세상을 떠났지만, 『결혼 · 여름』이 지닌 청춘의 생명력은 읽는 이로 하여금 젊음을 마주한 느낌, 다시 젊음을 되찾는 기분을 선사할 것이다."
      }
    ]
  }
}
```
