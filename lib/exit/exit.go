// Package exit provides a convenient and readable set of exit functions.
package exit

import (
	"fmt"
	"os"
)

var (
	// code is the exit code to be used.
	code int
)

func init() {
	// Assuming healthy exit.
	code = 0
}

// WithCode is an empty struct used to tie together the "Code" and "With"
// commands allowing "exit.Code(-1).With(err)".
type WithCode struct{}

// Code assigns the exit code.
func Code(exitCode int) (wc *WithCode) {
	code = exitCode
	return
}

// With is a convenience wrapper allowing "exit.Code(-1).With(err)".
func (wc *WithCode) With(a ...interface{}) {
	With(a...)
}

// Withf is a convenience wrapper allowing `exit.Code(-1).Withf("%v\n", err)`.
func (wc *WithCode) Withf(format string, a ...interface{}) {
	Withf(format, a...)
}

// With writes the passed values to the screen like Println and exits.
func With(a ...interface{}) {
	fmt.Println(a...)
	os.Exit(code)
}

// Withf writes the passed values to the screen like Printf and exits.
func Withf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
	os.Exit(code)
}
