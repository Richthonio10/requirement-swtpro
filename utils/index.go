package utils

const (
	HttpErrorCode       = 500
	ValidationErrorCode    = 422
	AuthorizationErrorCode = 401
)

func ErrorMessage(input error) string {
	if input != nil {
		return input.Error()
	}
	return ""
}