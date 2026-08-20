# Specification — App Directory Structure

| Field       | Value                                       |
| ----------- | ------------------------------------------- |
| SpecID      | RVF-CCI0N                              |
| Title       | App Project Directory Structure Standard    |
| Status      | Draft                                       |
| Date        | 2026-08-19                                  |
| Author      | Runvil Contributors                         |
| Domain      | Framework — Project Structure              |

## 1. Context

Runvil applications are fullstack monoliths: one binary serves server-rendered
pages (RVF-0F2EB), JSON APIs (RVF-230KF), static sites (RVF-PT8OD), and
embedded assets composed through `web.App` (RVF-C4087). Modern fullstack
frameworks define a canonical project layout — Laravel's `app/`, `config/`,
`resources/`, Astro/SvelteKit's `src/` — so new projects start conventional and
tooling behaves predictably. Runvil defines the same kind of standard for Go,
then keeps it a *convention* rather than a hard rule.

Today `runvil new` emits only a bare CLI skeleton and earlier example projects
predate any layout standard. Without a canonical structure every app invents
its own arrangement, divergence grows, and `runvil run`/`dev`/`init`
(RVN-6K41E, RVN-RD3WS, RVN-1QGI2) lose the predictability they depend on.

## 2. Problem Statement

There is no defined, *ideal* directory structure for a Runvil app. Teams copy
idiosyncratic layouts between projects; the devtool cannot offer a
conventional default; and newcomers have no canonical answer for "where do the
templates, styles, handlers, and entry point live?". A standard that is
idiomatic Go while matching the expectations of modern frameworks closes that
gap — provided it stays a default users may reorganize, never a prison.

## 3. Goals

- G1 — Define a canonical app directory structure that is idiomatically Go and
  conventional for modern fullstack frameworks.
- G2 — Scaffold it deterministically with `runvil new` (app) and `runvil init`
  (site/book).
- G3 — Keep the structure a convention: users may relocate directories when
  they declare the move explicitly; no tool hardcodes magic paths.
- G4 — Keep every project single-binary friendly: entry point, internal
  packages, UI sources, and static assets all embed via stdlib `embed`.
- G5 — Preserve existing shape detection (app/site/book) and the static
  book/site pipeline (RVN-1QGI2) unchanged.

## 4. Non-Goals

- NG1 — Enforcing the layout across all projects; it is a recommended default
  and scaffold output, not a linter or a hard gate.
- NG2 — Codifying multi-process or microservice layouts; this standard covers
  the single-process monolith.
- NG3 — Reflection-driven or auto-discovery assembly; paths are declared, not
  guessed beyond stable markers.
- NG4 — Client framework specifics; the JS/TS bridge (RVF-F2TQC) covers how
  compiled frontends attach, not where they must live.

## 5. Requirements

### 5.1 Canonical App Layout

A `runvil new <name>` app project follows this layout by default:

```
<name>/
├── cmd/
│   └── <name>/
│       └── main.go          # thin entry: load config, build app, run
├── internal/
│   ├── app/                 # assembly: New(deps) *web.App (all wiring here)
│   ├── config/              # typed config: struct + Load() + Validate()
│   ├── http/                # handlers, routes, middleware (API + SSR)
│   └── …                    # domain/store/service packages (choice)
├── web/
│   ├── components/          # reusable components (templates/registry entries)
│   ├── layouts/             # page shells
│   ├── pages/               # SSR page definitions (data + mounts)
│   └── theme/               # palette, tokens, toggle styles
├── public/                  # static assets served as-is (embedded)
├── manuscript/              # optional book source (docs/marketing)
├── runvil.yaml              # project config at the repo root
├── go.mod
├── go.sum
├── README.md
└── .gitignore
```

| ID          | Requirement                                                       | Priority |
| ----------- | ----------------------------------------------------------------- | -------- |
| FRK-STR-001 | App projects default to the canonical layout: `cmd/<name>`, `internal/`, `web/`, `public/`, root `runvil.yaml`. | Must |
| FRK-STR-002 | The entry point is `cmd/<name>/main.go` (a root `main.go` is allowed for single-command apps); it must stay thin — load config, build the app, call `App.Run` — never host business logic. | Must |
| FRK-STR-003 | App logic lives in `internal/` packages by default; anything importable by the outside world belongs under `pkg/` and is the project's choice. | Must |
| FRK-STR-004 | Full assembly composes in `internal/app` (e.g. `app.New(cfg, store) *web.App`); `cmd` only calls it. | Should |
| FRK-STR-005 | Configuration is a typed struct in `internal/config`, loaded through `libs/config` (RVL-2X1QZ) and validated through `libs/validate` (RVL-LHANF); the config file is `runvil.yaml` at the repo root. | Must |
| FRK-STR-006 | Frontend sources (components, layouts, pages, theme, styles) live under `web/` by default and register into the app at assembly time; built-in components come from the `ui` registry (RVF-PPUWX). | Must |
| FRK-STR-007 | Static assets served as-is live under `public/` and are embedded via `//go:embed public/...` into the binary. | Must |
| FRK-STR-008 | An optional book source is `manuscript/` (RVN-1QGI2); an optional exported SSG site is `site/`. Apps may serve either through `App.Static`. | Should |
| FRK-STR-009 | The layout is a convention, not a contract: no framework or devtool code may hardcode these paths beyond the stable markers in FRK-STR-030 through FRK-STR-032. | Must |
| FRK-STR-010 | A project may relocate any directory by declaring the move centrally in `runvil.yaml` (e.g. `project.dirs.web: src/web`); un-declared relocations are the author's responsibility. | Should |

