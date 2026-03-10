# SESSION.md — Current Work State

## How to Use This File
- Read this at the START of every new session BEFORE touching any code
- Update this at the END of every session (or when stopping mid-phase)
- This + CLAUDE.md + git log = full context, no history needed

## Last Session Summary
**Date:** 2026-03-10
**Completed:** Phase 3 — ALL 11 Tasks DONE.
**Branch:** `feature/phase3-applications`
**Worktree:** `.worktrees/feature-phase3`

**What was done:**
- Task 9: CloseWindow (6 tests green) + GetWindowState (4 tests green) — committed
- Task 10: Portfolio context — all 5 use cases (21 tests green) — committed
- Task 11: Full integration build + lint — all clean, CLAUDE.md updated — committed

**Phase 3 is COMPLETE.**

---

## Current Phase
**Phase 4: Infrastructures** — NOT STARTED

## Phase Checklist
- [x] Phase 0: Project Setup
- [x] Phase 1: Commons
- [x] Phase 2: Domains
- [x] Phase 3: Applications
  - [x] Research complete
  - [x] Architecture complete
  - [x] Task 1: shared/ (windowstate, signal, feecalc, FeeRateProvider port)
  - [x] Task 2: pricing/ports/price_repository.go
  - [x] Task 3: pricing RecordPrice command
  - [x] Task 4: pricing GetCurrentSignal query
  - [x] Task 5: pricing ComputeFee query (formula: p*(1-p)*0.0625)
  - [x] Task 6: marketwatch ports + RefreshMarkets
  - [x] Task 7: marketwatch UpdateTickSize + GetActiveMarket + IsMarketTradeable
  - [x] Task 8: trading ports + StartWindow + Heartbeat
  - [x] Task 9: trading PlaceOrder + CancelOrder + CloseWindow + GetWindowState
  - [x] Task 10: portfolio context (all 5 commands/queries)
  - [x] Task 11: full integration build + lint + update CLAUDE.md/SESSION.md
- [ ] Phase 4: Infrastructures
- [ ] Phase 5: Interfaces

---

## Git State
- Main branch: `main`
- Feature branch: `feature/phase3-applications`
- Worktree: `.worktrees/feature-phase3`
- Last commit: `chore(progress): mark Phase 3 complete in CLAUDE.md and SESSION.md`

## Key Decisions — FINAL

| Decision | Value |
|----------|-------|
| Module | `github.com/darmayasa221/polymarket-go` |
| Architecture | Microservice, clean architecture |
| Market type | 5-minute Up/Down (BTC, ETH, SOL, XRP) |
| Outcomes | **"Up" / "Down"** — never "Yes"/"No" |
| Signature type | EOA = 0 |
| USDC | USDC.e only — `0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174` |
| Fee formula | `fee(p) = p * (1-p) * 0.0625` — VERIFIED against live API |
| EIP-712 signing | Applications layer computes hash ONLY. Interfaces layer signs with private key. |
| Polygon chain ID | 137 |
| Min order size | 5 tokens |

---

## Key Constraints (from research)
- Heartbeat: POST every 5s or all orders auto-cancel after 10s
- GTD expiration: `windowEnd + 60s` mandatory buffer
- `feeRateBps`: fetch from `/fee-rate` before every order — never hardcode
- `tick_size_change` WS event: MUST handle — valid tick sizes: 0.1, 0.01, 0.001, 0.0001
- Matching engine restarts: Tuesdays 7AM ET, HTTP 425 → exponential backoff

---

## Trading Philosophy — LOCKED (never change)

**"Defend the Money. Profit is Bonus."**

Priority order (highest to lowest):
1. Never lose more than $3 on any single trade (stop loss at −$0.20/token from entry)
2. Exit at +$2 profit minimum — do not chase $1.00 resolution
3. Never hold a losing position to resolution if exit price > $0.10
4. If total capital < $16 → stop all trading, protect reserve
5. Sit out windows with confidence < 0.30 — no signal = no trade

**Mid-window exit triggers (runs every 30s during window):**
- STOP LOSS: token price drops 20 cents below entry → sell immediately
- TAKE PROFIT: token price rises 20 cents above entry → sell, lock it
- TIME STOP: T+4:30 and position is underwater → sell, cut loss

**Why this matters (math):**
- Without stop loss: break-even = 56% win rate
- With stop loss (rescues 30% of losses): break-even drops to 48%
- That 8% difference is the margin between a bot that survives and one that busts

**Capital rules:**
- Max deployed per window: 80% of total capital ($32 if capital = $40)
- Reserve: always keep minimum 20% ($8) — covers settlement gap + next window entry
- Selective: only trade 2-3 assets per window when signal confidence > 0.30
- Never trade all 4 assets simultaneously unless all 4 have strong signals (rare)

---

## Agent Roles (always active)
- Executor: writes code following skills
- Knowledge: reads codebase for context
- Monitor: asks YOU when decisions needed (not routine work)
- Thinker: reviews after each phase
- Architect: designs domains/services BEFORE implementation
- Researcher: uses MCP tools to research libraries and external APIs

## Mandatory Workflow (NEVER skip)
1. Researcher → library/API research before any external dependency decision
2. Architect → domain/service blueprint before any domain code
3. Executor → implements following blueprint + skills
4. Thinker → reviews after each phase
5. Monitor → surfaces decisions to you, not routine progress

---

## Start Next Session With
```
Read SESSION.md and CLAUDE.md. Check git log --oneline.
Phase 3 is COMPLETE. Next: Phase 4 (Infrastructures).
Use superpowers:finishing-a-development-branch to wrap up feature/phase3-applications first.
Then write the Phase 4 plan using superpowers:writing-plans before any implementation.
Phase 4 will implement: repository adapters (PostgreSQL/Redis),
CLOB HTTP client, WebSocket listeners (market/user/RTDS),
Chainlink reader, heartbeat service, and chain-reading infrastructure.
```
