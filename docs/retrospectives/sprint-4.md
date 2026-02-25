# Sprint 4 Retrospective (2026-03-02 ~ 2026-03-06)

## Sprint Summary
- **Goal**: 백엔드와 프론트엔드에 구조화된 로깅 정책을 적용하여, 일관된 로그 포맷과 레벨 기반 로깅 인프라를 확보한다.
- **Goal Achieved**: Yes — 2개 티켓 모두 완료. Backend log/slog 구조화 로깅, Frontend 로거 유틸리티 + ErrorBoundary 적용.
- **Planned**: 5 points (2 tickets)
- **Completed**: 5 points (2 tickets)
- **Velocity**: 5 points
- **Completion Rate**: 100%

## Queue Metrics
- **QA Pass Rate**: 100% (2/2 passed, 0 rework)
- **Rework Count**: 0 tickets sent back for rework
- **Bottleneck Stage**: None
- **Avg Time in QA Queue**: 즉시 처리 (Parallel 모드 — Dev 완료 후 QA Agent 바로 투입)

## Queue Stage at Sprint Close
| Stage | Count | Tickets |
|-------|-------|---------|
| Verified (Done) | 2 | #21, #22 |
| DEV Queue | 0 | — |
| In Progress | 0 | — |
| QA Queue | 0 | — |
| In Review | 0 | — |
| Blocked | 0 | — |

## Completed Work
- #21 [Story] Backend 구조화된 로깅 정책 적용 - log/slog (3pts, Backend) — `internal/logging` 패키지 신규, JSON handler, HTTP 미들웨어, 전 핸들러 slog 마이그레이션, 100% 테스트 커버리지
- #22 [Story] Frontend 로깅 유틸리티 및 에러 리포팅 적용 (2pts, Frontend) — `src/utils/logger.ts` 유틸리티, ErrorBoundary 컴포넌트, API 에러 로깅, 25개 신규 테스트

## Carry-over Items
- None

## What Went Well
- **100% 완료율 4스프린트 연속**: Sprint 1~4 모두 전량 완료. 안정적인 실행 패턴 지속
- **Parallel 모드 첫 성공**: Sprint 3 회고에서 제안된 Parallel 모드를 이번 스프린트에서 실증. Backend Dev + Frontend Dev가 독립 worktree에서 동시 작업하여 성공적으로 완료
- **Worktree 격리 모드 효과**: 각 에이전트가 독립 worktree에서 작업하여 git 충돌 없이 병렬 개발 가능 확인
- **QA 100% 일회 통과**: 두 티켓 모두 rework 없이 한 번에 QA 통과. 명확한 AC 덕분
- **Sprint 3 회고 액션 아이템 해소**: Parallel 모드 실증(해결), QA 후 즉시 main 병합(유지)
- **외부 의존성 제로**: 두 티켓 모두 외부 라이브러리 없이 구현(Go 표준 slog, 자체 logger 유틸리티)

## What Didn't Go Well
- **스프린트 라벨 날짜 불일치 지속**: Sprint 4 라벨은 `2026-03-02 ~ 2026-03-06`으로 설정했으나, 실제 실행은 2026-02-25에 수행. 4스프린트 연속 미해결. 스프린트 실행 시점에 라벨 날짜를 맞추는 프로세스 부재
- **Worktree PROGRESS.md 동기화 이슈**: Backend Dev의 worktree에서 PROGRESS.md 업데이트가 main 워킹카피에 자동 반영되지 않아, Scrum Master가 수동 머지 충돌 해결 필요. Frontend Dev의 업데이트만 main에 반영됨
- **Lint 경고 미해결**: Backend `logging_test.go`에 SA1012 (nil Context) lint 경고 잔존. 기능에 영향 없으나 코드 품질 표준에 미달
- **리모트 브랜치 미정리 지속**: Sprint 3에서 지적된 origin 리모트 브랜치 정리가 여전히 미수행

## Lessons Learned
- **Parallel 모드는 독립 티켓에 효과적**: 컴포넌트가 다른(backend/frontend) 독립 티켓은 Parallel 모드에 최적. 의존성 있는 티켓은 여전히 Sequential이 안전할 수 있음
- **Worktree 격리의 양면성**: git 충돌은 방지하지만, 공유 파일(PROGRESS.md) 동기화 문제 발생. 공유 상태 업데이트는 GitHub Issue/PR 등 외부 시스템 활용이 더 안정적
- **로깅 인프라는 초기 구축이 이상적**: 4번째 스프린트에 도입했지만, 1스프린트부터 있었으면 디버깅에 유용했을 것. 다음 프로젝트에서는 초기 설정에 포함 권장
- **AC 명확성 = QA 효율**: 두 티켓 모두 Given/When/Then 형식의 구체적 AC가 있어 QA가 명확한 기준으로 검증 가능. rework 0의 핵심 요인

## Next Sprint Recommendations

### Carry-over Tickets (priority)

없음 (100% 완료)

### Suggested New Work
| # | Title | Priority | Points | Component | Recommendation |
|---|-------|----------|--------|-----------|----------------|
| #7 | 사용자 인증 시스템 | medium | 5 | backend+frontend | #8의 선행 조건. Sprint 2부터 3회 deferred — 다음 스프린트에서 반드시 착수 |
| #8 | 요약 결과 저장 및 히스토리 조회 | medium | 5 | backend+frontend | #7에 의존. #7 완료 후 진행 가능. 같은 스프린트 포함 시 Sequential 권장 |
| new | SA1012 lint 경고 수정 | low | 1 | backend | `logging_test.go`에서 nil Context → `context.TODO()` 교체 |
| new | 리모트 브랜치 정리 | low | 1 | infra | Sprint 3부터 2회 연속 미수행. `git push origin --delete` 일괄 정리 |

### Process Improvements
- [ ] 스프린트 라벨 날짜를 `/sprint:start` 실행 시점에 자동 업데이트하도록 워크플로우 개선 (4스프린트 연속 미해결)
- [ ] Parallel 모드에서 PROGRESS.md 동기화 전략 개선: 각 에이전트가 GitHub Issue comment로 상태 업데이트하고, Scrum Master가 PROGRESS.md를 일괄 갱신하는 방식 검토
- [ ] 의존성 있는 티켓(#7→#8) Parallel 모드 전략 수립: #7 Backend+Frontend를 먼저 완료 후 #8 착수, 또는 #7 Backend 완료 시점에 #8 Backend 착수 등
- [ ] lint/vet 경고를 QA 검증 기준에 포함 여부 결정 (현재는 기능 테스트만 통과하면 pass)

## Velocity Trend
| Sprint | Planned | Completed | Rate |
|--------|---------|-----------|------|
| Sprint 1 | 5 tickets | 5 tickets | 100% |
| Sprint 2 | 24 pts (5 tickets) | 24 pts (5 tickets) | 100% |
| Sprint 3 | 4 pts (2 tickets) | 4 pts (2 tickets) | 100% |
| Sprint 4 | 5 pts (2 tickets) | 5 pts (2 tickets) | 100% |
