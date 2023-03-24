package errors

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

type errorTest struct {
	name     string
	err      error
	expected bool
}

func testError(t *testing.T, f func(error) bool, tests ...errorTest) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := f(test.err)
			if result != test.expected {
				t.Errorf("Expected %v, got %v for %v", test.expected, result, test.err)
			}
		})
	}
}

func TestIsIllegalArgument(t *testing.T) {
	testError(t, IsIllegalArgument,
		errorTest{name: "nil", err: nil, expected: false},
		errorTest{name: "IllegalArgument", err: IllegalArgument(), expected: true},
		errorTest{name: "OutOfBounds", err: OutOfBounds(), expected: false},
		errorTest{name: "StackEmpty", err: StackEmpty(), expected: false},
		errorTest{name: "UnsupportedError", err: UnsupportedError("unsupported"), expected: false},
		// A random error generated outside of this package
		errorTest{name: "random", err: fmt.Errorf("%v", time.Now()), expected: false},
		errorTest{name: "new", err: errors.New(time.Now().Format(time.RFC3339)), expected: false},
	)
}

func TestIsOutOfBounds(t *testing.T) {
	testError(t, IsOutOfBounds,
		errorTest{name: "nil", err: nil, expected: false},
		errorTest{name: "IllegalArgument", err: IllegalArgument(), expected: false},
		errorTest{name: "OutOfBounds", err: OutOfBounds(), expected: true},
		errorTest{name: "StackEmpty", err: StackEmpty(), expected: false},
		errorTest{name: "UnsupportedError", err: UnsupportedError("unsupported"), expected: false},
		// A random error generated outside of this package
		errorTest{name: "random", err: fmt.Errorf("%v", time.Now()), expected: false},
		errorTest{name: "new", err: errors.New(time.Now().Format(time.RFC3339)), expected: false},
	)
}

func TestIsStackEmpty(t *testing.T) {
	testError(t, IsStackEmpty,
		errorTest{name: "nil", err: nil, expected: false},
		errorTest{name: "IllegalArgument", err: IllegalArgument(), expected: false},
		errorTest{name: "OutOfBounds", err: OutOfBounds(), expected: false},
		errorTest{name: "StackEmpty", err: StackEmpty(), expected: true},
		errorTest{name: "UnsupportedError", err: UnsupportedError("unsupported"), expected: false},
		// A random error generated outside of this package
		errorTest{name: "random", err: fmt.Errorf("%v", time.Now()), expected: false},
		errorTest{name: "new", err: errors.New(time.Now().Format(time.RFC3339)), expected: false},
	)
}

func TestIsUnsupportedError(t *testing.T) {
	testError(t, IsUnsupportedError,
		errorTest{name: "nil", err: nil, expected: false},
		errorTest{name: "IllegalArgument", err: IllegalArgument(), expected: false},
		errorTest{name: "OutOfBounds", err: OutOfBounds(), expected: false},
		errorTest{name: "StackEmpty", err: StackEmpty(), expected: false},
		errorTest{name: "UnsupportedError", err: UnsupportedError("unsupported"), expected: true},
		// A random error generated outside of this package
		errorTest{name: "random", err: fmt.Errorf("%v", time.Now()), expected: false},
		errorTest{name: "new", err: errors.New(time.Now().Format(time.RFC3339)), expected: false},
	)
}
