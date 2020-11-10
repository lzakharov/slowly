package slowly

// Status is a response status.
type Status string

const (
	// StatusOK is a success response status.
	StatusOK Status = "ok"
)

// SlowRequest is a 'slow' request.
type SlowRequest struct {
	Timeout int64 `json:"timeout"`
}

// SlowResponse is a 'slow' action response.
type SlowResponse struct {
	Status Status `json:"status"`
}

var (
	// TimeoutTooLongResponse is the error response returned if the timeout is exceeded.
	TimeoutTooLongResponse = ErrorResponse{Error: "timeout too long"}
)

// ErrorResponse is a error response.
type ErrorResponse struct {
	Error string `json:"error"`
}
