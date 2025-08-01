# Home Library Management Server

- 본 프로젝트는 집에 있는 다양한 서적들을 효율적으로 관리할 수 있도록 도와줍니다.
- 집에 있는 도서들에 대한 독후감과 평가를 할 수 있는 기능을 비롯하여, 내가 완독하였던 책들의 목록을 가져옵니다.

## 사용한 기술

- Golang
  - Session with `gofiber`
  - ORM with `entgo`
- Redis(in Docker)
- MySQL(in Docker)

## 보안성

- 본 프로젝트의 모든 쿠키들은 암호화 됩니다. / 세션 쿠키 등을 비롯하여 모든 쿠키

  - [Encrypt Cookie](https://docs.lou2.kr/go-fiber/home/api/middleware/encrypt-cookie?q=) 관련 문서

- 본 프로젝트에서의 모든 암호는 안전하게 암호화 됩니다.

- 본 프로젝트는 JWT(JSON Web Token)을 사용하지 않고, Session 인증 방식을 사용합니다.
  > Stateless JWT tokens cannot be invalidated or updated, and will introduce either size issues or security issues depending on where you store them. Stateful JWT tokens are functionally the same as session cookies, but without the battle-tested and well-reviewed implementations or client support.
  >
  > Unless you work on a Reddit-scale application, there's no reason to be using JWT tokens as a session mechanism. Just use sessions.  
  > [Stop using JWT for sessions](http://cryto.net/~joepie91/blog/2016/06/13/stop-using-jwt-for-sessions/)