# Specification — CLI Configuration

| Field       | Value                                       |
| ----------- | ------------------------------------------- |
| SpecID      | RVF-UUQ3X                                   |
| Title       | CLI Configuration                          |
| Status      | Draft                                       |
| Version     | 0.1.0                                       |
| Date        | 2026-08-18                                  |
| Author      | Runvil Contributors                         |
| Domain      | Frameworks — CLI                            |

## 1. Context

Runvil CLI applications source configuration from the environment using Go
standard-library primitives (`os` environment variables, `strconv` coercion).
This specification defines the convention: naming, typing, defaults, and
precedence.

## 2. Goals

- G1 — Standardize environment-variable naming across applications.
- G2 — Provide type-safe coercion with explicit defaults.
- G3 — Establish a single precedence order for flags, environment, and defaults.

## 3. Non-Goals

- NG1 — File-based configuration (TOML/YAML/JSON) in the initial phase.
- NG2 — Remote/secret management or vault integration.
- NG3 — A configuration schema/DTO validation framework.

## 4. Requirements

### 4.1 Naming & Sources

| ID         | Requirement                                                       | Priority |
| ---------- | ----------------------------------------------------------------- | -------- |
| FRK-CF-001 | Load configuration from environment variables via `os`.           | Must     |
| FRK-CF-002 | Use a `<APP>_` prefix (e.g. `GREET_`) with `SCREAMING_SNAKE_CASE` keys. | Must |
| FRK-CF-003 | Expose each setting with a documented default.                    | Must     |

### 4.2 Typing & Validation

| ID         | Requirement                                                       | Priority |
| ---------- | ----------------------------------------------------------------- | -------- |
| FRK-CF-004 | Coerce values to typed settings using `strconv`.                  | Must     |
| FRK-CF-005 | Invalid values surface as usage errors (exit code 2).             | Must     |
| FRK-CF-006 | Missing settings fall back to defaults without error.             | Must     |

### 4.3 Precedence

| ID         | Requirement                                                       | Priority |
| ---------- | ----------------------------------------------------------------- | -------- |
| FRK-CF-007 | Precedence order: flags > environment variables > defaults.       | Must     |

## 5. Non-Functional Requirements

- NFR1 — **Memory safety.** The `unsafe` package must not be used.
- NFR2 — **Simplicity.** No external configuration dependencies.
- NFR3 — **Portability.** Linux, macOS, and Windows.
- NFR4 — **Quality gates.** `gofmt`, `go vet ./...`, and `go test ./...` must pass.

## 6. Success Criteria

- S1 — Coercion and validity rules are covered by tests.
- S2 — The `greet` example demonstrates env-based configuration end-to-end.
- S3 — No framework package imports a non-stdlib configuration library.

## 7. Related Specifications

| SpecID    | Title                                           |
| --------- | ----------------------------------------------- |
| [RVF-M8SSR](./RVF-M8SSR-cli-application-model.md) | CLI Application Model                |
| [RVF-QZTY2](./RVF-QZTY2-cli-errors-diagnostics.md) | CLI Errors & Diagnostics             |
| [RVF-EHVF8](./RVF-EHVF8-cli-scaffolding.md)         | CLI Scaffolding                      |

## 8. References

- [RVF-QOFJK](./RVF-QOFJK-runvil-meta-framework.md) — Runvil Meta-Framework initial specification.