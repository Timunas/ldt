package model

import (
	"time"
)

type Todo struct {
	ID       string    `json:"id" bson:"_id"`
	CreateAt time.Time `json:"created_at" bson:"createAt"`
	UpdateAt time.Time `json:"updated_at" bson:"updateAt"`

	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	UserID      string `json:"-" bson:"user_id"`
}

func NewTodo(name string, description string, userID string) *Todo {
	t := Todo{Name: name, Description: description, UserID: userID}
	return &t
}

type TodoRepository interface {
	FindAll() ([]*Todo, error)
	FindByID(id string) (*Todo, error)
	FindByUserID(id string) ([]*Todo, error)
	Save(todo *Todo) (*Todo, error)
	Delete(todo *Todo) error
}
