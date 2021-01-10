package tools

const (
	ConstInternalErrorMessage = "internal server error"
	ConstNotFoundMessage      = "not found message"
	ConstSomeMessage          = "some message"
	ConstUserNotFoundError    = "23503"
	ConstParentNotFound       = "12345"
)

type BadResponse struct {
	Message string `json:"message"`
}
