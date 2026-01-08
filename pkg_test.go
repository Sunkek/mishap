package mishap_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/sunkek/mishap"
)

func TestErr_Error_NoSource(t *testing.T) {
	e := mishap.New("hello", mishap.CodeInternal)
	if got := e.Error(); got != "hello" {
		t.Fatalf("Error() = %q, want %q", got, "hello")
	}
}

func TestErr_Error_WithSource(t *testing.T) {
	src := errors.New("db down")
	e := mishap.Wrap(src, "load user")
	want := "load user: db down"
	if got := e.Error(); got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
}

func TestErr_Unwrap(t *testing.T) {
	src := errors.New("x")
	e := mishap.Wrap(src, "y")
	if got := errors.Unwrap(e); got != src {
		t.Fatalf("Unwrap() = %v, want %v", got, src)
	}
}

func TestErr_Is_CodeTarget(t *testing.T) {
	e := mishap.New("nope", mishap.CodeNotFound)
	if !errors.Is(e, mishap.CodeNotFound) {
		t.Fatalf("errors.Is(e, CodeNotFound) = false, want true")
	}
	if errors.Is(e, mishap.CodeInternal) {
		t.Fatalf("errors.Is(e, CodeInternal) = true, want false")
	}
}

func TestErr_Is_ErrTarget(t *testing.T) {
	e := mishap.New("nope", mishap.CodeNotFound)
	target := mishap.New("nope", mishap.CodeNotFound)
	if !errors.Is(e, target) {
		t.Fatalf("errors.Is(e, target-with-same-code) = false, want true")
	}
}

func TestIs_ScansChain(t *testing.T) {
	inner := mishap.New("row not found", mishap.CodeNotFound)
	outer := mishap.Wrap(inner, "load user", mishap.WithCode(mishap.CodeInternal))

	if outer.Is(mishap.CodeNotFound) {
		t.Fatalf("outer.Is(CodeNotFound) = true, want false (topmost is INTERNAL)")
	}
	if !errors.Is(outer, mishap.CodeNotFound) {
		t.Fatalf("errors.Is(outer, CodeNotFound) = false, want true (inner has NOT_FOUND)")
	}
}

func TestWrap_NilErrReturnsNil(t *testing.T) {
	if got := mishap.Wrap(nil, "x"); got != nil {
		t.Fatalf("Wrap(nil, ...) = %v, want nil", got)
	}
}

func TestWrap_InheritsCode_FromInnerErr(t *testing.T) {
	inner := mishap.New("row not found", mishap.CodeNotFound)
	outer := mishap.Wrap(inner, "load user")

	if outer.Code() != mishap.CodeNotFound {
		t.Fatalf("outer.Code = %q, want %q", outer.Code(), mishap.CodeNotFound)
	}
}

func TestWrap_InheritsCode_ThroughFmtWrapped(t *testing.T) {
	inner := mishap.New("row not found", mishap.CodeNotFound)
	wrapped := fmt.Errorf("db: %w", inner)

	outer := mishap.Wrap(wrapped, "load user")
	if outer.Code() != mishap.CodeNotFound {
		t.Fatalf("outer.Code = %q, want %q", outer.Code(), mishap.CodeNotFound)
	}
}

func TestWrap_WithCodeOverridesInheritance(t *testing.T) {
	inner := mishap.New("row not found", mishap.CodeNotFound)
	outer := mishap.Wrap(inner, "load user", mishap.WithCode(mishap.CodeInternal))

	if outer.Code() != mishap.CodeInternal {
		t.Fatalf("outer.Code = %q, want %q", outer.Code(), mishap.CodeInternal)
	}
}

func TestWrap_DefaultCodeUsedWhenNoInheritableCode(t *testing.T) {
	inner := errors.New("boom") // not mishap.Err
	outer := mishap.Wrap(inner, "failed", mishap.WithDefaultCode(mishap.CodeBadRequest))

	if outer.Code() != mishap.CodeBadRequest {
		t.Fatalf("outer.Code = %q, want %q", outer.Code(), mishap.CodeBadRequest)
	}
}

func TestWrap_DefaultCodeIgnoredWhenInheritableCodeExists(t *testing.T) {
	inner := mishap.New("row not found", mishap.CodeNotFound)
	outer := mishap.Wrap(inner, "load user", mishap.WithDefaultCode(mishap.CodeBadRequest))

	// Precedence: WithCode > inherited src.Code > DefaultCode
	if outer.Code() != mishap.CodeNotFound {
		t.Fatalf("outer.Code = %q, want %q", outer.Code(), mishap.CodeNotFound)
	}
}