### 5.2 Scaffolding

| ID          | Requirement                                                       | Priority |
| ----------- | ----------------------------------------------------------------- | -------- |
| FRK-STR-020 | `runvil new <name>` scaffolds the canonical app layout (FRK-STR-001) with a runnable hello page, one JSON endpoint, and `runvil.yaml`; it builds and runs without modification. | Must |
| FRK-STR-021 | `runvil init` scaffolds the site/book layout: `manuscript/` plus `runvil.yaml` with shared fields (title, input, output, theme, optional `ssg:`). | Must |
| FRK-STR-022 | Scaffolded projects must work with `runvil run`, `runvil dev`, `runvil build`, and `runvil serve` exactly as documented (RVN-6K41E). | Must |
| FRK-STR-023 | Scaffolding keeps its safety contract: refuse non-empty targets and invalid names with exit code 2, list created files deterministically (RVN-RD3WS). | Must |
| FRK-STR-024 | Scaffolded output pins the framework/libs versions the devtool was built with (RVN-RD3WS RND-SC-006). | Must |

### 5.3 Discovery & Tooling

| ID          | Requirement                                                       | Priority |
| ----------- | ----------------------------------------------------------------- | -------- |
| FRK-STR-030 | Shape detection (RVN-6K41E RND-SHD-001..004) stays marker-based: `main.go`/`cmd/*` with `func main` ⇒ app; `runvil.yaml#ssg:` or `ssg.yaml` ⇒ site; `manuscript/` ⇒ book. It does not probe the canonical layout. | Must |
| FRK-STR-031 | The `runvil dev` watched set derives from the canonical layout — `**/*.go`, `web/**`, `public/**`, `runvil.yaml`, `ssg.yaml`, `manuscript/**` — with a `runvil.yaml` override for relocated dirs. | Should |
| FRK-STR-032 | `runvil info` reports the detected kind and, when resolved, the effective directories for `cmd`, `internal`, `web`, `public`, and `manuscript`. | Should |
| FRK-STR-033 | Applications resolve layout directories only through the declared locations (defaults or `runvil.yaml` overrides); nothing reflects on the filesystem beyond markers and declared paths. | Must |

## 6. Non-Functional Requirements

- NFR1 — **Memory safety.** The `unsafe` package must not be used.
- NFR2 — **Dependencies.** The scaffolding lives in `runvil` (stdlib + `libs/core`); framework packages may use stdlib + `libs`.
- NFR3 — **Portability.** Linux, macOS, Windows.
- NFR4 — **Quality gates.** `gofmt`, `go vet ./...`, `go test ./...` pass in every affected repo.
- NFR5 — **Determinism.** Identical input always yields the identical layout; no time-, locale-, or filesystem-order-dependent output.
- NFR6 — **Embeddability.** The canonical layout embeds through stdlib `embed` only; no generated code, no external bundler required.

## 7. Success Criteria

- S1 — `runvil new demo` produces a project where `runvil run` serves a page and a JSON endpoint, and `runvil dev` restarts after a `web/` or `.go` edit.
- S2 — `runvil init` produces a `manuscript/` + `runvil.yaml` that `runvil build` and `runvil serve` accept unchanged.
- S3 — Declaring `project.dirs.web: src/web` in `runvil.yaml` lets `runvil dev` watch the relocated directory without code changes.
- S4 — `runvil info` reports `Project kind: app` and the effective directories.
- S5 — No existing static/book project changes behavior.

## 8. Related Specifications

| SpecID      | RVF-CCI0N                              |
| --------- | ----------------------------------------------- |
| [RVF-C4087](./RVF-APP-C4087-runvil-app-framework.md) | Runvil App Framework (assembly) |
| [RVF-C9WLJ](./RVF-DI-C9WLJ-app-container-service-providers.md) | App Container & Service Providers |
| [RVF-PPUWX](./RVF-UI-PPUWX-layout-ui-system.md) | Layout & UI System |
| [RVF-PT8OD](./RVF-SSG-PT8OD-static-site-generator.md) | Static Site Generator |
| [RVF-0F2EB](./RVF-WEB-0F2EB-server-frontend-pipeline.md) | Server & Frontend Rendering Pipeline |
| [RVF-230KF](./RVF-HTTP-230KF-http-api-pipeline.md) | HTTP & API Pipeline |
| [RVF-F2TQC](./RVF-JS-F2TQC-js-ts-framework-integration.md) | JS/TS Framework Integration (client assets) |
| [RVN-6K41E](https://github.com/runvil/runvil/blob/main/docs/specs/RVN-RUN-6K41E-runvil-run-dev-deploy.md) | runvil run/dev/deploy |
| [RVN-RD3WS](https://github.com/runvil/runvil/blob/main/docs/specs/RVN-SCAFFOLD-RD3WS-project-scaffolding.md) | Project Scaffolding |
| [RVN-1QGI2](https://github.com/runvil/runvil/blob/main/docs/specs/RVN-BUILD-1QGI2-project-building.md) | Project Building |
| [RVL-2X1QZ](https://github.com/runvil/libs/blob/main/docs/specs/RVL-CONFIG-2X1QZ-configuration-loading.md) | Configuration Loading |
| [RVL-LHANF](https://github.com/runvil/libs/blob/main/docs/specs/RVL-VALIDATE-LHANF-struct-validation.md) | Struct Validation |

## 9. References

- [RVF-CMBZJ](./RVF-META-CMBZJ-runvil-meta-framework.md) — Runvil Meta-Framework (module architecture).
- Go standard library `embed`, `internal/`, and `cmd/` project-layout idioms.