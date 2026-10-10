# Go backend roadmap

Course: "Backend Engineering with Go" (Tiago Taquelim). Project: social_blog.
Read only when I ask for a level check or "what's next". Don't load it otherwise.

## How to use (teacher protocol)

Per course section:

1. Before I watch: ask me 2 questions about the topic. Don't explain yet.
2. While I build: hint ladder from CLAUDE.md.
3. After: review my code for the "Done when" items. File:line, why, no rewrite.
4. I tick the box myself, only when I can explain it without looking. You don't edit this file.

Level check: read this file + the repo, then tell me: ready / weak / build next. Max 10 lines.

## Course sections

- [ ] **Architecture**: REST design, transport/service/storage layers
      Done when: I can explain why handlers don't touch SQL.
- [ ] **Advanced Go**: errors, interfaces, pointers, goroutines, context, channels, maps, mutexes
      Done when: I can explain when a data race happens and how to detect it (`go test -race`).
- [ ] **TCP to HTTP**: net, net/http, JSON encode/decode
      Done when: I can explain what a handler, a mux and a middleware are.
- [ ] **Scaffolding**: config, env vars, router, health check, recovery middleware
      Done when: server has timeouts and I know why each one exists.
- [ ] **Databases**: repository pattern, pool config, migrations
      Done when: I can explain pool settings and write an up/down migration.
- [ ] **Posts CRUD**: validation, errors package, comments, PATCH vs PUT, optimistic concurrency, query timeouts
      Done when: I can explain how the version column prevents lost updates.
- [ ] **User feed**: profile, followers, indexes, feed query
      Done when: I can read an `EXPLAIN` and justify each index.
- [ ] **Filtering, sorting, pagination**
      Done when: I can explain offset vs cursor pagination and the tradeoff.
- [ ] **Docs**: swagger
- [ ] **Remaining sections** (course page lists: auth with email activation, middleware, rate limiting, optimisations, CI/CD, deploy to Google Cloud)
      Done when: I can explain how a token is issued and checked, and what rate limiting protects.

## After the course (only if needed, one at a time)

- [ ] Tests against a real Postgres, not only mocks
- [ ] Graceful shutdown + context propagation end to end
- [ ] Structured logging (`log/slog`)
- [ ] Transactions in pgx
- [ ] Read `net/http` and `database/sql` source for one hour each

## Not now

Microservices, Kubernetes, gRPC, frameworks. Finish the monolith well first.
