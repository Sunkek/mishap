# Mishap

This package generalizes my error handling approach.

## Install

Add it to your project with:

```sh
go get github.com/sunkek/mishap
```

## Usage

### Create

* Create a custom code

You can use the default error codes (`mishap.Code...`) or create your own codes.

```go
var ErrCodeOverheat = mishap.Code("OVERHEAT")
```

* Create a root error

You can create new errors easily, wherever you need them or in your common errors package to be imported across your project.

```go
err := mishap.New("temperature over 9000C", ErrCodeOverheat)
```

### Wrap

Wrap any error with a `mishap` error. Code precedence when wrapping: 
1. `WithCode`
2. Inherited code from the first `*mishap.Err` in chain
3. `WithDefaultCode`
4. `CodeInternal`


```go
err := errors.New("forge overheat") // No `mishap.Err` in chain
err = mishap.Wrap(err, "critical failure") // Code defaults to `mishap.CodeInternal`
```

```go
err := errors.New("forge overheat") // No `mishap.Err` in chain
var ErrCodeCriticalFailure = mishap.Code("CRITICAL_FAILURE") // Use a custom code if you want
err = mishap.Wrap( // Code defaults to the provided ErrCodeCriticalFailure
    err, 
    "critical failure", 
    mishap.WithDefaultCode(ErrCodeCriticalFailure),
)
```

```go
var ErrCodeOverheat = mishap.Code("OVERHEAT") // Use a custom code if you want
err := mishap.New("forge overheat", ErrCodeOverheat) // Create a new `mishap.Err`
var ErrCodeCriticalFailure = mishap.Code("CRITICAL_FAILURE") // Use a custom code if you want
err = mishap.Wrap( // Results with `ErrCodeOverheat`. Inheritance takes precedence over `mishap.WithDefaultCode` 
    err, 
    "critical failure",
    mishap.WithDefaultCode(ErrCodeCriticalFailure),
)
```

```go
var ErrCodeOverheat = mishap.Code("OVERHEAT") // Use a custom code if you want
err := mishap.New("forge overheat", ErrCodeOverheat) // Create a new `mishap.Err`
err = mishap.Wrap(err, "critical failure") // Inherits the ErrCodeOverheat from the previous err
```

```go
var ErrCodeOverheat = mishap.Code("OVERHEAT") // Use a custom code if you want
err := mishap.New("forge overheat", ErrCodeOverheat) // Create a new `mishap.Err`

var ErrCodeCriticalFailure = mishap.Code("CRITICAL_FAILURE") // Use a custom code if you want
err = mishap.Wrap( // Force ErrCodeCriticalFailure instead of code inheritance
    err, 
    "critical failure", 
    mishap.WithCode(ErrCodeCriticalFailure),
)
```

### Handle

A `mishap.Err` carries a `Code` that you can match with `errors.Is`. You can lookup the whole error chain:

```go
func handle(err error) {
	if err == nil {
		return
	}

	switch {
	case errors.Is(err, ErrCodeOverheat):
        // Sound the alarms
	case errors.Is(err, mishap.CodeNotFound):
		// return 404
	case errors.Is(err, mishap.CodeBadRequest), errors.Is(err, mishap.CodeValidation):
		// return 400
	case errors.Is(err, mishap.CodeUnauthorized):
		// return 401
	case errors.Is(err, mishap.CodeForbidden):
		// return 403
	default:
		// return 500 / log as error
	}
}
```

Or check only the topmost error code:

```go
func handle(err error) {
	if err == nil {
		return
	}

    mErr, ok := err.(*mishap.Err)
    if !ok {
        // return 500 / log as error
    }

	switch mErr.Code() {
	case ErrCodeOverheat:
        // Sound the alarms
	case mishap.CodeNotFound:
		// return 404
	case mishap.CodeBadRequest, mishap.CodeValidation:
		// return 400
	case mishap.CodeUnauthorized:
		// return 401
	case mishap.CodeForbidden:
		// return 403
	default:
		// return 500 / log as error
	}
}
```

### Log

`mishap.Err.Error()` prints the full error chain:

```go
var ErrCodeOverheat = mishap.Code("OVERHEAT")
err := mishap.New("temperature is over 9000C", ErrCodeOverheat)
var ErrCodeCriticalFailure = mishap.Code("CRITICAL_FAILURE")
err = mishap.Wrap(err, "extreme danger", mishap.WithDefaultCode(ErrCodeCriticalFailure))
fmt.Println(err.Code(), err.Error()) // OVERHEAT extreme danger: temperature is over 9000C
```