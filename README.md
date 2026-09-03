# mishap

[![CI](https://github.com/sunkek/mishap/actions/workflows/ci.yml/badge.svg)](https://github.com/sunkek/mishap/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/sunkek/mishap.svg)](https://pkg.go.dev/github.com/sunkek/mishap)
[![Go Report Card](https://goreportcard.com/badge/github.com/sunkek/mishap)](https://goreportcard.com/report/github.com/sunkek/mishap)

A small, explicit error handling package for Go — zero external dependencies.

`mishap` gives every error a typed code so you can route, match, and log errors
without parsing strings. It wraps any `error` value, inherits codes through the
chain, and integrates cleanly with the standard `errors` package.

```
go get github.com/sunkek/mishap
```

---

## Core concept

Every error carries a **Code** — a typed string you can match with `errors.Is`:

```go
var ErrCodeOverheat = mishap.Code("OVERHEAT")

err := mishap.New("temperature over 9000°C", ErrCodeOverheat)

if errors.Is(err, ErrCodeOverheat) {
    // sound the alarms
}
```

Codes survive wrapping, including through `fmt.Errorf("%w", ...)`:

```go
inner := mishap.New("row not found", mishap.CodeNotFound)
outer := mishap.Wrap(inner, "load user") // inherits CodeNotFound

fmt.Println(outer.Code())  // NOT_FOUND
fmt.Println(outer.Error()) // load user: row not found
```

---

## Creating errors

```go
// New requires a non-empty message and code — panics otherwise.
err := mishap.New("user not found", mishap.CodeNotFound)
```

Use the built-in codes for common HTTP/gRPC scenarios, or define your own:

```go
var ErrCodeOverheat = mishap.Code("OVERHEAT")
```

---

## Wrapping errors

`Wrap` wraps any existing error, returning an `error`. Code resolution follows
this precedence:

1. `WithCode` — explicit override
2. The code of the first `*Err` found anywhere in the cause chain
3. `WithDefaultCode` — caller-supplied fallback
4. `CodeInternal`

```go
// Inherit code from the cause chain (most common case)
outer := mishap.Wrap(inner, "load user")

// Force a specific code, ignoring the chain
outer := mishap.Wrap(inner, "load user", mishap.WithCode(mishap.CodeInternal))

// Provide a fallback only when the chain has no code
outer := mishap.Wrap(err, "load user", mishap.WithDefaultCode(mishap.CodeBadRequest))
```

`Wrap` returns a nil `error` when `err` is `nil` — safe to use in one-liners:

```go
return mishap.Wrap(repo.Find(id), "find user")
```

This is why `Wrap` returns `error` rather than `*Err`. A concrete pointer
return would make the line above produce a non-nil `error` interface holding a
nil `*Err`, so the caller's `err != nil` check would take the failure branch on
success — Go's typed-nil trap. Reach the `*Err` with `mishap.As` when you need
its `Code()` or `Message()`:

```go
if mErr, ok := mishap.As(mishap.Wrap(err, "load user")); ok {
    log.Println(mErr.Code())
}
```

---

## Handling errors

Match against the full chain with `errors.Is`:

```go
func handle(err error) {
    switch {
    case errors.Is(err, ErrCodeOverheat):
        // sound the alarms
    case errors.Is(err, mishap.CodeNotFound):
        // return 404
    case errors.Is(err, mishap.CodeBadRequest),
         errors.Is(err, mishap.CodeValidation):
        // return 400
    case errors.Is(err, mishap.CodeUnauthorized):
        // return 401
    case errors.Is(err, mishap.CodeForbidden):
        // return 403
    default:
        // return 500
    }
}
```

Or extract the topmost `*Err` directly with `mishap.As`:

```go
if mErr, ok := mishap.As(err); ok {
    switch mErr.Code() {
    case mishap.CodeNotFound:
        // return 404
    // ...
    }
}
```

`errors.Is` scans the whole chain. `mErr.Code()` returns only the topmost code.
Choose based on whether you care about where in the chain the code appears.

---

## Built-in codes

### Client errors (4xx)

| Code | Value |
|---|---|
| `CodeBadRequest` | `BAD_REQUEST` |
| `CodeValidation` | `VALIDATION_ERROR` |
| `CodeUnauthorized` | `UNAUTHORIZED` |
| `CodeForbidden` | `FORBIDDEN` |
| `CodeNotFound` | `NOT_FOUND` |
| `CodeConflict` | `CONFLICT` |
| `CodeGone` | `GONE` |
| `CodeMethodNotAllowed` | `METHOD_NOT_ALLOWED` |
| `CodeRequestTimeout` | `REQUEST_TIMEOUT` |
| `CodeRequestEntityTooLarge` | `REQUEST_ENTITY_TOO_LARGE` |
| `CodeUnsupportedMediaType` | `UNSUPPORTED_MEDIA_TYPE` |
| `CodeUnprocessableEntity` | `UNPROCESSABLE_ENTITY` |
| `CodeTooManyRequests` | `TOO_MANY_REQUESTS` |

### Server errors (5xx)

| Code | Value |
|---|---|
| `CodeInternal` | `INTERNAL_ERROR` |
| `CodeNotImplemented` | `NOT_IMPLEMENTED` |
| `CodeServiceUnavailable` | `SERVICE_UNAVAILABLE` |
| `CodeGatewayTimeout` | `GATEWAY_TIMEOUT` |

### General / transport-layer

| Code | Value |
|---|---|
| `CodeUnknown` | `UNKNOWN` |
| `CodeCancelled` | `CANCELLED` |
| `CodeDeadlineExceeded` | `DEADLINE_EXCEEDED` |
| `CodeResourceExhausted` | `RESOURCE_EXHAUSTED` |
| `CodeAborted` | `ABORTED` |
| `CodeDataLoss` | `DATA_LOSS` |

---

## API reference

```go
// Create a new structured error. Panics if message or code is empty.
func New(message string, code Code) *Err

// Wrap any error with a message. Returns a nil error if err is nil.
// Panics if message is empty. Returns error, not *Err, so the nil case
// survives assignment — use As to reach the *Err.
func Wrap(err error, message string, opts ...WrapOption) error

// Extract the first *Err from the chain. Convenience wrapper over errors.As.
func As(err error) (*Err, bool)

// WrapOptions
func WithCode(code Code) WrapOption        // force a code
func WithDefaultCode(code Code) WrapOption // fallback code

// *Err methods
func (e *Err) Code() Code    // topmost code
func (e *Err) Message() string // topmost message
func (e *Err) Error() string   // full chain: "outer: inner: cause"
func (e *Err) Unwrap() error   // cause, for errors.Is/As traversal
```

---

## Notes on Code

`Code` implements `error` so that `errors.Is(err, mishap.CodeNotFound)` works
without sentinel variables. **Do not return a bare `Code` as an error from your
own functions** — always use `New` or `Wrap`. Returning a `Code` directly would
compile silently but lose the message and chain.
