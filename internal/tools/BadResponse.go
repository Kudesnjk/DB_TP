package tools

const (
	ConstInternalErrorMessage = "internal server error"
	ConstNotFoundMessage      = "not found"
	ConstSomeMessage          = "some message"
)

type BadResponse struct {
	Message string `json:"message"`
}
