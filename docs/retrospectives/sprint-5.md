# Sprint 5 Retrospective (2026-03-09 ~ 2026-03-13)

## Sprint Summary
- **Goal**: LLM 인프라를 확장하여 .env 기반 credential 관리를 도입하고, Gemini provider를 추가하며, API 키 미설정 provider를 UI에서 비활성화한다.
- **Goal Achieved**: Yes — 3개 티켓 전량 완료. .env credential 관리, Gemini provider, Provider 가용성 UI 모두 구현 및 검증 완료.
- **Planned**: 10 points (3 tickets)
- **Completed**: 10 points (3 tickets)
- **Velocity**: 10 points
- **Completion Rate**: 100%

## Queue Metrics
- **QA Pass Rate**: 100% (3/3 passed, 0 rework)
- **Rework Count**: 0 tickets sent back for rework
- **Bottleneck Stage**: None
- **Avg Time in QA Queue**: 즉시 처리 (Sequential 모드 — 각 티켓 Dev 완료 즉시 QA 진행)

## Queue Stage at Sprint Close
| Stage | Count | Tickets |
|-------|-------|---------|
| Verified (Done) | 3 | #25, #26, #27 |
| DEV Queue | 0 | — |
| In Progress | 0 | — |
| QA Queue | 0 | — |
| In Review | 0 | — |
| Blocked | 0 | — |

## Completed Work
- #25 [Story] .env 파일 기반 credential 관리 (2pts, Backend) — godotenv 의존성, loadEnv 헬퍼 (graceful fallback), .env.example 문서화, 5개 테스트. PR #28 merged.
- #26 [Story] Gemini LLM provider 추가 (5pts, Backend+Frontend) — GeminiProvider (gemini.go), ProviderGemini 상수, DefaultGeminiConfig, Gemini 라디오 버튼 (UrlInput.tsx), api.ts 타입 업데이트, 4 backend + 1 frontend 테스트. PR #29 merged.
- #27 [Story] API 키 미설정 provider UI 비활성화 (3pts, Backend+Frontend) — GET /api/providers 엔드포인트, disabled 라디오 버튼 + tooltip, 첫 available provider 자동 선택, 미가용 시 경고 표시, 6 backend + 12 frontend 테스트. PR #30 merged.

## Carry-over Items
- None

## What Went Well
- **100% 완료율 5스프린트 연속**: Sprint 1~5 전량 완료. 팀의 추정 정확도와 실행력이 지속적으로 안정적임을 입증
- **Sequential 모드의 의존성 체인 처리 적합**: #25→#26→#27 순차 의존성을 Sequential 모드로 정확히 처리. Sprint 4 회고에서 "의존성 있는 티켓은 Sequential이 안전" 교훈을 적용한 결과
- **스프린트 최대 규모 성공 (10pts)**: 이전 최대치(Sprint 2: 24pts 5티켓, Sprint 4: 5pts 2티켓)와 비교하여, 의존성 체인이 있는 10pts 3티켓을 rework 없이 완수
- **QA 100% 일회 통과 유지**: 5스프린트 연속 rework 0. Given/When/Then AC 형식의 효과 지속 확인
- **외부 API 통합 성공**: Gemini provider 구현에서 Google Generative AI API (v1beta) 통합을 mock 기반 테스트로 검증. 실제 API 키 없이도 안전하게 개발/테스트 가능
- **Cross-component 일관성**: Backend /api/providers 엔드포인트와 Frontend provider selector가 깔끔하게 연동. Handoff Notes를 통한 context 공유 유효

## What Didn't Go Well
- **스프린트 라벨 날짜 불일치 5회 연속**: Sprint 5 라벨은 `2026-03-09 ~ 2026-03-13`이나 실제 실행은 2026-02-25. Sprint 1부터 미해결 — 프로세스 개선이 반복 제안되었으나 실행되지 않음
- **Deferred 티켓 누적**: #7(사용자 인증)이 Sprint 2부터 5회 연속 deferred. #8(히스토리 조회)은 #7에 의존하여 함께 계속 밀림. 기능적 부채가 누적 중
- **OpenAI provider 미완성**: `openaiKey` 변수를 읽지만 실제 OpenAI provider 인스턴스를 생성하지 않음 (main.go:41-44에서 로그만 출력). Gemini와 달리 provider가 등록되지 않아 /api/providers에서 available:false로만 표시
- **godotenv 외부 의존성 도입**: Sprint 4에서 "외부 의존성 제로"를 장점으로 꼽았으나, 이번 스프린트에서 godotenv를 도입. 표준 라이브러리만으로 .env 파싱이 가능했을 수 있음 — 의존성 최소화 원칙과 상충

## Lessons Learned
- **Sequential 모드는 의존성 체인에 필수**: Sprint 4(Parallel)와 Sprint 5(Sequential) 비교 경험을 통해, 티켓 간 의존성 여부가 모드 선택의 핵심 기준임을 재확인
- **Provider 패턴의 확장성 확인**: LLM adapter 패턴(Provider interface)이 새 provider 추가 시 효과적. Claude→OpenAI→Gemini로 확장하면서 기존 코드 변경 최소화
- **환경 변수 기반 feature toggle 유효**: API 키 유무로 provider 가용성을 동적 제어하는 패턴이 .env + /api/providers + UI disable의 3단 구조로 깔끔하게 구현됨
- **Handoff Notes의 가치**: Sequential 모드에서도 PROGRESS.md Handoff Notes가 context 누적에 유용. 각 티켓 완료 후 다음 티켓이 참조할 수 있는 기술 context 제공

## Next Sprint Recommendations

### Carry-over Tickets (priority)
없음 (100% 완료)

### Suggested New Work
| # | Title | Priority | Points | Component | Recommendation |
|---|-------|----------|--------|-----------|----------------|
| #7 | 사용자 인증 시스템 | medium | 5 | backend+frontend | **Sprint 2부터 5회 deferred — 반드시 착수 필요**. 더 이상 미룰 경우 기능적 부채 심각 |
| #8 | 요약 결과 저장 및 히스토리 조회 | medium | 5 | backend+frontend | #7에 의존. #7 완료 후 같은 스프린트 또는 다음 스프린트에서 진행 |
| new | OpenAI provider 완성 | low | 3 | backend | main.go에 openaiKey 로딩 로직이 있으나 provider 미등록. Gemini 패턴 참조하여 완성 |
| new | 리모트 브랜치 정리 | low | 1 | infra | Sprint 3부터 3회 연속 미수행 |

### Process Improvements
- [ ] 스프린트 라벨 날짜를 `/sprint:start` 실행 시점에 자동 업데이트 (5스프린트 연속 미해결 — 최우선 프로세스 개선)
- [ ] Deferred 티켓 정책 수립: N회 이상 deferred된 티켓은 자동 우선순위 상향 또는 drop 결정 강제
- [ ] lint/vet 경고를 QA 기준에 포함 (Sprint 4 미해결 — SA1012 등)
- [ ] 외부 의존성 도입 시 표준 라이브러리 대안 검토 절차 추가

## Velocity Trend
| Sprint | Planned | Completed | Rate |
|--------|---------|-----------|------|
| Sprint 1 | 5 tickets | 5 tickets | 100% |
| Sprint 2 | 24 pts (5 tickets) | 24 pts (5 tickets) | 100% |
| Sprint 3 | 4 pts (2 tickets) | 4 pts (2 tickets) | 100% |
| Sprint 4 | 5 pts (2 tickets) | 5 pts (2 tickets) | 100% |
| Sprint 5 | 10 pts (3 tickets) | 10 pts (3 tickets) | 100% |
