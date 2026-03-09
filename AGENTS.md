# AGENTS.md — Mandatory Conventions for polymarket-go

> Every agent MUST read this file before writing any code.

---

## Agent Roster & Invocation Rules

| Agent | When to invoke |
|-------|---------------|
| `researcher` | Before choosing any library or integrating any external API |
| `architect` | Before adding any new domain, service, or major feature |
| `executor` | For all implementation tasks |
| `knowledge` | When executor needs context about existing code |
| `thinker` | After each phase is complete |
| `monitor` | Surfaces blockers and decisions — not routine progress |

**Rule:** architect → executor order is MANDATORY. Never skip architect for new domains.

---

## Mandatory Skill Invocation Table

| Task | Skill to invoke |
|------|-----------------|
| Creating a domain entity | `.claude/skills/domains/entity.md` |
| Adding a value object | `.claude/skills/domains/value-object.md` |
| Adding a domain event | `.claude/skills/domains/events.md` |
| Defining a repository interface | `.claude/skills/domains/repository.md` |
| Creating a command use case | `.claude/skills/usecases/command.md` |
| Creating a query use case | `.claude/skills/usecases/query.md` |
| Creating an HTTP handler | `.claude/skills/interfaces/handler.md` |
| Adding middleware | `.claude/skills/interfaces/middleware.md` |
| Implementing a factory | `.claude/skills/patterns/factory.md` |
| Implementing a decorator | `.claude/skills/patterns/decorator.md` |
| Implementing a repository (infrastructure) | `.claude/skills/domains/repository.md` |
| Defining error keys | `.claude/skills/patterns/errors.md` |
| Using value objects | `.claude/skills/patterns/value-objects.md` |
| Domain events pattern | `.claude/skills/patterns/domain-events.md` |
| Unit of work pattern | `.claude/skills/patterns/unit-of-work.md` |
| CQRS structure | `.claude/skills/patterns/cqrs.md` |
| Wiring DI | `.claude/skills/architecture/providers.md` |
| Before any code | `.claude/skills/workflow/pre-work.md` |
| Development loop | `.claude/skills/workflow/development.md` |
| Creating a commit | `.claude/skills/workflow/commit.md` |

---

## Layer Hierarchy

```
┌─────────────────────────────────────────────────────────┐
│  interfaces  (HTTP handlers, middleware)                 │
│  ↓ imports: applications + domains + commons            │
├─────────────────────────┬───────────────────────────────┤
│  applications           │  infrastructures              │
│  (use cases, CQRS)      │  (DB, cache, JWT impls)       │
│  ↓ imports:             │  ↓ imports:                   │
│  domains + commons      │  domains + commons            │
├─────────────────────────┴───────────────────────────────┤
│  domains  (entities, repo interfaces, value objects)    │
│  ↓ imports: commons                                     │
├─────────────────────────────────────────────────────────┤
│  commons  (utilities, errors, constants)                │
│  ↓ imports: stdlib + intra-commons only                 │
└─────────────────────────────────────────────────────────┘

applications and infrastructures are PEER layers — neither imports the other.
```

Enforced by golangci-lint `depguard` in `.golangci.yml`.

---

## Forbidden Import Patterns

```go
// WRONG: commons importing domains
import "github.com/darmayasa221/polymarket-go/internal/domains/market"

// WRONG: domains importing applications
import "github.com/darmayasa221/polymarket-go/internal/applications/markets"

// WRONG: applications importing infrastructures
import "github.com/darmayasa221/polymarket-go/internal/infrastructures/sqlite"

// CORRECT: domains importing commons/errors/types
import "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
```

---

## Naming Conventions

### Packages
- All plural: `commons`, `domains`, `applications`, `infrastructures`, `interfaces`
- Entity packages named after what it IS: `market`, `position`, `order`
- NOT how it's created: never `newmarket`, `createposition`

### Files per Layer
| Layer | Files |
|-------|-------|
| Domain entity | `entity.go`, `new.go`, `validate.go`, `validation.go`, `errors.go` |
| Domain events | `names.go` (event name constants — REQUIRED), one file per event |
| Repository | `interfaces.go` (in `domains/{context}/repository/`), `repository.go` (in infra) |
| Use case | `interface.go`, `usecase.go`, `middleware.go`, `dto/input.go`, `dto/output.go`, `usecase_test.go` |
| Handler | `handler.go` (struct only), `routes.go` (RegisterRoutes only), one file per action |
| Middleware | each middleware in own sub-package: `auth/middleware.go`, etc. |

### Constants
| Prefix | Use |
|--------|-----|
| `Err` | Error keys: `ErrMarketNotFound = "MARKET.NOT_FOUND"` |
| `Msg` | Log/response messages |
| `Field` | Field names for validation |
| `Action` | Action names |
| `Reason` | Reason strings |
| `Prefix` | Prefix strings |
| `Type` | Type constants |
| `Purpose` | Purpose constants |

---

## Error Key Format

```
DOMAIN.ERROR_CODE
```

Examples:
```go
const (
    ErrMarketNotFound    = "MARKET.NOT_FOUND"
    ErrPositionInvalid   = "POSITION.INVALID"
    ErrOrderRejected     = "ORDER.REJECTED"
)
```

---

## SRP File Rules

```
domains/market/
  entity.go       → struct definition only
  new.go          → factory constructors only
  validate.go     → validation logic only
  validation.go   → validation constants only
  errors.go       → error key constants only
```

---

## CQRS Structure

```
applications/
  markets/
    commands/
      createmarket/
        interface.go
        usecase.go
        middleware.go
        usecase_test.go
        dto/input.go
        dto/output.go
    queries/
      getmarket/
        ...
      listmarkets/
        ...
```

---

## Verification Commands

```bash
make fmt          # format before every commit
make lint         # must pass before every commit
make test         # all tests must pass
make build        # must compile
make check        # all of the above
```

---

## Commit Message Format

```
<type>(<scope>): <description>
```

Types: `feat`, `fix`, `refactor`, `test`, `docs`, `chore`

Rules:
- Never add author
- Stage specific files only — never `git add .`
- Run `make check` before every commit
- Never use `--no-verify`
