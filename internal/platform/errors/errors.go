package errors

type ErrorCode int

const (
	InternalServerError ErrorCode = iota - 1001
	BadRequest          ErrorCode = iota - 2001
)

var errorMessages = map[ErrorCode]string{
	BadRequest:          "invalid request",
	InternalServerError: "internal error",
}

type CustomError struct {
	Message string    `json:"message"`
	Code    ErrorCode `json:"code"`
}

func NewCustomError(code ErrorCode) *CustomError {
	return &CustomError{
		Message: errorMessages[code],
		Code:    code,
	}
}

func GetErrorMessage(code ErrorCode) string {
	return errorMessages[code]
}
