# Specification — Runvil Static Site Generator

| Field       | Value                                       |
| ----------- | ------------------------------------------- |
| SpecID      | RVF-PN41Q                                   |
| Title       | Runvil Static Site Generator                |
| Status      | Draft                                       |
| Version     | 0.1.0                                       |
| Date        | 2026-08-18                                  |
| Author      | Runvil Contributors                         |
| Domain      | Frameworks — Web — SSG                      |

## 1. Context

The `web` package provides a low-level routing, templating, and static-export
layer, while `mdbind` generates sites from markdown. Neither offers the
component-based authoring model popularized by Astro and SvelteKit: reusable
components, layouts that wrap page content, and styles scoped to the component
that declares them.

This specification defines an Astro/Svelte-inspired static site generator
(`web/ssg`) that composes `html/template` components and layouts into a static
website with a single build step.

## 2. Problem Statement

Sites built with plain `html/template` repeat navigation, headers, and footers
in every page. Global stylesheets grow unbounded and leak between pages, and
there is no shared notion of a "component" with its own markup and styles.
Every consumer reinvents the same composition and asset-collection logic.

## 3. Goals

- G1 — Model a site as a tree of named components composed through templates.
- G2 — Scope a component's CSS to that component's rendered root element.
- G3 — Support layouts that wrap page content in a shared shell.
- G4 — Produce a complete static site (HTML plus collected CSS) in one build.
- G5 — Serve the same pages during development through a plain HTTP handler.

## 4. Non-Goals

- NG1 — Markdown authoring, collections, or content pipelines (owned by mdbind).
- NG2 — Client-side hydration, islands, or JavaScript frameworks.
- NG3 — Runtime data fetching; all data is provided at build time.

## 5. Requirements

| ID          | Requirement                                                       | Priority |
| ----------- | ----------------------------------------------------------------- | -------- |
| FRK-SSG-001 | Provide a `Site` builder with named `Component`s (template + scoped style). | Must |
| FRK-SSG-002 | Components compose by invoking other named components within templates. | Must |
| FRK-SSG-003 | Each component's rendered HTML carries a unique scope attribute on its root element, injected automatically. | Must |
| FRK-SSG-004 | Each component's CSS is rewritten so every rule applies only under its scope, with `:root` rules left global. | Must |
| FRK-SSG-005 | Provide `Layout`s with a content slot wrapping page output.         | Must |
| FRK-SSG-006 | Pages declare an output path, title, layout, root component, and data. | Must |
| FRK-SSG-007 | `Site.Build` writes pages as `index.html` (clean URLs) plus assets and a single collected `style.css`. | Must |
| FRK-SSG-008 | Only CSS of components and layouts actually used by built pages is emitted. | Must |
| FRK-SSG-009 | `Site` serves the same pages and assets over HTTP for development.  | Must |

## 6. Non-Functional Requirements

- NFR1 — **Memory safety.** The `unsafe` package must not be used.
- NFR2 — **Minimal dependencies.** HTML parsing uses `golang.org/x/net/html`; no JavaScript or CSS tooling.
- NFR3 — **Quality gates.** `gofmt`, `go vet ./...`, and `go test ./...` must pass.

## 7. Success Criteria

- S1 — A page composed of nested components builds to static HTML where component styles are prefixed by their scope.
- S2 — Building a site with unused components omits their CSS from `style.css`.
- S3 — The development handler renders the same output as the build.

## 8. Related Specifications

| SpecID    | Title                                           |
| --------- | ----------------------------------------------- |
| [RVF-8G3WQ](./RVF-8G3WQ-runvil-web-framework.md) | Runvil Web Framework             |
| [RVF-5K7PZ](./RVF-5K7PZ-web-theming-system.md) | Runvil Web Theming System       |

## 9. References

- [RVM-5F9TL](https://github.com/runvil/mdbind/blob/main/docs/specs/RVM-5F9TL-mdbind-site-builder.md) — mdbind Site Builder (markdown-driven counterpart).