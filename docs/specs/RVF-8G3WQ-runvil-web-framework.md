# Specification — Runvil Web Framework

| Field       | Value                                       |
| ----------- | ------------------------------------------- |
| SpecID      | RVF-8G3WQ                                   |
| Title       | Runvil Web Framework                       |
| Status      | Draft                                       |
| Version     | 0.1.0                                       |
| Date        | 2026-08-18                                  |
| Author      | Runvil Contributors                         |
| Domain      | Frameworks — Web                            |

## 1. Context

The `web` package is the Runvil web framework, hosted inside the framework
monorepo (`github.com/runvil/framework/web`). It composes the Go standard
library's `net/http` and `html/template` into a routing and rendering layer
used for both live serving and static site generation.

It is the shared foundation for ecosystem site builders such as `mdbind`:
pages are defined as routes, rendered through templates, and either served over
HTTP or exported to a directory as a complete static website.

## 2. Problem Statement

Web work in Go sits between two extremes: raw `net/http` boilerplate (manual
method/path matching, parameter extraction, static file handling, template
plumbing) or heavy third-party web frameworks that couple applications to a
single vendor and duplicate the standard library. The Runvil ecosystem lacks a
web foundation of its own that dogfoods `net/http` and `html/template` and
unifies live serving with static export.

The result: ecosystem tools hand-roll routing and page generation differently,
and site builders cannot share a common static-export primitive.

## 3. Goals

- G1 — Provide method-aware routing with `{param}` path segments.
- G2 — Provide template loading and execution over `html/template`.
- G3 — Allow any page tree to be exported as a deterministic static website.
- G4 — Build only on the standard library (and permitted `runvil-libs` packages).
- G5 — Fit within the framework monorepo as an additive package.

## 4. Non-Goals

- NG1 — A full HTTP client, websockets, or streaming protocols in this phase.
- NG2 — A middleware ecosystem or dependency-injection container.
- NG3 — Server-side session, auth, or templating-specific template DSLs beyond `html/template`.

## 5. Requirements

### 5.1 Routing & Serving

| ID         | Requirement                                                       | Priority |
| ---------- | ----------------------------------------------------------------- | -------- |
| FRK-WEB-001 | Register routes by HTTP method and path pattern.                  | Must     |
| FRK-WEB-002 | Support `{param}` segments extracted into route parameters.       | Must     |
| FRK-WEB-003 | Dispatch requests through an `http.Handler` (`ServeHTTP`).        | Must     |
| FRK-WEB-004 | Return `http.NotFound` for unmatched routes.                      | Must     |
| FRK-WEB-005 | Serve a directory of files under a URL prefix.                    | Must     |

### 5.2 Rendering

| ID         | Requirement                                                       | Priority |
| ---------- | ----------------------------------------------------------------- | -------- |
| FRK-WEB-006 | Load named templates from an `fs.FS` (embed-compatible).          | Must     |
| FRK-WEB-007 | Execute a named template into any `io.Writer`.                    | Must     |
| FRK-WEB-008 | Escape unsafe HTML by default per `html/template`.                | Must     |

### 5.3 Static Export

| ID         | Requirement                                                       | Priority |
| ---------- | ----------------------------------------------------------------- | -------- |
| FRK-WEB-009 | Write each page as `<outDir>/<path>/index.html`.                  | Must     |
| FRK-WEB-010 | Write the root page as `<outDir>/index.html`.                     | Must     |
| FRK-WEB-011 | Emit deterministic ordering and identical output for identical input. | Must |
| FRK-WEB-012 | Sanitize output paths to prevent directory traversal.             | Must     |
| FRK-WEB-013 | Copy asset bytes to their declared URL (e.g. `assets/style.css`). | Must     |

## 6. Non-Functional Requirements

- NFR1 — **Memory safety.** The `unsafe` package must not be used.
- NFR2 — **Performance.** Routing must stay in the microseconds; export must not buffer unboundedly.
- NFR3 — **Portability.** Linux, macOS, and Windows.
- NFR4 — **Quality gates.** `gofmt`, `go vet ./...`, and `go test ./...` must pass.

## 7. Success Criteria

- S1 — Parameterized routes match and dispatch correctly, covered by tests.
- S2 — `Export` produces a root `index.html`, nested pages, and assets with safe paths.
- S3 — A site builder (`mdbind`) serves and exports the same page tree through this package.
- S4 — The package depends only on the standard library (`net/http`, `html/template`, `io`, `fs`).

## 8. Related Specifications

| SpecID    | Title                                           |
| --------- | ----------------------------------------------- |
| [RVF-5K7PZ](./RVF-5K7PZ-web-theming-system.md)        | Web Theming System               |
| [RVF-PN41Q](./RVF-PN41Q-static-site-generator.md)     | Static Site Generator            |
| [RVF-M8SSR](./RVF-M8SSR-cli-application-model.md) | CLI Application Model                |
| [RVF-WXQQ5](./RVF-WXQQ5-cli-output-formatting.md) | CLI Output & Formatting              |
| [RVF-QOFJK](./RVF-QOFJK-runvil-meta-framework.md) | Runvil Meta-Framework                |

## 9. References

- [RVF-M8SSR](./RVF-M8SSR-cli-application-model.md) — CLI Application Model (framework entry patterns).
- [RVL-4Y8UP](https://github.com/runvil/libs/blob/main/docs/specs/RVL-4Y8UP-runvil-libraries.md) — Runvil Libraries initial specification.
- [RVL-N459G](https://github.com/runvil/libs/blob/main/docs/specs/RVL-N459G-terminal-io-rendering.md) — Terminal I/O & Rendering.