package mishap_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/sunkek/mishap"
)

// ── New ───────────────────────────────────────────────────────────────────────

func TestNew_Error_NoSource(t *testing.T) {
	e := mishap.New("hello", mishap.CodeInternal)
	if got := e.Error(); got != "hello" {
		t.Fatalf("Error() = %q, want %q", got, "hello")
	}
}

func TestNew_PanicsOnEmptyCode(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("New with empty code should panic")
		}
	}()
	mishap.New("something", "")
}

func TestNew_PanicsOnEmptyMessage(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("New with empty message should panic")
		}
	}()
	mishap.New("", mishap.CodeInternal)
}

// ── Wrap ──────────────────────────────────────────────────────────────────────

func TestWrap_Error_WithSource(t *testing.T) {
	src := errors.New("db down")
	e := mishap.Wrap(src, "load user")
	want := "load user: db down"
	if got := e.Error(); got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
}

func TestWrap_Unwrap(t *testing.T) {
	src := errors.New("x")
	e := mishap.Wrap(src, "y")
	if got := errors.Unwrap(e); got != src {
		t.Fatalf("Unwrap() = %v, want %v", got, src)
	}
}

func TestWrap_NilErrReturnsNil(t *testing.T) {
	if got := mishap.Wrap(nil, "x"); got != nil {
		t.Fatalf("Wrap(nil, ...) = %v, want nil", got)
	}
}

// TestWrap_NilErrIsNilThroughErrorReturn guards the typed-nil trap: when Wrap
// returned *Err, the documented one-liner below produced a non-nil error
// interface holding a nil pointer, so a caller's err != nil check took the
// failure branch on success. Asserting on Wrap's result directly does not
// catch this — the comparison has to happen after the value has passed
// through an error return.
func TestWrap_NilErrIsNilThroughErrorReturn(t *testing.T) {
	find := func() error { return nil }
	load := func() error { return mishap.Wrap(find(), "load user") }

	if err := load(); err != nil {
		t.Fatalf("load() = %v (%T), want nil", err, err)
	}
}

func TestWrap_PanicsOnEmptyMessage(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Wrap with empty message should panic")
		}
	}()
	mishap.Wrap(errors.New("cause"), "")
}

func TestWrap_InheritsCode_FromInnerErr(t *testing.T) {
	inner := mishap.New("row not found", mishap.CodeNotFound)
	outer := mustAs(t, mishap.Wrap(inner, "load user"))
	if outer.Code() != mishap.CodeNotFound {
		t.Fatalf("outer.Code = %q, want %q", outer.Code(), mishap.CodeNotFound)
	}
}

func TestWrap_InheritsCode_ThroughFmtWrapped(t *testing.T) {
	inner := mishap.New("row not found", mishap.CodeNotFound)
	wrapped := fmt.Errorf("db: %w", inner)
	outer := mustAs(t, mishap.Wrap(wrapped, "load user"))
	if outer.Code() != mishap.CodeNotFound {
		t.Fatalf("outer.Code = %q, want %q", outer.Code(), mishap.CodeNotFound)
	}
}

func TestWrap_WithCodeOverridesInheritance(t *testing.T) {
	inner := mishap.New("row not found", mishap.CodeNotFound)
	outer := mustAs(t, mishap.Wrap(inner, "load user", mishap.WithCode(mishap.CodeInternal)))
	if outer.Code() != mishap.CodeInternal {
		t.Fatalf("outer.Code = %q, want %q", outer.Code(), mishap.CodeInternal)
	}
}

func TestWrap_DefaultCodeUsedWhenNoInheritableCode(t *testing.T) {
	inner := errors.New("boom")
	outer := mustAs(t, mishap.Wrap(inner, "failed", mishap.WithDefaultCode(mishap.CodeBadRequest)))
	if outer.Code() != mishap.CodeBadRequest {
		t.Fatalf("outer.Code = %q, want %q", outer.Code(), mishap.CodeBadRequest)
	}
}

func TestWrap_DefaultCodeIgnoredWhenInheritableCodeExists(t *testing.T) {
	inner := mishap.New("row not found", mishap.CodeNotFound)
	outer := mustAs(t, mishap.Wrap(inner, "load user", mishap.WithDefaultCode(mishap.CodeBadRequest)))
	if outer.Code() != mishap.CodeNotFound {
		t.Fatalf("outer.Code = %q, want %q", outer.Code(), mishap.CodeNotFound)
	}
}

func TestWrap_FallsBackToCodeInternalByDefault(t *testing.T) {
	outer := mustAs(t, mishap.Wrap(errors.New("raw"), "oh no"))
	if outer.Code() != mishap.CodeInternal {
		t.Fatalf("outer.Code = %q, want %q", outer.Code(), mishap.CodeInternal)
	}
}

// ── errors.Is ─────────────────────────────────────────────────────────────────

func TestIs_CodeTarget(t *testing.T) {
	e := mishap.New("nope", mishap.CodeNotFound)
	if !errors.Is(e, mishap.CodeNotFound) {
		t.Fatalf("errors.Is(e, CodeNotFound) = false, want true")
	}
	if errors.Is(e, mishap.CodeInternal) {
		t.Fatalf("errors.Is(e, CodeInternal) = true, want false")
	}
}

func TestIs_ErrTarget_MatchesByCodeOnly(t *testing.T) {
	e := mishap.New("message A", mishap.CodeNotFound)
	target := mishap.New("message B", mishap.CodeNotFound) // different message, same code
	if !errors.Is(e, target) {
		t.Fatalf("errors.Is should match on code, not message")
	}
}

func TestIs_ScansChain(t *testing.T) {
	inner := mishap.New("row not found", mishap.CodeNotFound)
	outer := mustAs(t, mishap.Wrap(inner, "load user", mishap.WithCode(mishap.CodeInternal)))

	// topmost code is INTERNAL
	if outer.Is(mishap.CodeNotFound) {
		t.Fatalf("outer.Is(CodeNotFound) = true, want false (topmost is INTERNAL)")
	}
	// but the chain contains NOT_FOUND
	if !errors.Is(outer, mishap.CodeNotFound) {
		t.Fatalf("errors.Is(outer, CodeNotFound) = false, want true (inner has NOT_FOUND)")
	}
}

// ── As ────────────────────────────────────────────────────────────────────────

func TestAs_FindsMishapErr(t *testing.T) {
	inner := mishap.New("not found", mishap.CodeNotFound)
	wrapped := fmt.Errorf("service: %w", inner)

	mErr, ok := mishap.As(wrapped)
	if !ok {
		t.Fatal("As() returned false, want true")
	}
	if mErr.Code() != mishap.CodeNotFound {
		t.Fatalf("As().Code() = %q, want %q", mErr.Code(), mishap.CodeNotFound)
	}
}

func TestAs_ReturnsFalseForPlainError(t *testing.T) {
	_, ok := mishap.As(errors.New("plain"))
	if ok {
		t.Fatal("As() returned true for a plain error, want false")
	}
}

func TestAs_ReturnsFalseForNil(t *testing.T) {
	_, ok := mishap.As(nil)
	if ok {
		t.Fatal("As() returned true for nil, want false")
	}
}

// mustAs extracts the *Err from err, failing the test if there is none. These
// tests assert on the topmost code and on Is, both of which need the concrete
// type; Wrap returns error, so As is the documented way back to it.
func mustAs(t *testing.T, err error) *mishap.Err {
	t.Helper()
	e, ok := mishap.As(err)
	if !ok {
		t.Fatalf("As(%v) = false, want an *Err", err)
	}
	return e
}
