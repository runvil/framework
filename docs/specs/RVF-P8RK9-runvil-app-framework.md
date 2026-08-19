# Specification — Runvil App Framework

| Field       | Value                                       |
| ----------- | ------------------------------------------- |
| SpecID      | RVF-P8RK9                                   |
| Title       | Runvil App Framework — Monolith Assembly    |
| Status      | Draft                                       |
| Date        | 2026-08-19                                  |
| Author      | Runvil Contributors                         |
| Domain      | Frameworks — App / Monolith                 |

## 1. Context

The full stack is now addressable: static sites (RVF-PN41Q), JSON APIs
(RVF-H3QD8), and server-rendered frontends (RVF-L6NJ5). **Monolith Assembly**
defines how a *project* becomes one runnable application: the project's
`main.go` composes `web.App` from routing, middleware, pages, static assets, and
configuration, then ships as a single binary that serves HTML + API + assets.

This spec is the contract between `framework` and the projects that assemble
apps. It supersedes the "no per-project `cmd`" convention **for app projects
only**: static/book sites keep the config-driven `runvil build` path, while
fullstack apps own an explicit `main.go` so backend logic lives in Go, is
`go run`-able, and has no generated or hidden code.

## 2. Problem Statement

Runvil can build sites but not *applications*. A fullstack monolith needs:

- Backend logic (handlers, middleware, DB wiring) expressed in Go — which the
  config-only model cannot host.
- Server-rendered pages and JSON APIs in one binary.
- Production-ready lifecycle beyond `http.ListenAndServe`.
- A single, documented assembly point so `runvil run`/`dev` (RVN-K2SQ7) can
  drive any project predictably.

Existing conventions were built for static delivery; extending them with a
small, explicit `main.go` per app keeps projects idiomatic without reintroducing
the "magic main" that other fullstack Go frameworks were criticized for
(RVF-ZK9LQ §1.1).

## 3. Goals

- G1 — Define the app-project contract: an explicit, idiomatic `main.go` that builds a `web.App`.
- G2 — Keep static/book sites entirely config-driven (unchanged convention); exempt only app projects from the no-`cmd` rule.
- G3 — Make app assembly declarative enough that `runvil run`/`dev` can predict project shape and lifecycle.
- G4 — Support configuration via `libs/config` + `libs/validate` with env overrides.
- G5 — Support single-binary deployment with embedded assets, static export, and API in one artifact.
- G6 — Keep the framework dependency direction: projects → framework → libs; never reverse.

## 4. Non-Goals

- NG1 — Generating `main.go` or any build-time code (NG6 of RVF-ZK9LQ still holds).
- NG2 — A plugin/hot-swap system for app components.
- NG3 — Microservices orchestration; this is single-process monolith.
- NG4 — Storage/ORM coupling: DB choice stays a project decision.
- NG5 — Changing behavior of the existing static/book `runvil build` path.

## 5. Requirements

### 5.1 App Project Contract

| ID          | Requirement                                                       | Priority |
| ----------- | ----------------------------------------------------------------- | -------- |
| FRK-APP-001 | An app project ships an explicit `main.go` (package `main`) that builds and runs its `web.App`. | Must |
| FRK-APP-002 | App projects may also live under `cmd/<name>`, keeping `main` thin and logic in packages. | Should |
| FRK-APP-003 | The project shape is detectable: a `main.go`/`cmd/*` with `web.App` marks an *app*; `runvil.yaml`+`ssg:` marks a *site*; `manuscript/` marks a *book*. | Must |
| FRK-APP-004 | No generated code, no injected `main`, no reflection-based app assembly. | Must |

### 5.2 Assembly & Configuration

| ID          | Requirement                                                       | Priority |
| ----------- | ----------------------------------------------------------------- | -------- |
| FRK-APP-010 | Provide `web.NewApp()` as the sole assembly entry; middleware/API via RVF-H3QD8, pages via RVF-L6NJ5. | Must |
| FRK-APP-011 | App configuration loads through `libs/config` (RVL-X7C4M) and validates through `libs/validate` (RVL-B9TW2) before use. | Must |
| FRK-APP-012 | Config precedence is env > file > zero (as specified in RVL-X7C4M §5.3). | Must |
| FRK-APP-013 | A startup configuration/validation error exits non-zero with a clear message (mapped through `core` exit codes). | Must |

### 5.3 Rendering & Assets

