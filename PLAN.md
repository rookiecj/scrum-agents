# Sprint 5 Plan (2026-03-09 ~ 2026-03-13)

## Sprint Goal

LLM 인프라를 확장하여 .env 기반 credential 관리를 도입하고, Gemini provider를 추가하며, API 키 미설정 provider를 UI에서 비활성화한다.

## Tickets

| # | Title | Type | Priority | Points | Component |
|---|-------|------|----------|--------|-----------|
| #25 | .env 파일 기반 credential 관리 | story | medium | 2pts | backend |
| #26 | Gemini LLM provider 추가 | story | medium | 5pts | backend+frontend |
| #27 | API 키 미설정 provider UI 비활성화 | story | medium | 3pts | backend+frontend |

## Sprint Capacity

- **Total Story Points**: 10 pts
- **Backend**: 10 pts (3 tickets 모두 backend 포함)
- **Frontend**: 8 pts (#26, #27)

## Execution Order (Sequential 권장)

티켓 간 의존성이 있으므로 순차 실행:

1. **#25** .env credential 관리 (2pts, backend only) — 선행 인프라
2. **#26** Gemini provider 추가 (5pts, backend+frontend) — #25 완료 후
3. **#27** Provider UI 비활성화 (3pts, backend+frontend) — #25, #26 완료 후

## Risks & Dependencies

- **#25 → #26 → #27 순차 의존성**: Parallel 모드 불가, Sequential 모드로 진행
- **Gemini API 접근**: GOOGLE_API_KEY 필요 — QA 시 실제 API 호출 또는 mock 필요
- **외부 라이브러리**: .env 로딩에 `godotenv` 사용 시 go.mod 의존성 추가

## Deferred to Next Sprint

- #7 사용자 인증 시스템 (5pts) — Sprint 2부터 5회 deferred
- #8 요약 결과 저장 및 히스토리 조회 (5pts) — #7에 의존
