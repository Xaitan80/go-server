package api

// ErrorResponse is a standard JSON error format for all API endpoints
type ErrorResponse struct {
	Error string `json:"error"`
}
