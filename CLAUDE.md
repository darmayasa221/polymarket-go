# CLAUDE.md — polymarket-go

## Current Build State
- [ ] Phase 0: Project Setup
- [ ] Phase 1: Commons
- [ ] Phase 2: Domains
- [ ] Phase 3: Applications
- [ ] Phase 4: Infrastructures
- [ ] Phase 5: Interfaces

## Project
Polymarket trading bot — microservice architecture.
Module: `github.com/darmayasa221/polymarket-go`
Branch convention: `feature/*`, `fix/*` — PRs to `main`

## Confirmed Domains
| Domain | Purpose |
|--------|---------|
| `market` | Market discovery, lifecycle, 5-minute window management |
| `oracle` | Price feeds (Chainlink + Binance), resolution signal |
| `order` | Order creation, EIP-712 signing, GTD expiration logic |
| `position` | Position tracking, unrealised/realised PnL |

## Domain Constraints (from API research)
- Outcomes are **"Up" / "Down"** — never "Yes" / "No"
- Slug pattern: `{ticker}-updown-5m-{floor(unix/300)*300}` — predictable, no API call needed
- `feeRateBps` — always fetched live from CLOB `/fee-rate`, never hardcoded
- GTD expiration: `now + 60 + seconds_remaining_in_window`
- Tick size can change mid-market — handle `tick_size_change` WS event
- Three WebSocket connections required: market, user, RTDS (separate keepalive intervals)

## API Research
Complete. See `docs/decisions/polymarket-api-summary.md`.

## Layer Rules (NEVER violate)
```
commons         → stdlib + intra-commons only
domains         → commons only
applications    → domains + commons (no infrastructures/interfaces)
infrastructures → domains + commons (no interfaces)
interfaces      → applications + domains + commons
```
Enforced by golangci-lint depguard rules.

## Critical Rules
1. Entity named after what it IS: `market/` not `newmarket/`
2. SRP at file level: entity.go, new.go, validate.go, validation.go, errors.go
3. CQRS: applications/{context}/commands/ and applications/{context}/queries/
4. Value Objects — never raw string/int for domain concepts
5. Factory is ONLY way to create entity: `New()` always calls `Validate()`
6. Interface assertion always: `var _ Interface = (*Impl)(nil)`
7. Use `timeutil.Now()` not `time.Now()`
8. Use `crypto.GenerateUUID()` for IDs
9. Error key format: `DOMAIN.ERROR_CODE`
10. No author in commits

## Microservice Rules
- Each service owns its data — no cross-service DB joins
- Services communicate via events (async) or typed interfaces (sync)
- Never share domain entities across service boundaries — use DTOs/contracts
- Each service has its own `cmd/{service}/` entry point
- Document every service boundary decision in `docs/decisions/`

## Constant Prefixes
Err=error keys, Msg=messages, Field=field names, Action=actions,
Reason=reasons, Prefix=prefixes, Type=type constants, Purpose=purpose constants

## Build Order (per service)
commons → domains → applications → infrastructures → interfaces

## Source Framework
Based on: `github.com/darmayasa221/go-base-framework` (clean-code branch)
Commons layer copied directly — do not diverge from base patterns.

## Skills Reference
See AGENTS.md for mandatory skill invocation table.
