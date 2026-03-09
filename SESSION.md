# SESSION.md — Current Work State

## How to Use This File
- Read this at the START of every new session BEFORE touching any code
- Update this at the END of every session (or when stopping mid-phase)
- This + CLAUDE.md + git log = full context, no history needed

## Last Session Summary
**Date:** 2026-03-10
**Completed:** Project agent setup — CLAUDE.md, AGENTS.md, SESSION.md, .claude/ agents + skills
**Branch:** master
**Next Action:** Phase 0 — initialize Go module, copy commons from base framework, set up tooling

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

- [ ] Phase 1: Commons — copy + verify from base framework
- [ ] Phase 2: Domains — depends on confirmed domain list
- [ ] Phase 3: Applications
- [ ] Phase 4: Infrastructures
- [ ] Phase 5: Interfaces

## Domains Status
> PENDING — still being designed. Use architect + researcher agents to define.
> Run architect agent before starting Phase 2.

## Key Decisions Made
- Module: github.com/darmayasa221/polymarket-go
- Architecture: microservice approach
- Base framework: go-base-framework `clean-code` branch (kept as-is, not merged)
- Agent setup: executor, knowledge, monitor, thinker, architect, researcher
- Domains: still being analyzed — update CLAUDE.md when confirmed

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

## Start Next Session With
```
Read SESSION.md and CLAUDE.md. Check git log --oneline.
Current phase: Phase 0 (not started).
Next: start Phase 0 — go module init and tooling setup.
```
