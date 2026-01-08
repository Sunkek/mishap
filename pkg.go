package mishap

import (
	"errors"
	"fmt"
)

// Err represents a structured error with optional fields.
type Err struct {
	sourceErr error  // Underlying error (optional)
	message   string // Required human-readable message
	code      Code   // String error code
}

// Message returns the error's message.
func (e *Err) Message() string { return e.message }

// Code returns the error's code.
func (e *Err) Code() Code { return e.code }

// Error implements the error interface.
func (e *Err) Error() string {
	if e.sourceErr != nil {
		return fmt.Sprintf("%s: %v", e.message, e.sourceErr.Error())
	}
	return e.message
}

// Unwrap allows errors.Is and errors.As to work with Err.
func (e *Err) Unwrap() error {
	return e.sourceErr
}

// Is compares the topmost Err's code to the target code
func (e *Err) Is(target error) bool {
	if t, ok := target.(*Err); ok {
		return e.code != "" && e.code == t.code
	}
	if c, ok := target.(Code); ok {
		return e.code != "" && e.code == c
	}
	return false
}

// New creates an *Err with message and code
func New(message string, code Code) *Err {
	if code == "" {
		code = CodeInternal
	}
	if message == "" {
		message = "error"
	}
	return &Err{
		message: message,
		code:    code,
	}
}

type wrapConfig struct {
	code        Code
	defaultCode Code
}

type WrapOption func(*wrapConfig)

// Enforces the passed code to the error. Overrides possible code inheritance
func WithCode(code Code) WrapOption {
	return func(c *wrapConfig) {
		c.code = code
	}
}

// Enforces the passed code to the error IF WithCode option wasn't specified AND the source error has no code
func WithDefaultCode(code Code) WrapOption {
	return func(c *wrapConfig) {
		c.defaultCode = code
	}
}

// Wrap creates an *Err around err with message and optional behavior.
// Sets error code based on the following order: `WithCode` > err.Code (if *Err) > `DefaultCode`
func Wrap(err error, message string, opts ...WrapOption) *Err {
	if err == nil {
		return nil
	}

	cfg := wrapConfig{
		defaultCode: CodeInternal,
	}
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
	if message == "" {
		message = "error"
	}

	return &Err{
		sourceErr: err,
		message:   message,
		code:      cfg.code,
	}
}
