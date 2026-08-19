# Specification — CLI Scaffolding

| Field       | Value                                       |
| ----------- | ------------------------------------------- |
| SpecID      | RVF-EHVF8                                   |
| Title       | CLI Scaffolding                             |
| Status      | Draft                                       |
| Date        | 2026-08-18                                  |
| Author      | Runvil Contributors                         |
| Domain      | Frameworks — CLI                            |

## 1. Context

Adoption begins with creation. This specification defines how Runvil
bootstraps new CLI projects: generation of a Go module wired to the framework,
a working example command, and project conventions — reducing time-to-first-run.

## 2. Problem Statement

Starting a new Runvil CLI project means copying boilerplate, hand-wiring the
framework dependency, and re-deriving conventions for exit codes, output, and
configuration — slow, error-prone, and inconsistent across teams. Teams
converge on diverging project shapes instead of a repeatable, idiomatic
starting point.

## 3. Goals

- G1 — Generate a new, buildable Runvil CLI project in a single command.
- G2 — Keep generated projects minimal, idiomatic, and understandable.
- G3 — Ensure generated projects follow the framework conventions (exit codes, output, config).

## 4. Non-Goals

- NG1 — Package management, tagging, or release automation.
- NG2 — Interactive project wizards beyond minimal prompts/flags.
- NG3 — Generating full application templates beyond a CLI skeleton.

## 5. Requirements

### 5.1 Generation

| ID         | Requirement                                                       | Priority |
| ---------- | ----------------------------------------------------------------- | -------- |
| FRK-SC-001 | Provide `init`/`new` scaffolding commands for CLI projects.      | Should   |
| FRK-SC-002 | Require, at minimum, a project name; infer the module path when omitted. | Must |
| FRK-SC-003 | Emit a `main` package wired to the `cli` application model.      | Must     |
| FRK-SC-004 | Include a runnable example command after generation.             | Must     |

### 5.2 Output & Safety

| ID         | Requirement                                                       | Priority |
| ---------- | ----------------------------------------------------------------- | -------- |
| FRK-SC-005 | Generated projects pass `gofmt`, `go vet`, and `go test`.        | Must     |
| FRK-SC-006 | Refuse to overwrite an existing non-empty target directory.      | Must     |
| FRK-SC-007 | Report generated files clearly and deterministically.            | Should   |

## 6. Non-Functional Requirements

- NFR1 — **Memory safety.** The `unsafe` package must not be used.
- NFR2 — **Idempotence.** Re-running on the same target is safe and fails predictably.
- NFR3 — **Portability.** Linux, macOS, and Windows.
- NFR4 — **Quality gates.** `gofmt`, `go vet ./...`, and `go test ./...` must pass.

## 7. Success Criteria

- S1 — A generated project builds and runs its example command without edits.
- S2 — Generation never corrupts or overwrites pre-existing files.
- S3 — Generated output contains no placeholders or TODOs left unfilled.

## 8. Related Specifications

| SpecID    | Title                                           |
| --------- | ----------------------------------------------- |
| [RVF-M8SSR](./RVF-M8SSR-cli-application-model.md) | CLI Application Model                |
| [RVF-UUQ3X](./RVF-UUQ3X-cli-configuration.md)     | CLI Configuration                    |
| [RVF-QOFJK](./RVF-QOFJK-runvil-meta-framework.md) | Runvil Meta-Framework                |

## 9. References

- [RVL-4Y8UP](https://github.com/runvil/libs/blob/main/docs/specs/RVL-4Y8UP-runvil-libraries.md) — Runvil Libraries initial specification.