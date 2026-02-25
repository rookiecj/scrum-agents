# Sprint 3 Retrospective (2026-02-25 ~ 2026-02-27)

## Sprint Summary
- **Goal**: 기반 인프라 안정화: 미병합 피처 브랜치를 정리하고 CI 파이프라인을 구축하며, 미등록 API 엔드포인트를 라우터에 등록하여 전체 서비스 통합 기반을 확보한다.
- **Goal Achieved**: Yes — 2개 티켓 모두 완료. 피처 브랜치 전량 main 병합, CI 파이프라인 구축, classify/summarize 엔드포인트 등록 완료.
- **Planned**: 4 points (2 tickets)
- **Completed**: 4 points (2 tickets)
- **Velocity**: 4 points
- **Completion Rate**: 100%

## Queue Metrics
- **QA Pass Rate**: 100% (2/2 passed, 0 rework)
- **Rework Count**: 0 tickets sent back for rework
- **Bottleneck Stage**: None
- **Avg Time in QA Queue**: 즉시 처리 (Sequential 모드)

## Queue Stage at Sprint Close
| Stage | Count | Tickets |
|-------|-------|---------|
| Verified (Done) | 2 | #19, #20 |
| DEV Queue | 0 | — |
| In Progress | 0 | — |
| QA Queue | 0 | — |
| In Review | 0 | — |
| Blocked | 0 | — |

## Completed Work
- #19 [Task] 피처 브랜치 main 병합 및 CI 구성 (2pts, Backend+Frontend) — 3개 미머지 브랜치 반영, 10개 로컬 브랜치 정리, GitHub Actions CI workflow 추가
- #20 [Task] API 엔드포인트 등록 (2pts, Backend) — `/api/classify`, `/api/summarize` 라우트 등록, 7개 E2E 테스트 추가

## Carry-over Items
- None

## What Went Well
- **100% 완료율 3스프린트 연속**: Sprint 1~3 모두 전량 완료. 안정적인 실행 패턴 확립
- **Sprint 2 회고 이슈 해결**: 피처 브랜치 미병합(#19)과 API 미등록(#20) 문제를 이번 스프린트에서 직접 대응
- **기술 부채 해소**: 누적된 인프라 이슈(미병합 브랜치, CI 부재, 미등록 엔드포인트)를 전용 스프린트로 일괄 정리
- **CI 파이프라인 구축**: PR 생성 시 자동 빌드/테스트가 실행되어 코드 품질 게이트 확보
- **스프린트 규모 적정성**: 4pts의 소규모 스프린트로 기반 작업에 집중 — 오버커밋 없이 깔끔하게 완료
- **QA 완료 후 즉시 main 병합**: Sprint 2 회고에서 지적된 "병합 시점 불명확" 이슈를 이번 스프린트에서 실천

## What Didn't Go Well
- **스프린트 라벨 날짜 불일치 지속**: Sprint 3 라벨은 `2026-02-25 ~ 2026-02-27`로 설정했으나, 이전 스프린트들의 날짜도 `03-02~03-06`으로 불일치 상태. 3스프린트 연속 미해결
- **단일 세션(Sequential) 모드 지속**: Parallel 모드 실증 테스트를 여전히 진행하지 못함
- **리모트 브랜치 미정리**: 로컬 브랜치는 정리했으나 `origin/feature/1-url-type-detection` 등 리모트 브랜치는 남아있음
- **벨로시티 편차 큼**: Sprint 2(24pts) → Sprint 3(4pts)로 벨로시티 변동이 크지만, 의도적인 기반 작업 스프린트이므로 실제 문제는 아님

## Lessons Learned
- 기술 부채 해소를 전용 스프린트 목표로 설정하면 명확한 완료 기준과 집중도를 확보할 수 있음
- 회고에서 발견된 이슈를 다음 스프린트 티켓으로 직접 전환하면 실행력이 높아짐 (Sprint 2 회고 → Sprint 3 티켓)
- 소규모 스프린트도 유효함 — 모든 스프린트가 대규모일 필요 없음
- CI 구축은 프로젝트 초기에 해두는 것이 이상적이나, 늦더라도 한번 설정하면 이후 모든 스프린트에 효과

## Next Sprint Recommendations

### Carry-over Tickets (priority)

없음 (100% 완료)

### Suggested New Work
| # | Title | Priority | Points | Component | Recommendation |
|---|-------|----------|--------|-----------|----------------|
| #7 | 사용자 인증 시스템 | medium | 5 | backend+frontend | #8의 선행 조건. Sprint 2부터 deferred |
| #8 | 요약 결과 저장 및 히스토리 조회 | medium | 5 | backend+frontend | #7에 의존. 함께 진행 |
| new | Frontend-Backend 통합 E2E 테스트 | medium | 3 | full-stack | Mock 대신 실서버 연동 검증. Sprint 2 회고에서 제안 |
| new | 리모트 브랜치 정리 | low | 1 | infra | origin에 남아있는 오래된 피처 브랜치 삭제 |

### Process Improvements
- [ ] QA 완료 직후 피처 브랜치 main 병합을 프로세스로 고정 (Sprint 3에서 실천, 다음 스프린트도 유지)
- [ ] 스프린트 라벨 날짜를 실제 실행 일정에 맞춰 설정 (3스프린트 연속 미해결 — 다음 스프린트에서 반드시 수정)
- [ ] Parallel 모드 실증 테스트 — #7(backend)과 #8(frontend) 동시 진행 가능 시 시도
- [ ] `git push` 후 CI 결과 확인을 스프린트 워크플로우에 포함

## Velocity Trend
| Sprint | Planned | Completed | Rate |
|--------|---------|-----------|------|
| Sprint 1 | 5 tickets | 5 tickets | 100% |
| Sprint 2 | 24 pts (5 tickets) | 24 pts (5 tickets) | 100% |
| Sprint 3 | 4 pts (2 tickets) | 4 pts (2 tickets) | 100% |
