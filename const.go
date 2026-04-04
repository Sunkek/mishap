package mishap

// Code is a string error code that can be used with errors.Is.
//
// Code intentionally implements the error interface so that
// errors.Is(err, mishap.CodeNotFound) works without a sentinel variable.
// Do NOT return a bare Code as an error from your own functions — always
// use [New] or [Wrap] to produce a proper *Err value.
type Code string

func (c Code) Error() string { return string(c) }

const (
	// CodeInternal is the catch-all for unexpected server-side failures.
	CodeInternal = Code("INTERNAL_ERROR")

	// ── Client errors (4xx) ─────────────────────────────────────────────────

	CodeBadRequest   = Code("BAD_REQUEST")
	CodeValidation   = Code("VALIDATION_ERROR")
	CodeUnauthorized = Code("UNAUTHORIZED") // 401 — not authenticated
	CodeForbidden    = Code("FORBIDDEN")    // 403 — authenticated but not permitted
	CodeNotFound     = Code("NOT_FOUND")
	CodeConflict     = Code("CONFLICT")
	CodeGone         = Code("GONE")

	CodeMethodNotAllowed      = Code("METHOD_NOT_ALLOWED")
	CodeRequestTimeout        = Code("REQUEST_TIMEOUT")
	CodeRequestEntityTooLarge = Code("REQUEST_ENTITY_TOO_LARGE")
	CodeUnsupportedMediaType  = Code("UNSUPPORTED_MEDIA_TYPE")
	CodeUnprocessableEntity   = Code("UNPROCESSABLE_ENTITY")
	CodeTooManyRequests       = Code("TOO_MANY_REQUESTS")

	// ── Server errors (5xx) ─────────────────────────────────────────────────

	CodeNotImplemented     = Code("NOT_IMPLEMENTED")
	CodeServiceUnavailable = Code("SERVICE_UNAVAILABLE")
	CodeGatewayTimeout     = Code("GATEWAY_TIMEOUT")

	// ── General / transport-layer codes ─────────────────────────────────────

	CodeUnknown           = Code("UNKNOWN")            // catch-all for unclassified errors
	CodeCancelled         = Code("CANCELLED")          // context cancelled
	CodeDeadlineExceeded  = Code("DEADLINE_EXCEEDED")  // deadline / timeout
	CodeResourceExhausted = Code("RESOURCE_EXHAUSTED") // quota, memory, etc.
	CodeAborted           = Code("ABORTED")            // conflict / retry
	CodeDataLoss          = Code("DATA_LOSS")          // unrecoverable corruption
)
