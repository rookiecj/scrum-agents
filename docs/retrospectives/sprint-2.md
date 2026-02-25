# Sprint 2 Retrospective (2026-03-02 ~ 2026-03-06)

## Sprint Summary
- **Goal**: 종류별 최적 요약 프롬프트 엔진 구현 및 PDF/트위터/뉴스레터 콘텐츠 타입 확장으로 Link Summarizer 핵심 기능 완성. 빌드 스크립트로 개발 생산성 확보. E2E 테스트로 품질 검증 강화.
- **Goal Achieved**: Yes — 원래 계획된 3개 티켓(18pts) 모두 완료 후, 미드스프린트 추가된 E2E 테스트 2개(6pts)까지 전량 완료.
- **Planned**: 24 points (5 tickets — 원래 18pts/3tickets + 미드스프린트 6pts/2tickets)
- **Completed**: 24 points (5 tickets)
- **Velocity**: 24 points
- **Completion Rate**: 100%

## Queue Metrics
- **QA Pass Rate**: 100% (5/5 passed, 0 rework)
- **Rework Count**: 0 tickets sent back for rework
- **Bottleneck Stage**: None — 모든 티켓이 원활히 파이프라인 통과
- **Avg Time in QA Queue**: 즉시 처리 (Sequential 모드에서 DEV 완료 직후 QA 진행)

## Queue Stage at Sprint Close
| Stage | Count | Tickets |
|-------|-------|---------|
| Verified (Done) | 5 | #5, #16, #3, #17, #18 |
| DEV Queue | 0 | — |
| In Progress | 0 | — |
| QA Queue | 0 | — |
| In Review | 0 | — |
| Blocked | 0 | — |

## Completed Work
- #5 [Story] 종류별 최적 요약 프롬프트 엔진 (8pts, Backend) — 7 JSON 프롬프트 템플릿, TemplateRegistry, Summarizer, /api/summarize 핸들러
- #16 [Task] 빌드 스크립트 작성 (2pts, Backend+Frontend) — Makefile with build, test, run, lint, clean, help
- #3 [Story] PDF/트위터/뉴스레터 콘텐츠 추출 (8pts, Backend) — PDFExtractor, TwitterExtractor, NewsletterExtractor + 페이월 감지
- #17 [Task] Backend E2E 테스트 작성 (3pts, Backend) — 16 E2E 테스트 (health, detect, extract, pipeline, error cases)
- #18 [Task] Frontend E2E 테스트 작성 (3pts, Frontend) — Playwright 9 E2E 테스트 (happy path, provider, progress, errors)

## Carry-over Items
- None

## What Went Well
- **100% 완료율 2스프린트 연속**: 원래 계획 3티켓(18pts) + 미드스프린트 추가 2티켓(6pts) = 24pts 전량 완료
- **Sprint 1 회고 반영**: PLAN.md/PROGRESS.md 자동 생성 프로세스 도입, Story Point 라벨 사용으로 벨로시티 정량 측정 가능해짐
- **핸드오프 기록 활용**: PROGRESS.md의 Handoff Notes가 티켓 간 컨텍스트 공유에 효과적으로 작동
- **미드스프린트 유연성**: `/sprint:add`로 E2E 테스트 티켓을 스프린트 중 추가하여 품질 향상 달성
- **코드 품질**: 모든 구현에 테이블 기반 테스트 작성, QA rework 0건
- **프롬프트 엔진 외부화**: 7개 JSON 템플릿을 코드와 분리하여 비개발자도 프롬프트 수정 가능한 구조
- **Twitter API 리스크 회피**: API 대신 OG meta tag scraping으로 우회 — 리스크 식별 후 적절한 대안 선택

## What Didn't Go Well
- **피처 브랜치 미병합**: #5, #16, #3 피처 브랜치가 main에 병합되지 않은 상태로 E2E 테스트 브랜치 작업 시 수동 머지 필요
- **스프린트 날짜 불일치 지속**: 라벨 날짜(03-02~03-06)와 실제 실행일(02-25) 불일치 — Sprint 1 회고에서 지적된 이슈 미해결
- **단일 세션 실행 지속**: Sequential 모드만 사용 — Parallel 모드 실증 테스트 미진행
- **API 엔드포인트 미등록**: `/api/classify`, `/api/summarize` 핸들러는 구현되었으나 main.go에 미등록 상태
- **Frontend 백엔드 연동 미검증**: Frontend E2E는 API mock으로만 테스트 — 실제 백엔드 연동 통합 테스트 부재

## Lessons Learned
- Story Point 기반 벨로시티 측정이 가능해져 스프린트 용량 계획이 정밀해짐 (Sprint 1: 측정 불가 → Sprint 2: 24pts)
- 미드스프린트 티켓 추가는 용량 여유가 있을 때 효과적 — `/sprint:add` 명령이 유용함
- 피처 브랜치를 main에 병합하는 시점을 QA 완료 직후로 명확히 정해야 함
- E2E 테스트를 별도 티켓으로 분리하면 품질 작업이 가시적으로 추적 가능
- PROGRESS.md 핸드오프 노트가 에이전트 간 정보 공유에 핵심 역할 수행

## Next Sprint Recommendations

### Carry-over Tickets (priority)

없음 (100% 완료)

### Suggested New Work
| # | Title | Priority | Points | Component | Recommendation |
|---|-------|----------|--------|-----------|----------------|
| #7 | 사용자 인증 시스템 | medium | 5 | backend | #8의 선행 조건. Sprint 2에서 deferred |
| #8 | 요약 결과 저장 및 히스토리 조회 | medium | 5 | backend+frontend | #7에 의존. 함께 진행 |
| new | API 엔드포인트 등록 (classify, summarize) | high | 2 | backend | main.go에 /api/classify, /api/summarize 등록 |
| new | 피처 브랜치 main 병합 및 CI 구성 | high | 2 | full-stack | 미병합 브랜치 정리 + GitHub Actions CI |
| new | Frontend-Backend 통합 E2E 테스트 | medium | 3 | full-stack | Mock 없이 실제 서버 연동 검증 |

### Process Improvements
- [ ] QA 완료 직후 피처 브랜치를 main에 자동 병합하는 워크플로우 도입 (PR + merge)
- [ ] 스프린트 라벨 날짜를 실제 실행 일정에 맞춰 설정 (Sprint 1 회고에서 미해결)
- [ ] Parallel 모드 실증 테스트 — 다음 스프린트에서 backend/frontend 티켓이 있을 때 시도
- [ ] `/sprint:start`에서 `sprint:end auto` 패턴처럼 인수 기반 자동화 확장

## Velocity Trend
| Sprint | Planned | Completed | Rate |
|--------|---------|-----------|------|
| Sprint 1 | 5 tickets | 5 tickets | 100% |
| Sprint 2 | 24 pts (5 tickets) | 24 pts (5 tickets) | 100% |
