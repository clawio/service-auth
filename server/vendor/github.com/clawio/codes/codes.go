package codes

import "fmt"

// A Code is an unsigned 32-bit error code.
type Code uint32

const (
	// InvalidToken is returned when the auth token is invalid or has expired
	InvalidToken = 89

	// Unauthenticated is returned when authentication is needed for execution.
	Unauthenticated Code = 132

	// BadAuthenticationData is returned when the authentication fails.
	BadAuthenticationData Code = 215

	// BadInputData is returned when the input parameters are not valid.
	BadInputData Code = 400

	// Internal is returned when there is an unexpected/undesired problem
	Internal = 500
)

// String returns a string representation of the Code
func (c Code) String() string {
	switch c {
	case InvalidToken:
		return "Invalid or expired token"
	case Unauthenticated:
		return "Unauthenticated request"
	case BadAuthenticationData:
		return "Bad authentication data"
	case BadInputData:
		return "Bad input data"
	case Internal:
		return "Internal error. Please submit a query to the support team"
	default:
		return "FIXME: this should be a helpful message"
	}
}

// APIErr represents the error returned from the logical layer.
// It can be converted to JSON in REST layer to give end users a better feedback.
type APIErr struct {
	Message string `json:"message"`
	Code    Code   `json:"code"`
}

// NewAPIErr is a usefull function to create APIErrs with the corresponding Code message.
func NewAPIErr(c Code) *APIErr {
	return &APIErr{c.String(), c}
}

// Error() implements the Error interface.
func (e *APIErr) Error() string {
	return fmt.Sprintf("apierror: code=%d message=%s", e.Code, e.Code.String())
}
