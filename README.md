# blueprint-vet

Static analysis for services built from [go-blueprint](https://github.com/nhalm/go-blueprint) patterns.

A multichecker binary plus a SQL file checker that enforce blueprint conformance rules — handler response writing, repository executor routing, model ID typing, soft-delete defaults, and other patterns documented in the blueprint that compile and pass tests but produce silent-but-wrong behavior.

> **Status:** early development (`v0.x`). Rule set may shift; pin to a tag once `v1.0.0` ships.

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
| R-8 | `nofmtprint`      | `fmt.Print*` / `fmt.Fprint*` / `fmt.Sprint*` outside `cmd/` and `internal/config/`. |
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

Each analyzer lives in its own package under `internal/analysis/<name>/` with a sibling `_test.go` and `testdata/` directory. Tests use `analysistest.Run` against `// want "..."` annotations in testdata source.
