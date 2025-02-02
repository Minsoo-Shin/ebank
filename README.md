# 프로젝트
입금, 인출하는 GRPC 서버를 구현합니다. 유저마다 계좌는 여러 계좌를 가질 수 있습니다. 그리고 영속성은 파일로 관리합니다.

| 디렉토리            | 설명 |
|---------------------|-------------------------------------------------------------------------------------------------------|
| api/                 | API 프로토콜 버퍼 및 gRPC 서비스 정의를 포함합니다. |
| build/               | 빌드를 위한 Dockerfile 저장합니다. |
| cmd/                | 메인 애플리케이션 엔트리 포인트를 포함합니다. 여러 실행 파일이 필요한 경우 여기에 추가할 수 있습니다. |
| internal/           | 프로젝트 내부에서만 사용하는 코드를 포함합니다. |
| pkg/                 | 외부로 공개하는 라이브러리 코드를 포함한다. |
| services/           | 서비스 로직을 포함합니다. |
| services/domain/model | 서비스 로직을 포함합니다. |
| services/domain/repository | 계좌 서비스 로직을 포함합니다. |
| services/domain/service | 계좌 서비스 로직을 포함합니다. |
| data/               | JSON 파일 등의 데이터 파일을 저장합니다.    |

- log, recover interceptor 추가

# 실행 방법
`make run`

# 설명
## 도메인 구분

### User(Auth)
- 유저 생성
- 유저 업데이트
- 유저 조회
- 유저 삭제
- 로그인

### Account
- 계좌 생성
- 계좌 조회
- 계좌 업데이트
- 계좌 삭제

### Transaction
- 계좌 입급
- 계좌 인출
- 계좌 입출금 내역 조회

## api 구현

