package handler

type RequestContextUserIDKey struct{}

type TodoRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