| ID          | Requirement                                                       | Priority |
| ----------- | ----------------------------------------------------------------- | -------- |
| FRK-APP-020 | App pages reuse the `ui` registry (RVF-ZK9LQ): built-in components out of the box, custom components registerable. | Must |
| FRK-APP-021 | Apps may embed the SSG-exported `site/` (built via RVF-PN41Q) and serve it with `App.Static`, keeping docs/marketing pages static. | Should |
| FRK-APP-022 | All frontend rendering flows through the unified SSR render path (RVF-L6NJ5 FRK-SRV-014). | Must |
| FRK-APP-050 | Apps may ship compiled client bundles (JS/TS output) under `public/` or `site/` and serve them via `App.Static`; the framework never compiles, transforms, or blocks them. | Must |
| FRK-APP-051 | The default monolith remains zero-JS; a client integration is an opt-in, additive layer (RVF-2TK4X) that never alters server rendering or the theme-toggle default. | Must |
| FRK-APP-052 | `examples/app` must demonstrate serving a static `public/` with a `<script>` mount point alongside server-rendered pages, kept byte-identical without the script. | Should |

### 5.4 Lifecycle & Deployment

| ID          | Requirement                                                       | Priority |
| ----------- | ----------------------------------------------------------------- | -------- |
| FRK-APP-030 | `main.go` calls `App.Run(addr)` for graceful start/shutdown (RVF-L6NJ5 FRK-SRV-021). | Must |
| FRK-APP-031 | The deliverable is one binary: embedded assets + API + SSR on one listener. | Must |
| FRK-APP-032 | Every exported identifier composing apps has a doc comment (rendered by `go doc`). | Must |

### 5.5 Org & Testing

| ID          | Requirement                                                       | Priority |
| ----------- | ----------------------------------------------------------------- | -------- |
| FRK-APP-040 | Provide `examples/app` demonstrating a monolith: config, DTO + validation, middleware, dynamic + static pages, embedded assets. | Must |
| FRK-APP-041 | Quality gates (`gofmt`, `go vet ./...`, `go test ./...`) pass for `examples/app`. | Must |

## 6. Non-Functional Requirements

- NFR1 — **Memory safety.** The `unsafe` package must not be used.
- NFR2 — **Dependencies.** Framework uses stdlib + `libs`; apps add only what their backend needs (DB drivers, etc.).
- NFR3 — **Portability.** Linux, macOS, Windows.
- NFR4 — **Quality gates.** `gofmt`, `go vet ./...`, `go test ./...` pass.
- NFR5 — **Determinism.** Identical config/env yield identical app startup.

## 7. Success Criteria

- S1 — `examples/app` starts from `go run .`, serves a dynamic page, a JSON endpoint, and an embedded static page on one port.
- S2 — `runvil identify` (RVN-K2SQ7) classifies the example as an *app* automatically.
- S3 — Rebuilding `site/` with RVF-PN41Q and embedding it does not change app startup.
- S4 — A misconfigured app exits non-zero with a before-first-request error message.

## 8. Related Specifications

| SpecID    | Title                                           |
| --------- | ----------------------------------------------- |
| [RVF-8G3WQ](./RVF-8G3WQ-runvil-web-framework.md) | Runvil Web Framework (host) |
| [RVF-H3QD8](./RVF-H3QD8-http-api-pipeline.md) | HTTP & API Pipeline |
| [RVF-L6NJ5](./RVF-L6NJ5-server-frontend-pipeline.md) | Server & Frontend Rendering Pipeline |
| [RVF-PN41Q](./RVF-PN41Q-static-site-generator.md) | Static Site Generator (site/ export) |
| [RVF-ZK9LQ](./RVF-ZK9LQ-layout-ui-system.md) | Layout & UI System |
| [RVF-2TK4X](./RVF-2TK4X-js-ts-framework-integration.md) | JS/TS Framework Integration (client asset seam) |
| [RVN-K2SQ7](https://github.com/runvil/runvil/blob/main/docs/specs/RVN-K2SQ7-runvil-run-dev-deploy.md) | runvil run/dev/deploy |
| [RVL-X7C4M](https://github.com/runvil/libs/blob/main/docs/specs/RVL-X7C4M-configuration-loading.md) | Configuration Loading |
| [RVL-B9TW2](https://github.com/runvil/libs/blob/main/docs/specs/RVL-B9TW2-struct-validation.md) | Struct Validation |

## 9. References

- [RVF-QOFJK](./RVF-QOFJK-runvil-meta-framework.md) — Runvil Meta-Framework (module architecture).
- Go stdlib `net/http`, `os/signal`.