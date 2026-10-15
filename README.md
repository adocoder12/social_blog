# social_blog

A REST API for a small social blog. It's written in Go with Postgres, and follows the Udemy course _Backend Engineering with Go_.

**Stack:** Go · chi (router) · pgx (Postgres driver) · golang-migrate · go-playground/validator · slog · Docker Compose · Air (live reload)

---

## Running it

```sh
# create .env first (see Config below)
make dev                # docker compose up --build: API + Postgres, with live reload
```

The API listens on `:$PORT` (default `8080`). Migrations run automatically at startup.

| Command                      | What it does                                |
| ---------------------------- | ------------------------------------------- |
| `make build` / `make run`    | Build `bin/app` / run it                    |
| `make migrate-create name=x` | New numbered up/down migration pair         |
| `make migrate-up` / `-down`  | Apply all migrations / roll back one        |
| `make db-shell`              | `psql` inside the Postgres container        |
| `make reset-docker`          | Wipe volumes and rebuild from scratch       |
| `go fmt ./...`               | Format all code (do this before committing) |
| `go vet ./...`               | Catch suspicious code the compiler allows   |

### Config

Read from environment variables (and `.env` through godotenv) in `internal/config`. The main ones are `PORT`, `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`, `DB_MAX_CONNS` and `DB_MIN_CONNS`. Config is validated at startup, and the app exits early if it's invalid.

---

## Project layout

```
cmd/api/main.go          entry point: config → DB pool → migrations → wiring → HTTP server
internal/
  config/                env vars → Config struct (+ Validate)
  db/                    pgx pool + migration runner; SQL files in db/migrate/migrations
  model/                 domain types (Post, User, PostPatch): no HTTP, no SQL
  store/                 repositories: interfaces (store.go) + Postgres implementations
  services/              business logic / use cases; depends on store interfaces
  dto/                   request/response shapes (JSON + validate tags) and model mapping
  handler/               HTTP layer: routes, handlers, JSON/error helpers, Application struct
```

`internal/` means no other module can import these packages. The Go compiler enforces this.

### Request flow

```
HTTP request
  → chi middleware (RequestID, Logger, Recoverer, Timeout 8s)
  → handler        parse URL/body → validate DTO → call service → map error to status
  → service        use case logic (e.g. fetch → apply patch → save)
  → store          SQL through pgx, with r.Context() passed all the way down
  → Postgres
```

Dependencies point **one way only**: handler → service → store. Lower layers never import higher ones.

---

## API

Base path: `/api/v1`

| Method   | Path            | Success | Errors                          |
| -------- | --------------- | ------- | ------------------------------- |
| `GET`    | `/health`       | 200     |                                 |
| `GET`    | `/posts`        | 200     | 500                             |
| `POST`   | `/posts`        | 201     | 400 bad body / validation / user missing |
| `GET`    | `/posts/{postID}` | 200   | 400 bad id · 404                |
| `PATCH`  | `/posts/{postID}` | 200   | 400 bad id / body / validation · 404 |
| `DELETE` | `/posts/{postID}` | 204   | 400 bad id · 404                |

**Status code rule of thumb:** 4xx means the client must fix the request. 5xx means the server must be fixed. Never return 500 for bad input.

---

## Design decisions (and why)

### Layers

- **Handler** only deals with HTTP: it decodes, validates, calls **one** service method, and writes JSON plus a status code. No business logic.
- **Service** is where use cases live. Its method signatures follow the use case, not the database. For example, `UpdatePost(ctx, id, patch)` in the service and `UpdatePost(ctx, *Post)` in the repo are different on purpose.
- **Store (repository)** is persistence only: "save what I'm given".

### Interfaces: accept them, implement them implicitly

- `PostService` holds a `store.PostRepository` interface, not the concrete `*PostStore`. Tests can pass a fake repo, with no Postgres needed.
- Only the type **implementing** an interface (`*PostStore`) has to match its method set. A type that **uses** it (`PostService`) can have any methods it likes.
- Go has no `implements` keyword. A type satisfies an interface just by having the methods.

### DTOs vs models

- `dto` types are the JSON contract with clients. `model` types are what the app works with internally.
- Validation tags go on **DTOs**, because that's where untrusted input enters, and each endpoint can have its own rules.
- Mapping functions: `RequestPost.ToModel()`, `UpdatePost.ToPatch()`, `FromModel()`. The service never imports `dto`.
- Partial updates use **pointer fields** (`*string`, `*[]string`): `nil` means "not sent", which is different from "sent as empty".

