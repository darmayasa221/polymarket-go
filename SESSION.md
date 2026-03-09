# SESSION.md — Current Work State

## How to Use This File
- Read this at the START of every new session BEFORE touching any code
- Update this at the END of every session (or when stopping mid-phase)
- This + CLAUDE.md + git log = full context, no history needed

## Last Session Summary
**Date:** 2026-03-10
**Completed:** Full Polymarket API research — all 5 tasks done, all open questions resolved
**Branch:** main
**Artifacts created:**
- `docs/decisions/polymarket-api-summary.md` — complete API reference (Gamma, CLOB, Data, WS, RTDS, fees, auth, contracts, oracle)
- `docs/plans/2026-03-10-phase1-domain-plan.md` — Phase 1 + 2 implementation blueprint

**Next Action:** Phase 0 tooling setup, then Phase 1 commons, then Phase 2 domains

---

## Current Phase
**Phase 0: Project Setup** — NOT STARTED

## Phase Checklist
- [ ] Phase 0: Project Setup
  - [ ] git init + module init
  - [ ] Copy commons from go-base-framework
  - [ ] .golangci.yml (depguard layer rules)
  - [ ] lefthook.yml (pre-commit hooks)
  - [ ] Makefile
  - [ ] Dockerfile + docker-compose
  - [ ] .env.example
  - [ ] AGENTS.md confirmed with correct module path

- [ ] Phase 1: Commons
  - [ ] commons/timeutil — WindowStart/End, SecondsRemaining, Now()
  - [ ] commons/crypto — GenerateUUID, GenerateSalt (uint256)
  - [ ] commons/polyid — ConditionID, TokenID, OrderID, SlugID types
  - [ ] commons/slug — predictable slug builder (no API call)

- [ ] Phase 2: Domains
  - [ ] domains/market — Market entity, Asset enum (btc/eth/sol/xrp), Outcome enum (Up/Down)
  - [ ] domains/oracle — Price, PriceSource (Chainlink/Binance), PredictOutcome signal
  - [ ] domains/order — Order entity, EIP-712 signing, GTD expiration
  - [ ] domains/position — Position entity, UnrealisedPnL, RealisedPnL

- [ ] Phase 3: Applications
  - ⚠️ BLOCKED — requires research session first (see below)
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

## Phase 3 Pre-Requisite Research (REQUIRED before any Phase 3 code)

Before designing applications layer, researcher agent must answer:

1. **5m market liquidity** — real order book depth on live markets. Is $10-$100 tradeable without moving the price?
2. **Late-entry fillability** — at T-4:00 (last 60s), is there still an active order book?
3. **Actual feeRateBps** — call `/fee-rate` on a live 5m market. What is the real fee?
4. **Chainlink round selection** — which exact round resolves the market? Closest to windowEnd? Or last confirmed before?
5. **Settlement speed** — how many seconds after window close does USDC arrive?
6. **Historical volume** — check Data API for real trade counts on past 5m markets
7. **Bid/ask spread** — typical spread on 5m markets at different times in the window

Output: `docs/decisions/5m-market-mechanics.md` before any Phase 3 planning.

---

## Start Next Session With
```
Read SESSION.md and CLAUDE.md. Check git log --oneline.
Use superpowers:executing-plans on docs/plans/2026-03-10-phase0-1-2-implementation.md.
Start from Task 0 and execute every task in order.
```
