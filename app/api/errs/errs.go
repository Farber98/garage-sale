// Package errs provides types and support related to web error functionality.
package errs

import (
	"errors"
	"fmt"
)

// ErrCode represents an error code in the system.
type ErrCode struct {
	value int
}

// Value returns the integer value of the error code.
func (ec ErrCode) Value() int {
	return ec.value
}

// String returns the string representation of the error code.
func (ec ErrCode) String() string {
	return codeNames[ec]
}

// UnmarshalText implement the unmarshal interface for JSON conversions.
func (ec *ErrCode) UnmarshalText(data []byte) error {
	errName := string(data)

	v, exists := codeNumbers[errName]
	if !exists {
		return fmt.Errorf("err code %q does not exist", errName)
	}

	*ec = v

	return nil
}

// MarshalText implement the marshal interface for JSON conversions.
func (ec ErrCode) MarshalText() ([]byte, error) {
	return []byte(ec.String()), nil
}

// Equal provides support for the go-cmp package and testing.
func (ec ErrCode) Equal(ec2 ErrCode) bool {
	return ec.value == ec2.value
}

// =============================================================================

// Error represents an error in the system.
type Error struct {
	Code    ErrCode `json:"code"`
	Message string  `json:"message"`
}

// New constructs an error based on an app error.
func New(code ErrCode, err error) Error {
	return Error{
		Code:    code,
		Message: err.Error(),
	}
}

// Newf constructs an error based on a error message.
func Newf(code ErrCode, format string, v ...any) Error {
	return Error{
		Code:    code,
		Message: fmt.Sprintf(format, v...),
	}
}

// Error implements the error interface.
func (e Error) Error() string {
	return e.Message
}

// IsError tests the concrete error is of the Error type
func IsError(err error) bool {
	var er Error
	return errors.As(err, &er)
}

// GetError returns a copy of the Error
func GetError(err error) Error {
	var er Error
	if !errors.As(err, &er) {
		return Error{}
	}
	return er
}
