package mishap

type Code string

func (c Code) Error() string { return string(c) }

const (
	CodeInternal     = Code("INTERNAL_ERROR")
	CodeBadRequest   = Code("BAD_REQUEST")
	CodeValidation   = Code("VALIDATION_ERROR")
	CodeUnauthorized = Code("UNAUTHORIZED")
	CodeForbidden    = Code("FORBIDDEN")
	CodeNotFound     = Code("NOT_FOUND")
	CodeConflict     = Code("CONFLICT")

	// Additional client-side errors (4xx)
	CodeMethodNotAllowed      = Code("METHOD_NOT_ALLOWED")
	CodeRequestTimeout        = Code("REQUEST_TIMEOUT")
	CodeRequestEntityTooLarge = Code("REQUEST_ENTITY_TOO_LARGE")
	CodeUnsupportedMediaType  = Code("UNSUPPORTED_MEDIA_TYPE")
	CodeUnprocessableEntity   = Code("UNPROCESSABLE_ENTITY") // Often used for validation failures beyond basic bad request
	CodeTooManyRequests       = Code("TOO_MANY_REQUESTS")
	CodeGone                  = Code("GONE")

	// Additional server-side errors (5xx)
	CodeNotImplemented     = Code("NOT_IMPLEMENTED")
	CodeServiceUnavailable = Code("SERVICE_UNAVAILABLE")
	CodeGatewayTimeout     = Code("GATEWAY_TIMEOUT")

	// General non-HTTP codes
	CodeUnauthenticated   = Code("UNAUTHENTICATED")
	CodeUnknown           = Code("UNKNOWN")            // Catch-all for unspecified errors
	CodeCancelled         = Code("CANCELLED")          // Operation cancelled (e.g., context cancellation)
	CodeDeadlineExceeded  = Code("DEADLINE_EXCEEDED")  // Timeout or deadline
	CodeResourceExhausted = Code("RESOURCE_EXHAUSTED") // Out of memory, quota exceeded
	CodeAborted           = Code("ABORTED")            // Operation aborted due to conflict or retry
	CodeDataLoss          = Code("DATA_LOSS")          // Unrecoverable data corruption
)
