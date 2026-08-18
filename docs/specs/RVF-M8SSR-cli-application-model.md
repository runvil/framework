# Specification — CLI Application Model

| Field       | Value                                       |
| ----------- | ------------------------------------------- |
| SpecID      | RVF-M8SSR                                   |
| Title       | CLI Application Model                       |
| Status      | Draft                                       |
| Version     | 0.1.0                                       |
| Date        | 2026-08-18                                  |
| Author      | Runvil Contributors                         |
| Domain      | Frameworks — CLI                            |

## 1. Context

The `cli` package is the heart of the Runvil CLI ecosystem: a framework-level
application model built on the Go standard library, integrating argument
parsing (`flag`), logging (`log/slog`), terminal rendering (`term`), and the
shared error/exit-code model (`core`).

This specification defines the application model: `App`, `Command`, dispatch,
and lifecycle.

## 2. Goals

- G1 — Model CLI applications as `App` + `Command` composition.
- G2 — Support nested command/subcommand hierarchies.
- G3 — Delegate argument parsing entirely to the standard library `flag`.
- G4 — Guarantee that every execution path yields a canonical exit code.

## 3. Non-Goals

- NG1 — A bespoke argument parser (stdlib `flag` is mandated).
- NG2 — A plugin/extension runtime in the initial phase.
- NG3 — Async execution scheduling (deferred to the async phase).

## 4. Requirements

### 4.1 Application & Command Model

| ID         | Requirement                                                       | Priority |
| ---------- | ----------------------------------------------------------------- | -------- |
| FRK-AM-001 | Provide an `App` with a name, version, and known commands.        | Must     |
| FRK-AM-002 | Provide a `Command` with a name, description, flag registration, and handler. | Must |
| FRK-AM-003 | Support registering commands and returning the `App` for chaining. | Must     |
| FRK-AM-004 | Support nested command hierarchies (sub-commands under a parent). | Must     |

### 4.2 Dispatch & Execution

| ID         | Requirement                                                       | Priority |
| ---------- | ----------------------------------------------------------------- | -------- |
| FRK-AM-005 | Dispatch on the first token; unknown commands yield the usage exit code. | Must |
| FRK-AM-006 | Parse remaining tokens with `flag.FlagSet` per command; parse failures yield the usage exit code. | Must |
| FRK-AM-007 | Every handler returns an exit code; `Run` returns it to the caller. | Must     |
| FRK-AM-008 | Empty invocation renders help and returns the usage exit code.    | Must     |

### 4.3 Lifecycle & Diagnostics

| ID         | Requirement                                                       | Priority |
| ---------- | ----------------------------------------------------------------- | -------- |
| FRK-AM-009 | Emit structured logs through `log/slog` (dispatch, errors).       | Must     |
| FRK-AM-010 | Keep the `cli` package free of panics across dispatch and parsing. | Must    |
| FRK-AM-011 | Provide the framework meta-package as the single-dependency entry point. | Must |

## 5. Non-Functional Requirements

- NFR1 — **Memory safety.** The `unsafe` package must not be used.
- NFR2 — **Performance.** Dispatch overhead must stay in the microseconds.
- NFR3 — **Portability.** Linux, macOS, and Windows.
- NFR4 — **Quality gates.** `gofmt`, `go vet ./...`, and `go test ./...` must pass.

## 6. Success Criteria

- S1 — Nested commands dispatch correctly and are covered by tests.
- S2 — All non-success paths resolve to `ExitCodeUsage` (2) or `ExitCodeFailure` (1), never raw integers.
- S3 — Package documentation covers every exported identifier.
- S4 — The `greet` example demonstrates the full model end-to-end.

## 7. Related Specifications

| SpecID    | Title                                           |
| --------- | ----------------------------------------------- |
| [RVF-QOFJK](./RVF-QOFJK-runvil-meta-framework.md) | Runvil Meta-Framework                 |
| [RVF-LJWEB](./RVF-LJWEB-cli-help-usage.md)         | CLI Help & Usage                     |
| [RVF-QZTY2](./RVF-QZTY2-cli-errors-diagnostics.md) | CLI Errors & Diagnostics             |
| [RVL-CHBZ4](https://github.com/runvil/runvil-libs/blob/main/docs/specs/RVL-CHBZ4-errors-exit-codes.md) | Core Errors & Exit Codes |

## 8. References

- [RVL-4Y8UP](https://github.com/runvil/runvil-libs/blob/main/docs/specs/RVL-4Y8UP-runvil-libraries.md) — Runvil Libraries initial specification.
- [RVL-N459G](https://github.com/runvil/runvil-libs/blob/main/docs/specs/RVL-N459G-terminal-io-rendering.md) — Terminal I/O & Rendering.