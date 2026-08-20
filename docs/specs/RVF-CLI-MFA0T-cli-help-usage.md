# Specification — CLI Help & Usage

| Field       | Value                                       |
| ----------- | ------------------------------------------- |
| SpecID      | RVF-MFA0T                              |
| Title       | CLI Help & Usage                            |
| Status      | Draft                                       |
| Date        | 2026-08-18                                  |
| Author      | Runvil Contributors                         |
| Domain      | Frameworks — CLI                            |

## 1. Context

Help and usage output is a first-class part of a CLI's usability. This
specification defines how Runvil CLI commands expose help: an explicit `help`
sub-command, automatic `-h`/`--help`, and consistent usage text rendered
through `term`.

## 2. Problem Statement

Help and usage text is typically hand-maintained, drifts out of sync with the
flags and commands it documents, and is reachable through inconsistent routes
(`-h`, `--help`, `help`, or nothing at all). Users cannot discover a tool's
interface reliably, and nested subcommands become undiscoverable. Automation
and operators likewise have no stable way to interrogate a CLI's surface.

## 3. Goals

- G1 — Provide consistent, predictable help output across all commands.
- G2 — Make help available at every level of a command hierarchy.
- G3 — Render help through the `term` conventions (colors degrade gracefully).

## 4. Non-Goals

- NG1 — Downloadable man pages or shell completions in the initial phase.
- NG2 — Web-based documentation generation.
- NG3 — Localized/translated help strings.

## 5. Requirements

| ID         | Requirement                                                       | Priority |
| ---------- | ----------------------------------------------------------------- | -------- |
| FRK-HP-001 | Provide automatic `-h` and `--help` handling for every command.   | Must     |
| FRK-HP-002 | Provide an explicit `help [command...]` sub-command.              | Must     |
| FRK-HP-003 | Top-level help lists the application version, usage line, and commands. | Must |
| FRK-HP-004 | Per-command help lists its options, arguments, and description.   | Must     |
| FRK-HP-005 | Unknown command help requests resolve to the help exit code (2).  | Must     |
| FRK-HP-006 | Help text is derived from command/flag definitions, not hand-authored strings. | Should |
| FRK-HP-007 | Help output honors the `term` color policy (plain text when unsupported). | Must |

## 6. Non-Functional Requirements

- NFR1 — **Memory safety.** The `unsafe` package must not be used.
- NFR2 — **Performance.** Help generation must stay near-instant.
- NFR3 — **Portability.** Linux, macOS, and Windows.
- NFR4 — **Quality gates.** `gofmt`, `go vet ./...`, and `go test ./...` must pass.

## 7. Success Criteria

- S1 — `help` and `-h` produce equivalent, deterministic output for a given app.
- S2 — Nested-command help is reachable and correct.
- S3 — Help output contains no escape sequences when colors are disabled.

## 8. Related Specifications

| SpecID      | RVF-MFA0T                              |
| --------- | ----------------------------------------------- |
| [RVF-5XJFC](./RVF-CLI-5XJFC-cli-application-model.md) | CLI Application Model                |
| [RVF-NPFSE](./RVF-CLI-NPFSE-cli-output-formatting.md) | CLI Output & Formatting             |
| [RVL-R934Y](https://github.com/runvil/libs/blob/main/docs/specs/RVL-TERM-R934Y-terminal-io-rendering.md) | Terminal I/O & Rendering            |

## 9. References

- [RVF-CMBZJ](./RVF-META-CMBZJ-runvil-meta-framework.md) — Runvil Meta-Framework initial specification.