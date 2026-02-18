package query

import (
	"strings"
)

// Errno defines the error number and message
type Errno struct {
	Code    int
	Message string
}

func (err Errno) Error() string {
	return err.Message
}

// Err represents a detailed error
type Err struct {
	Code    int
	Message string
	Err     error
}

func (err *Err) Error() string {
	return err.Message
}

var (
	// Common errors
	OK                  = &Errno{Code: 10000, Message: "operation succeeded"}
	InternalServerError = &Errno{Code: 10001, Message: "Internal server error"}
	ErrBind             = &Errno{Code: 10002, Message: "Error occurred while binding the request body to the struct"}
	ErrParam            = &Errno{Code: 10003, Message: "Invalid parameters"}

	ErrDatabase = &Errno{Code: 20002, Message: "Database error"}

	// Resource errors
	ErrZoneNotFound   = &Errno{Code: 40004, Message: "Zone not found"}
	ErrRecordNotFound = &Errno{Code: 40005, Message: "Record not found"}

	// Auth errors (stubs for now)
	ErrNeedAuth     = &Errno{Code: 50107, Message: "Authentication required"}
	ErrNoPermission = &Errno{Code: 50116, Message: "No permission"}
)

// DecodeErr decodes an error into code and message
func DecodeErr(err error) (int, string) {
	if err == nil {
		return OK.Code, OK.Message
	}

	switch typed := err.(type) {
	case *Err:
		return typed.Code, typed.Message
	case *Errno:
		return typed.Code, typed.Message
	default:
		// Clean up standard error messages
		msg := strings.ReplaceAll(err.Error(), "\"", "")
		return InternalServerError.Code, msg
	}
}

// Response defines the standard API response structure
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}
