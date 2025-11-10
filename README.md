# Oolio — Backend Challenge

Simple Go backend for the Oolio coding challenge.

Live deployment: https://oolio.fly.dev/api

Contents
- `cmd/httpapi` — HTTP server entrypoint
- `internal/routes` — route wiring and handlers for product and order APIs
- `internal/data` — in-memory product data used by handlers
- `internal/binding` — JSON binding + validation helper
- `internal/validation` — validator and translator initialization

API Endpoints
- GET /stats — runtime memory stats (returns JSON)
- GET /api/product/ — list all products (returns 201)
- GET /api/product/{productId} — find product by id (200 or 404)
- POST /api/order/ — create an order (200 on success, 422 on validation errors)

Quick start (local)

Requires Go (see `go.mod`). From the project root:

1. Run the server:

```bash
go run ./cmd/httpapi
```

2. The server listens on the default address from the `cmd/httpapi` code (check `main.go`).

Environment and promocode files

The server reads a comma-separated list of promocode source paths from the environment variable `PROMO_FILES`.
You can provide these in a `.env` file or using your environment. An example `.env.example` is included in the repo. Example:

```env
PROMO_FILES=/path/to/couponbase1,/path/to/couponbase2,/path/to/couponbase3
```

Each path should point to either:
- a plain text file containing one promocode per line (the project will build Pebble DBs from these text files), or
- a previously-created Pebble DB directory produced by the project (the code detects existing DB manifests and re-uses them).

How the Pebble index and promocode validation work

1. Pebble index initialization

   - On startup (`cmd/httpapi/main.go`), the server reads `PROMO_FILES` and calls `index.NewPebbleIndex(paths)`.
   - For each path the app will call `EnsurePebble(path)`:
     - If a Pebble DB already exists for that path, it opens and re-uses it.
     - Otherwise it creates a Pebble DB and bulk-loads the promocodes from the text file into the DB.
   - The returned `PebbleIndex` contains a slice of `PebbleStore` entries, one per provided path.

2. Validation rule registration

   - After the index is created the app registers a custom validator named `promocode` via `internal/validation.RegisterPromocodeValidation`.
   - The `promocode` validator implementation calls `PebbleIndex.IsValid2of3(code)`.

3. How `IsValid2of3` validates a code

   - `IsValid2of3` checks the code against each `PebbleStore` using `PebbleStore.Has(code)`.
   - It counts hits across stores and returns `true` when at least 2 stores contain the code (2-of-3 strategy).
   - This gives robustness if some promo files overlap or are noisy — the code must be present in at least two sources to be considered valid.

4. Usage in requests

   - The `order.OrderRequest` struct has `CouponCode string `json:"couponCode" validate:"omitempty,promocode"``.
   - When a request contains a `couponCode`, the validator will call the registered `promocode` validator. If the code isn't valid per the PebbleIndex, validation fails and the request returns a validation error (422).

Notes and troubleshooting

- If you change `PROMO_FILES`, the server will attempt to build or open pebble DBs for the new paths on startup. Make sure the process has read/write permissions to the target directories.
- For local development you can point `PROMO_FILES` to small test files that contain a few promocodes (one per line) to exercise validation without heavy data.
- The code also contains a Bloom-filter based validator implementation (`internal/index/bloom.go`) used in tests and experiments; production mode uses Pebble for durability and faster exact lookups.

Design note

Initially the project experimented with an in-memory Bloom filter and static hash tables to validate promocodes. The Bloom filter approach required on the order of hundreds of megabytes (>= ~300MB) of RAM for realistic promo datasets, which made it unsuitable for constrained environments. To reduce memory usage and provide durable, on-disk indexes the project uses Pebble DB to store promocodes and performs fast lookups against multiple pebble stores (the `2-of-3` strategy). Pebble reduces memory pressure while keeping lookups performant.

Running tests

Run the full test suite:

```bash
go test ./...
```

What the tests cover
- Router tests exercise the public API surface using an in-memory HTTP server.
- Binding tests validate JSON binding behavior, unknown-field rejection, and validation error formatting.

Notes
- The project uses go-playground/validator for request validation. Some tests initialize the validator and translations to avoid global init dependencies.
- Product data is stored in memory in `internal/data/product.go` for simplicity.

Deployment

The service is deployed to Fly.io at: https://oolio.fly.dev/api

Improvements & Roadmap

This project is a compact backend with a clear starting point. Suggested areas to improve and scope for work:

- Logging
   - Replace ad-hoc log.Printf with a structured logger (zap, zerolog). Add levels (debug/info/warn/error) and request-scoped fields (request id, remote addr, latency).
   - Add centralized log formatting and optional JSON output for ingestion by logging systems.

- Documentation / OpenAPI
   - Add OpenAPI (Swagger) spec describing endpoints, request/response schemas and validation rules. Generate server/client stubs as needed.
   - Serve the OpenAPI UI (Swagger UI / Redoc) at `/docs` for easy exploration.

- Tests
   - Expand tests with table-driven tests for edge cases (invalid bodies, large requests, concurrency).
   - Add integration tests that start the server and use real pebble DBs generated from small promo files.

- Observability
   - Integrate Prometheus metrics (endpoint counts, request latencies, pebble index metrics). Expose `/metrics`.
   - Add distributed tracing hooks (OpenTelemetry) to trace request flow and pebble lookups.

- Validation & Promo index
   - Move validator & index initialization into a testable initializer function so tests can supply mock indexes.
   - Add an optional mode to pre-warm pebble DBs and report load progress.

- Performance & resource usage
   - Benchmark promo lookups and tune pebble options for read-heavy workloads.
   - Add a memory/profile mode to generate pprof profiles for heap and CPU.

- Security
   - Validate and sanitize all inputs more strictly. Add rate-limiting and basic auth or API keys for order endpoints if needed.
   - Hardening for running on shared hosts (drop privileges, limit file descriptors).

- Developer experience
   - Provide a small script to generate sample promo files and pre-build pebble DBs for local development.
   - Add Makefile or task runner for common tasks: test, lint, build, run, generate-swagger.

These improvements can be prioritized based on production needs. If you'd like I can implement one of the above (OpenAPI generation, structured logging, or a local pebble generator script) next.