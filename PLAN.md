# Sprint 3 Plan (2026-02-25 ~ 2026-02-27)

## Sprint Goal

기반 인프라 안정화: 미병합 피처 브랜치를 정리하고 CI 파이프라인을 구축하며, 미등록 API 엔드포인트를 라우터에 등록하여 전체 서비스 통합 기반을 확보한다.

## Tickets

| # | Title | Type | Priority | Points | Component |
|---|-------|------|----------|--------|-----------|
| #19 | 피처 브랜치 main 병합 및 CI 구성 | task | high | 2 | backend+frontend |
| #20 | API 엔드포인트 등록 (classify, summarize) | task | high | 2 | backend |

## Sprint Capacity

- **Total Story Points**: 4 pts
- **Backend**: 4 pts
- **Frontend**: 2 pts (#19 CI 구성에 포함)

## Risks & Dependencies

- **#19 브랜치 충돌 가능성**: 수동 머지된 브랜치와 main 간 diff 확인 필요. 코드는 이미 반영되었으므로 리스크 낮음.
- **#20 → #19 순서 권장**: main이 정리된 후 API 등록 작업 진행이 깔끔하나, 독립 진행도 가능.
- **핸들러 의존성**: classify/summarize 핸들러의 의존성(Classifier, Summarizer 인스턴스) 주입 방법 확인 필요.

## Deferred to Next Sprint

- #7 사용자 인증 시스템 (5pts) — #8의 선행 조건
- #8 요약 결과 저장 및 히스토리 조회 (5pts) — #7에 의존
- Frontend-Backend 통합 E2E 테스트 (3pts)
