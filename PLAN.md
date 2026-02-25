# Sprint 4 Plan (2026-03-02 ~ 2026-03-06)

## Sprint Goal

백엔드와 프론트엔드에 구조화된 로깅 정책을 적용하여, 일관된 로그 포맷과 레벨 기반 로깅 인프라를 확보한다.

## Tickets

| # | Title | Type | Priority | Points | Component |
|---|-------|------|----------|--------|-----------|
| #21 | Backend 구조화된 로깅 정책 적용 (log/slog) | story | medium | 3pts | backend |
| #22 | Frontend 로깅 유틸리티 및 에러 리포팅 적용 | story | medium | 2pts | frontend |

## Sprint Capacity

- **Total Story Points**: 5 pts
- **Backend**: 3 pts
- **Frontend**: 2 pts

## Risks & Dependencies

- #21과 #22는 독립적이므로 병렬 진행 가능
- Go 1.21+ 필요 (log/slog 사용) — 현재 go.mod에 go 1.22 설정되어 있어 문제 없음
- 외부 라이브러리 의존성 없음 (backend: 표준 라이브러리, frontend: 자체 구현)

## Deferred to Next Sprint

- #7 사용자 인증 시스템 (5pts) — #8의 선행 조건
- #8 요약 결과 저장 및 히스토리 조회 (5pts) — #7에 의존
