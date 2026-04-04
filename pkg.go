package mishap

import (
	"errors"
	"fmt"
)

// Err is a structured error that carries a [Code] and an optional cause chain.
type Err struct {
	sourceErr error
	message   string
	code      Code
}

// Message returns the error's human-readable message (the topmost layer only).
func (e *Err) Message() string { return e.message }

// Code returns the error's code (the topmost layer only).
// To search the entire chain use [errors.Is].
func (e *Err) Code() Code { return e.code }

// Error implements the error interface and prints the full chain.
func (e *Err) Error() string {
	if e.sourceErr != nil {
		return fmt.Sprintf("%s: %s", e.message, e.sourceErr)
	}
	return e.message
}

// Unwrap allows [errors.Is] and [errors.As] to traverse the cause chain.
func (e *Err) Unwrap() error { return e.sourceErr }

// Is matches only on code equality.
//
// When the target is a [Code] or a *[Err], the comparison is by code value
// alone — the message is irrelevant. This means:
//
//	errors.Is(err, mishap.CodeNotFound)            // matches any *Err with that code in the chain
//	errors.Is(err, mishap.New("x", CodeNotFound))  // same — message is ignored
//
// To check only the topmost code without scanning the chain, call Code() directly.
func (e *Err) Is(target error) bool {
	switch t := target.(type) {
	case *Err:
		return e.code != "" && e.code == t.code
	case Code:
		return e.code != "" && e.code == t
	}
	return false
}

// New creates an *[Err] with the given message and code.
//
// Both message and code must be non-empty; New panics if either is empty.
// If you need a zero-value fallback code, use [CodeInternal] explicitly.
func New(message string, code Code) *Err {
	if code == "" {
		panic("mishap.New: code must not be empty")
	}
	if message == "" {
		panic("mishap.New: message must not be empty")
	}
	return &Err{message: message, code: code}
}

// As extracts the first *[Err] found in err's chain.
// It is a convenience wrapper around [errors.As].
//
//	if mErr, ok := mishap.As(err); ok {
//	    switch mErr.Code() { ... }
//	}
func As(err error) (*Err, bool) {
	var e *Err
	ok := errors.As(err, &e)
	return e, ok
}

// ── Wrap ─────────────────────────────────────────────────────────────────────

type wrapConfig struct {
	code        Code
	defaultCode Code
}

// WrapOption configures the behaviour of [Wrap].
type WrapOption func(*wrapConfig)

// WithCode forces the given code onto the wrapped error, overriding any code
// that would otherwise be inherited from the cause chain.
func WithCode(code Code) WrapOption {
	return func(c *wrapConfig) { c.code = code }
}

// WithDefaultCode sets a fallback code that is used only when [WithCode] was
// not supplied AND no *[Err] with a code exists anywhere in the cause chain.
// If neither applies, [CodeInternal] is used.
func WithDefaultCode(code Code) WrapOption {
	return func(c *wrapConfig) { c.defaultCode = code }
}

// Wrap creates an *[Err] around err with the given message.
// It returns nil when err is nil.
//
// Code resolution order:
//  1. [WithCode] — explicit override
//  2. The code of the first *[Err] found anywhere in the cause chain
//  3. [WithDefaultCode] — caller-supplied fallback
//  4. [CodeInternal]
func Wrap(err error, message string, opts ...WrapOption) *Err {
	if err == nil {
		return nil
	}
	if message == "" {
		panic("mishap.Wrap: message must not be empty")
	}

	cfg := wrapConfig{defaultCode: CodeInternal}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}

	if cfg.code == "" {
		var src *Err
		if errors.As(err, &src) && src.Code() != "" {
			cfg.code = src.Code()
		}
	}
	if cfg.code == "" {
		cfg.code = cfg.defaultCode
	}

	return &Err{
		sourceErr: err,
		message:   message,
		code:      cfg.code,
	}
}
