# SESSION.md — Current Work State

## How to Use This File
- Read this at the START of every new session BEFORE touching any code
- Update this at the END of every session (or when stopping mid-phase)
- This + CLAUDE.md + git log = full context, no history needed

## Last Session Summary
**Date:** 2026-03-10
**Completed:** Phase 4 — ALL 16 Tasks DONE.
**Branch:** `feature/phase4-infrastructures`
**Worktree:** `.worktrees/feature-phase4`

**What was done:**
- Task 1: PostgreSQL adapter + schema migration
- Task 2: Domain Reconstitute constructors (market, order, position)
- Tasks 3-6: PostgreSQL repositories (pricing, marketwatch, trading, portfolio)
- Task 7: Redis WindowStateStore
- Tasks 8-11: CLOB HTTP client (auth, FeeRateProvider, OrderSubmitter, HeartbeatSender)
- Task 12: Gamma API MarketSource
- Tasks 13-15: WebSocket handlers (RTDS, Market, User)
- Task 16: Config extension (PostgreSQL/CLOB/Gamma) + polymarket container providers

**Phase 4 is COMPLETE.**

---

## Current Phase
**Phase 5: Interfaces** — NOT STARTED

## Phase Checklist
- [x] Phase 0: Project Setup
- [x] Phase 1: Commons
- [x] Phase 2: Domains
- [x] Phase 3: Applications
- [x] Phase 4: Infrastructures
  - [x] Task 1: PostgreSQL adapter + schema
  - [x] Task 2: Domain Reconstitute constructors
  - [x] Task 3: Price repository (PostgreSQL)
  - [x] Task 4: Market repository (PostgreSQL)
  - [x] Task 5: Order repository (PostgreSQL)
  - [x] Task 6: Position repository (PostgreSQL)
  - [x] Task 7: Redis WindowStateStore
  - [x] Task 8: CLOB HTTP client base (L2 auth)
  - [x] Task 9: FeeRateProvider
  - [x] Task 10: OrderSubmitter
  - [x] Task 11: HeartbeatSender
  - [x] Task 12: Gamma API MarketSource
  - [x] Task 13: RTDS WebSocket handler
  - [x] Task 14: Market WebSocket handler
  - [x] Task 15: User WebSocket handler
  - [x] Task 16: Config + container wiring + full build
- [ ] Phase 5: Interfaces

---

## Git State
- Main branch: `main`
- Feature branch: `feature/phase4-infrastructures`
- Worktree: `.worktrees/feature-phase4`
- Last commit: `feat(infra): extend config + wire polymarket providers into container`

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
Phase 4 is COMPLETE. Next: Phase 5 (Interfaces).
Use superpowers:finishing-a-development-branch to wrap up feature/phase4-infrastructures first.
Then write the Phase 5 plan using superpowers:writing-plans before any implementation.
Phase 5 will implement: HTTP handlers for Polymarket trading bot operations,
CLI entry points (cmd/bot/), and wiring all infrastructure adapters into the running bot.
```
