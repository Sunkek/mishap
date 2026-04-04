# Changelog

All notable changes to mishap are documented here.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
mishap uses [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased](https://github.com/sunkek/mishap/compare/v1.1.0...HEAD)

---

## [1.1.0](https://github.com/sunkek/mishap/compare/v1.0.0...v1.1.0) — 2026-04-04

### Changed

* **Breaking:** `New` now panics on empty message or code instead of silently
  substituting defaults (`"error"` / `CodeInternal`). Call-site mistakes are
  now caught at development time rather than silently reaching production logs.
* **Breaking:** `Wrap` now panics on empty message for the same reason.
* **Breaking:** `CodeUnauthenticated` removed — it duplicated `CodeUnauthorized`
  (HTTP 401). Use `CodeUnauthorized` for unauthenticated requests.
* `Error()` no longer calls `.Error()` explicitly on the source error; the `%s`
  verb already invokes it. No behaviour change, cleaner internals.

### Added

* `mishap.As(err) (*Err, bool)` — convenience wrapper over `errors.As` for
  extracting the first `*Err` from a chain without a typed variable declaration.
* CI via GitHub Actions (`go vet` + `go test -race` on push and PR).
* `CHANGELOG.md`, `CONTRIBUTING.md`, `SECURITY.md`, `.gitignore`.
* `make test` now runs with `-race`.

### Fixed

* `Code` type now has a doc comment explaining why it implements `error` and
  warning against using it as a direct error return value.
* `(*Err).Is` now has a doc comment clarifying that matching is by code only,
  not message.

---

## [1.0.0](https://github.com/sunkek/mishap/releases/tag/v1.0.0) — 2026-01-08

### Added

* Initial release.
* `Code` type — typed string error code implementing `error` for use with
  `errors.Is`.
* `Err` struct with `message`, `code`, and `sourceErr` fields.
* `New(message, code)` — create a root error.
* `Wrap(err, message, opts...)` — wrap any error, inheriting codes through the
  chain.
* `WithCode` and `WithDefaultCode` wrap options.
* Code inheritance: `WithCode` > chain's first `*Err` code > `WithDefaultCode`
  > `CodeInternal`.
* Built-in codes for common HTTP and gRPC scenarios.
* `errors.Is` / `errors.As` / `errors.Unwrap` compatibility.
* Zero external dependencies.
