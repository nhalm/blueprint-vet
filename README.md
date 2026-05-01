# blueprint-vet

[![Go Reference](https://pkg.go.dev/badge/github.com/nhalm/blueprint-vet.svg)](https://pkg.go.dev/github.com/nhalm/blueprint-vet)
[![Go Report Card](https://goreportcard.com/badge/github.com/nhalm/blueprint-vet)](https://goreportcard.com/report/github.com/nhalm/blueprint-vet)
[![CI](https://github.com/nhalm/blueprint-vet/actions/workflows/ci.yml/badge.svg)](https://github.com/nhalm/blueprint-vet/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/nhalm/blueprint-vet)](https://github.com/nhalm/blueprint-vet/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/nhalm/blueprint-vet)](go.mod)
[![License](https://img.shields.io/github/license/nhalm/blueprint-vet)](LICENSE)

Static analysis for services built from [go-blueprint](https://github.com/nhalm/go-blueprint) patterns.

A multichecker binary plus a SQL file checker that enforce blueprint conformance rules — handler response writing, repository executor routing, model ID typing, soft-delete defaults, and other patterns documented in the blueprint that compile and pass tests but produce silent-but-wrong behavior.

## Install

```sh
go install github.com/nhalm/blueprint-vet/cmd/blueprint-vet@latest
go install github.com/nhalm/blueprint-vet/cmd/blueprint-sql-check@latest
```

## Use

```sh
blueprint-vet ./...
blueprint-sql-check ./internal/repository/queries
```

### As a golangci-lint plugin

The Go AST analyzers can also run inside `golangci-lint` via the [Module
Plugin System](https://golangci-lint.run/docs/plugins/module-plugins/), which
collapses two Go-AST tools into a single `golangci-lint run` invocation.

Register the plugin in a `.custom-gcl.yml` next to your `.golangci.yml`:

```yaml
version: v2.11.4
plugins:
  - module: 'github.com/nhalm/blueprint-vet'
    import: 'github.com/nhalm/blueprint-vet/plugin'
    version: latest   # or pin to a specific release tag
```

Then enable it in `.golangci.yml`:

```yaml
version: "2"
linters:
  enable:
    - blueprint-vet
  settings:
    custom:
      blueprint-vet:
        type: module
        description: Blueprint conformance rules (R-1..R-7, R-11, R-12).
```

Build the custom binary once and run as usual:

```sh
golangci-lint custom        # produces ./custom-gcl
./custom-gcl run ./...
```

`blueprint-sql-check` is a SQL-file linter and stays a standalone binary —
golangci-lint only lints Go source.

## Rules

### Go AST analyzers (`blueprint-vet`)

| ID  | Name              | What it catches |
|-----|-------------------|-----------------|
| R-1 | `nowriteheader`   | Direct `WriteHeader` on `http.ResponseWriter` in `internal/api` (use `chikit.SetResponse`/`chikit.SetError`). |
| R-2 | `nojsonencode`    | `json.NewEncoder(w).Encode(...)` against the response writer in `internal/api`. |
| R-3 | `handlerjsonbool` | `chikit.JSON` / `chikit.Query` calls whose bool return is discarded. |
| R-4 | `mockgendirective`| `*_interface.go` files missing a `//go:generate mockgen` line. |
| R-5 | `repoexecutor`    | Repository methods passing `r.db` directly to generated calls instead of `executorFromContext(ctx, r.db)`. |
| R-6 | `idtypeuuid`      | Model ID fields under `internal/models` typed as something other than `uuid.UUID`. |
| R-7 | `apperroralias`   | `internal/errors` imported without the `apperrors` alias. |
| R-8 | `nofmtprint`      | `fmt.Print*` / `fmt.Fprint*` outside `cmd/` and `internal/config/`. (`fmt.Sprint*` is allowed — it returns a string and performs no I/O.) |
| R-11| `layerdirection`  | `internal/models` importing repository/service/api; `internal/api` importing `internal/repository` directly. |
| R-12| `errortranslate`  | Repository methods returning bare `err` from generated calls instead of wrapping through `translateError`. |

### SQL file rules (`blueprint-sql-check`)

| ID  | Name              | What it catches |
|-----|-------------------|-----------------|
| R-9 | `softdelete`      | `:one`/`:many`/`:paginated` SELECT queries missing `deleted_at` filter (opt-out: `*IncludingDeleted`, `*Audit`, `*Trash`, `*AllVersions`). |
| R-10| `paginatedorderby`| `:paginated` queries missing `ORDER BY`. |

Full statements and rationale: [`proposals/blueprint-vet.md`](https://github.com/nhalm/go-blueprint/blob/main/proposals/blueprint-vet.md).

## Develop

```sh
make tidy   # resolve module dependencies
make test   # run analyzer tests (golden testdata)
make build  # compile binaries
```

Each analyzer lives in its own package under `analysis/<name>/` with a sibling `_test.go` and `testdata/` directory. Tests use `analysistest.Run` against `// want "..."` annotations in testdata source. The golangci-lint Module Plugin entry point lives in `plugin/` and re-exports the same analyzer set used by `cmd/blueprint-vet`.
