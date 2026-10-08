# social_blog

## Choosing between `Exec`, `QueryRow`, and `Query`

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
