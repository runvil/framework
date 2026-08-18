# Runvil Framework

**Runvil** is a meta-framework written in [Go](https://go.dev/). Unlike a conventional framework built within a single, self-contained ecosystem, Runvil composes modules sourced across multiple ecosystems and repositories and orchestrates them into one cohesive, high-performance foundation for building a wide range of applications — from web services and CLI tools to background workers and desktop applications.

## Overview

Runvil is organized as a Go module monorepo that hosts the framework's own orchestration packages. These packages compose and unify components that live outside this repository — such as the libraries maintained in Runvil Libraries — into a single, consistent experience. Adopting only the modules you need remains a first-class concern.

## Features

- **Meta-framework design** — Composes modules from multiple ecosystems into one cohesive framework.
- **Modular architecture** — Individual packages for each concern, composable as needed.
- **Stdlib-first** — Argument parsing (`flag`), logging (`log/slog`), and configuration (`os`) come from the Go standard library.
- **Safe by design** — Go's memory safety with no manual memory management; `unsafe` is not used.
- **Future-ready** — Extensible with additional modules such as HTTP, async runtime integration, configuration, and logging.

## Workspace Layout

| Path          | Description                                                             |
| ------------- | ----------------------------------------------------------------------- |
| `framework/`  | Meta-framework package composing modules sourced across multiple ecosystems. |
| `cli/`        | Integrated command-line application model (initial ecosystem).          |
| `examples/`   | Example applications demonstrating the framework.                       |

Shared primitives and terminal utilities are provided by the [`runvil-libs`](https://github.com/runvil/runvil-libs) monorepo — one of the ecosystems Runvil integrates rather than re-implements.

## Related Repositories

- [runvil-libs](https://github.com/runvil/runvil-libs) — Modular, reusable libraries for the Runvil ecosystem (`core`, `term`).

## Getting Started

### Prerequisites

- Go toolchain 1.22 or newer — see [go.dev/dl](https://go.dev/dl/)

The framework depends on the published `github.com/runvil/runvil-libs` module (fetched via its git URL); no local checkout is required.

### Building

```bash
go build ./...
```

### Testing

```bash
go test ./...
```

### Running the example

```bash
go run ./examples/greet hello --name Alice
GREET_GREETING=Halo go run ./examples/greet hello --name Alice
```

## Roadmap

The following modules are planned for the framework:

- `runvil-async` — Async runtime abstractions and task scheduling.
- `runvil-http` — HTTP server and client abstractions.
- `runvil-router` — Routing layer.
- `runvil-server` — Integrated server (async + http + router).
- `runvil-worker` — Background jobs and scheduling.

## Contributing

Contributions are welcome. Please run `gofmt` and `go vet` before submitting changes, and ensure all tests pass.

## License

Runvil is distributed under the [MIT License](LICENSE).
