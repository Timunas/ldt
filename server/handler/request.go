package handler

type TodoRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
