package mishap_test

import (
	"errors"
	"fmt"

	"github.com/sunkek/mishap"
)

func ExampleNew() {
	err := mishap.New("user not found", mishap.CodeNotFound)
	fmt.Println(err.Error())
	fmt.Println(err.Code())

	// Output:
	// user not found
	// NOT_FOUND
}

func ExampleWrap_inheritCodeFromInnerMishap() {
	inner := mishap.New("row not found", mishap.CodeNotFound)

	outer := mishap.Wrap(inner, "load user")
	fmt.Println(outer.Error())
	fmt.Println(outer.Code())

	// Output:
	// load user: row not found
	// NOT_FOUND
}

func ExampleWrap_inheritThroughNonMishapWrapper() {
	inner := mishap.New("row not found", mishap.CodeNotFound)
	wrapped := fmt.Errorf("db: %w", inner)

	outer := mishap.Wrap(wrapped, "load user")
	fmt.Println(outer.Code())

	// Output:
	// NOT_FOUND
}

func ExampleWrap_withCodeOverridesInheritance() {
	inner := mishap.New("row not found", mishap.CodeNotFound)

	outer := mishap.Wrap(inner, "load user", mishap.WithCode(mishap.CodeInternal))
	fmt.Println(outer.Code())

	// Output:
	// INTERNAL_ERROR
}

func ExampleWrap_defaultCodeUsedWhenNoInnerCode() {
	inner := errors.New("boom")

	outer := mishap.Wrap(inner, "failed", mishap.WithDefaultCode(mishap.CodeBadRequest))
	fmt.Println(outer.Code())

	// Output:
	// BAD_REQUEST
}