### Validation

- There's one `*validator.Validate`, created in `NewApplication` and stored on `Application`. It caches info about each struct type and is safe for concurrent use, so it's created once, not per request.
- Order inside a handler: `readJSON` → `validator.Struct` → service call.
- Business rules ("user must exist", "only the author can edit") belong in the **service**, not in struct tags.

### Errors

- Wrap errors with context and `%w`: `fmt.Errorf("post service: get by id: %w", err)`.
- `%w` keeps the original error in the chain, so the handler can still do `errors.Is(err, store.ErrPostNotFound)` → 404.
- Sentinel errors (`ErrPostNotFound`, `ErrUserNotFound`) live in `store/errors.go`.
- Handle an error **once**: either return it (wrapped) or log it and respond. Don't do both.

### Context

- `r.Context()` goes from handler → service → store → pgx. If the client disconnects or the 8s timeout middleware fires, the SQL query is cancelled too.
- `context.Context` is always the **first parameter** and is never stored in a struct.
- Startup uses `context.WithTimeout(ctx, 5*time.Second)` for the DB connection, so a dead DB fails fast instead of hanging.

### Graceful shutdown

- `signal.NotifyContext` cancels a context on Ctrl+C / SIGTERM. The server then calls `Shutdown`, which lets in-flight requests finish.
- `ListenAndServe` runs in a goroutine, and errors come back over a buffered channel. `http.ErrServerClosed` is the normal result of a shutdown, not a failure.

### Migrations

- These are numbered `up`/`down` SQL pairs in `internal/db/migrate/migrations`. **Never edit a migration that has already run.** Add a new one instead.
- `down` must exactly undo `up`, so rollbacks work.

### Known gap: lost updates

`PostService.UpdatePost` reads, modifies, then writes. Two concurrent edits mean the second silently overwrites the first. The fix is optimistic concurrency: a `version` column plus `UPDATE ... WHERE id = $1 AND version = $2`. It belongs in the service and store, which is why the read and write live together in one service method.

---

## Go idioms to keep in mind

- **gofmt is not optional.** Run `go fmt ./...` or format on save. Plain `gofmt` with no arguments waits for stdin.
- **Initialisms stay uppercase:** `ID`, `URL`, `HTTP`, `JSON`, so `GetPostByID`, not `GetPostById`.
- **MixedCaps, no underscores:** `ToModel`, not `Tomodel` or `to_model`.
- **Exported** = starts with a capital letter. Keep things unexported unless another package needs them.
- **Return early** on errors. Keep the happy path unindented.
- **Nil slices encode to JSON `null`**, and empty slices encode to `[]`. Initialize with `make(..., 0)` if clients expect an array.
- **`:=` needs at least one new variable** on the left. Reusing a parameter name with `:=` reassigns it, which compiles but hurts readability.
- **Don't over-engineer:** add an interface when you have a second implementation or a test that needs it, not before.

---

## pgx: choosing between `Exec`, `QueryRow`, and `Query`

| Method     | Returns                    | Use for                                   |
| ---------- | -------------------------- | ----------------------------------------- |
| `Exec`     | `pgconn.CommandTag, error` | Statements where you don't need data back |
| `QueryRow` | `pgx.Row`                  | Exactly one row expected                  |
| `Query`    | `pgx.Rows, error`          | Zero to many rows                         |

### `Exec`

`INSERT`, `UPDATE`, `DELETE` without `RETURNING`. Use `result.RowsAffected()` to check whether anything changed.

```go
result, err := pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
```

### `QueryRow`

A single row: `SELECT ... WHERE id = $1`, or `INSERT/UPDATE ... RETURNING`. Always call `.Scan()`, which also closes it. A missing row returns `pgx.ErrNoRows`.

```go
err := pool.QueryRow(ctx, `SELECT id FROM users WHERE id = $1`, id).Scan(&u.ID)
```

### `Query`

A list of rows. Always `defer rows.Close()`, loop with `rows.Next()`, and check `rows.Err()` afterward.

```go
rows, err := pool.Query(ctx, `SELECT id FROM users`)
if err != nil {
	return err
}
defer rows.Close()

for rows.Next() {
	// rows.Scan(...)
}
return rows.Err()
```

### Gotchas

- `Exec` with `RETURNING` discards the returned data. Use `QueryRow` instead.
- `QueryRow` on a multi-row result silently returns only the first row.
- `QueryRow` errors surface at `Scan`, so never skip the `Scan` error check.
