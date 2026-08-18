# Specification — Runvil Web Theming System

| Field       | Value                                       |
| ----------- | ------------------------------------------- |
| SpecID      | RVF-5K7PZ                                   |
| Title       | Runvil Web Theming System                  |
| Status      | Draft                                       |
| Version     | 0.1.0                                       |
| Date        | 2026-08-18                                  |
| Author      | Runvil Contributors                         |
| Domain      | Frameworks — Web                            |

## 1. Context

The Runvil web framework ships a small theming system for static sites: a
configurable `Theme` type that emits an inline script and a toggle button.
Sites follow the system color scheme by default, let users switch between
light and dark, and remember the choice without any server round-trip.

## 2. Problem Statement

Light/dark support is re-implemented per site with ad-hoc scripts, inconsistent
storage keys, and a visible flash of the wrong theme. The framework provides no
shared, testable primitive, so every consumer — including mdbind — reinvents
the same toggle.

## 3. Goals

- G1 — Provide a `Theme` type with light/dark/auto preferences.
- G2 — Persist the user's choice across visits.
- G3 — Apply the theme before first paint, avoiding a flash of the wrong mode.
- G4 — Ship a ready-made toggle button wired by the script.

## 4. Non-Goals

- NG1 — Server-side theming, cookies, or themes beyond light/dark.
- NG2 — Theme-aware image generation or color transformation.

## 5. Requirements

| ID         | Requirement                                                       | Priority |
| ---------- | ----------------------------------------------------------------- | -------- |
| FRK-TH-001 | Support `auto`, `light`, and `dark` preferences, with `auto` following the system color scheme. | Must |
| FRK-TH-002 | Persist the user's choice in `localStorage` under a configurable key. | Must |
| FRK-TH-003 | Apply the theme synchronously from an inline script so no wrong-theme flash occurs. | Must |
| FRK-TH-004 | Expose a toggle on `window.runvilTheme` and wire any element with `data-theme-toggle`. | Must |
| FRK-TH-005 | Provide ready-made toggle button markup (`Theme.Button`).          | Must |
| FRK-TH-006 | Keep the `color-scheme` meta element in sync with the applied theme. | Must |

## 6. Non-Functional Requirements

- NFR1 — **Memory safety.** The `unsafe` package must not be used.
- NFR2 — **Dependency-free.** No external JavaScript or CSS libraries.
- NFR3 — **Portability.** Works in evergreen browsers; degrades gracefully to the system scheme.
- NFR4 — **Quality gates.** `gofmt`, `go vet ./...`, and `go test ./...` must pass.

## 7. Success Criteria

- S1 — `Theme{}.Script()` and `Theme{}.Button()` render embeddable markup.
- S2 — The script reads `localStorage`, follows the system scheme for `auto`, and wires `data-theme-toggle`.
- S3 — A site built with the theme renders before first paint without a wrong-theme flash.

## 8. Related Specifications

| SpecID    | Title                                           |
| --------- | ----------------------------------------------- |
| [RVF-8G3WQ](./RVF-8G3WQ-runvil-web-framework.md) | Runvil Web Framework             |
| [RVM-5F9TL](https://github.com/runvil/mdbind/blob/main/docs/specs/RVM-5F9TL-mdbind-site-builder.md) | mdbind Site Builder |

## 9. References

- [RVF-QOFJK](./RVF-QOFJK-runvil-meta-framework.md) — Runvil Framework initial specification.