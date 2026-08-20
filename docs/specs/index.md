# Specifications Index — Runvil Framework

This directory holds the formal specifications for the Runvil Framework.

Ordered by **impact-to-effort** (high impact, low effort first) and **dependency chain** (foundations first).

| SpecID    | Title                                      | Status | Depends On |
| --------- | ------------------------------------------ | ------ | ---------- |
| [RVF-CMBZJ](./RVF-META-CMBZJ-runvil-meta-framework.md) | Runvil Meta-Framework — Initial Specification | Draft | — |
| [RVF-M07QS](./RVF-WEB-M07QS-runvil-web-framework.md) | Runvil Web Framework | Draft | RVF-CMBZJ |
| [RVF-DF3PL](./RVF-SSG-DF3PL-file-site-pipeline.md) | File-Based Site Pipeline (`.rv` Modules) | Draft | RVF-M07QS, RVF-PT8OD |
| [RVF-PT8OD](./RVF-SSG-PT8OD-static-site-generator.md) | Static Site Generator | Draft | RVF-M07QS |
| [RVF-C4087](./RVF-APP-C4087-runvil-app-framework.md) | Runvil App Framework (assembly) | Draft | RVF-M07QS |
| [RVF-C9WLJ](./RVF-DI-C9WLJ-app-container-service-providers.md) | App Container & Service Providers | Draft | RVF-CMBZJ |
| [RVF-CCI0N](./RVF-STRUCT-CCI0N-app-directory-structure.md) | App Project Directory Structure Standard | Draft | RVF-CMBZJ |
| [RVF-5XJFC](./RVF-CLI-5XJFC-cli-application-model.md) | CLI Application Model | Draft | RVF-CMBZJ |
| [RVF-PZ5JU](./RVF-CLI-PZ5JU-cli-scaffolding.md) | CLI Scaffolding | Draft | RVF-5XJFC |
| [RVF-0F2EB](./RVF-WEB-0F2EB-server-frontend-pipeline.md) | Server & Frontend Rendering Pipeline | Draft | RVF-M07QS, RVF-PT8OD, RVF-C4087 |
| [RVF-D57UK](./RVF-SSG-D57UK-ssg-markdown-content-pipeline.md) | SSG Markdown Content Pipeline & Collections | Draft | RVF-PT8OD |
| [RVF-DR5YU](./RVF-SSG-DR5YU-ssg-asset-pipeline.md) | SSG Asset Pipeline | Draft | RVF-PT8OD |
| [RVF-99A63](./RVF-SSG-99A63-ssg-incremental-builds.md) | SSG Incremental Builds & Dependency Graph | Draft | RVF-PT8OD, RVF-D57UK, RVF-DR5YU |
| [RVF-209JV](./RVF-SSG-209JV-ssg-live-reload-hmr.md) | SSG Live Reload & HMR | Draft | RVF-PT8OD, RVF-99A63 |
| [RVF-5ZHQV](./RVF-ARCH-5ZHQV-modular-monolith-architecture.md) | Modular Monolith Architecture Default | Draft | RVF-C9WLJ, RVF-CCI0N |
| [RVF-0Z671](./RVF-UI-0Z671-runvil-ui-framework.md) | Runvil UI Framework | Draft | RVF-M07QS |
| [RVF-PPUWX](./RVF-UI-PPUWX-layout-ui-system.md) | Layout & UI System | Draft | RVF-0Z671 |
| [RVF-V0TMZ](./RVF-UI-V0TMZ-web-theming-system.md) | UI Theming System | Draft | RVF-0Z671 |
| [RVF-230KF](./RVF-HTTP-230KF-http-api-pipeline.md) | HTTP & API Pipeline | Draft | RVF-M07QS |
| [RVF-F2TQC](./RVF-JS-F2TQC-js-ts-framework-integration.md) | JS/TS Framework Integration — Future Bridge | Draft | RVF-0Z671, RVF-PT8OD |
| [RVF-5XJFC](./RVF-CLI-5XJFC-cli-application-model.md) | CLI Application Model | Draft | RVF-CMBZJ |
| [RVF-MFA0T](./RVF-CLI-MFA0T-cli-help-usage.md) | CLI Help & Usage | Draft | RVF-5XJFC |
| [RVF-NPFSE](./RVF-CLI-NPFSE-cli-output-formatting.md) | CLI Output & Formatting | Draft | RVF-5XJFC |
| [RVF-KAKQL](./RVF-CLI-KAKQL-cli-errors-diagnostics.md) | CLI Errors & Diagnostics | Draft | RVF-5XJFC |
| [RVF-FGNZ9](./RVF-CLI-FGNZ9-cli-configuration.md) | CLI Configuration | Draft | RVF-5XJFC |
| [RVF-PZ5JU](./RVF-CLI-PZ5JU-cli-scaffolding.md) | CLI Scaffolding | Draft | RVF-5XJFC |

## Conventions

- Each specification is stored as a single Markdown file named `{ProjectId}-{Scope}-{SpecID}-{slug}.md`.
- SpecIDs are unique, random 5-character alphanumeric codes (e.g., `RVF-CMBZJ`).
- New specifications must be added to this index when created.
- Ordering: impact-to-effort (high impact, low effort first), then dependency chain (foundations first).