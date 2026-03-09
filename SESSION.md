# SESSION.md — Current Work State

## How to Use This File
- Read this at the START of every new session BEFORE touching any code
- Update this at the END of every session (or when stopping mid-phase)
- This + CLAUDE.md + git log = full context, no history needed

## Last Session Summary
**Date:** 2026-03-10
**Completed:** Phase 3 pre-work — research (7 questions answered) + Architect blueprint (96 files designed)
**Branch:** main
**Artifacts created:**
- `docs/decisions/5m-market-mechanics.md` — live API research, 7 questions answered
- `docs/plans/2026-03-10-phase3-applications-plan.md` — 11-task plan, 96 files, ready to execute

**Key research findings:**
- Fee: `base_fee: 1000` → parabolic `fee(p) = p × (1-p) × curveConstant`; peak ~156 bps at p=0.50
- Settlement: ~2-3 min after window close (Polygon 64-block finality)
- Liquidity: $10-$200 sweet spot; books always populated; late entry viable near 50/50
- Chainlink: timestamp-matched Data Stream reports at windowStart and windowEnd
- Volume: 100-500 trades/window; bursty; algorithmically driven

**Next Action:** Execute `docs/plans/2026-03-10-phase3-applications-plan.md` (11 tasks, 5 batches)

---

## Current Phase
**Phase 3: Applications** — READY (research + architecture complete)

## Phase Checklist
- [x] Phase 0: Project Setup
  - [x] depguard module paths fixed to polymarket-go
  - [x] .env.example with Polymarket env vars
  - [x] shopspring/decimal dependency added

- [x] Phase 1: Commons
  - [x] commons/timeutil — WindowStart/End, SecondsRemaining
  - [x] commons/crypto — GenerateSalt (uint256)
  - [x] commons/polyid — ConditionID, TokenID, OrderID, SlugID types
  - [x] commons/slug — predictable slug builder (no API call)

- [x] Phase 2: Domains
  - [x] domains/market — Market entity, Asset enum (btc/eth/sol/xrp), Outcome enum (Up/Down)
  - [x] domains/oracle — Price, PriceSource (Chainlink/Binance), PredictOutcome signal
  - [x] domains/order — Order entity, EIP-712 signing, GTD expiration
  - [x] domains/position — Position entity, UnrealisedPnL, RealisedPnL

- [ ] Phase 3: Applications
  - [x] Research complete — `docs/decisions/5m-market-mechanics.md`
  - [x] Architecture complete — `docs/plans/2026-03-10-phase3-applications-plan.md`
  - [ ] Implementation — 11 tasks, 96 files, READY
- [ ] Phase 4: Infrastructures
- [ ] Phase 5: Interfaces

---

## Domains — CONFIRMED

All domains defined. See `docs/plans/2026-03-10-phase1-domain-plan.md` for full specs.

| Domain | Purpose |
|--------|---------|
| `market` | Market discovery, lifecycle, window management |
| `oracle` | Price feeds (Chainlink + Binance), resolution signal |
| `order` | Order creation, EIP-712 signing, GTD logic |
| `position` | Position tracking, PnL calculation |

---

## Key Decisions — FINAL

| Decision | Value |
|----------|-------|
| Module | `github.com/darmayasa221/polymarket-go` |
| Architecture | Microservice, clean architecture |
| Base framework | `go-base-framework` clean-code branch |
| Market type | 5-minute Up/Down (BTC, ETH, SOL, XRP) |
| Outcomes | **"Up" / "Down"** — never "Yes"/"No" |
| Signature type | EOA = 0 (funder = EOA address, needs POL for gas) |
| USDC | USDC.e only — `0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174` |
| Slug | `{ticker}-updown-5m-{floor(unix/300)*300}` — predictable |
| WS connections | 3 — market, user, RTDS (separate keepalive intervals) |
| EIP-712 | Port from `https://github.com/Polymarket/clob-client/blob/main/src/signing/eip712.ts` |

---

## Key Constraints (from research)
- Heartbeat: POST every 5s or all orders auto-cancel after 10s
- GTD expiration: `now + 60 + seconds_remaining` (60s mandatory buffer)
- `feeRateBps`: fetch from `/fee-rate` before every order — never hardcode
- `tick_size_change` WS event: MUST handle — tick size changes when price > 0.96 or < 0.04
- Matching engine restarts: Tuesdays 7AM ET, HTTP 425 → exponential backoff
- `enable_order_book: true` — filter required before trading any market

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

## Phase 3 Research — COMPLETE
All 7 questions answered. See `docs/decisions/5m-market-mechanics.md`.

---

## Start Next Session With
```
Read SESSION.md and CLAUDE.md. Check git log --oneline.
Use superpowers:executing-plans on docs/plans/2026-03-10-phase3-applications-plan.md.
Start from Task 1 and execute every task in order (11 tasks, 5 batches).
```
