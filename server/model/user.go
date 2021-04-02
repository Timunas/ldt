package model

import (
	"time"
)

type User struct {
	ID       string    `json:"id" bson:"_id"`
	CreateAt time.Time `json:"created_at" bson:"createAt"`
	UpdateAt time.Time `json:"updated_at" bson:"updateAt"`

	Name  string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
}

func NewUser(name string, email string) *User {
	u := User{Name: name, Email: email}
	return &u
}

type UserRepository interface {
	FindByID(id string) (*User, error)
	FindByEmail(email string) (*User, error)
	Save(user *User) (*User, error)
	Delete(user *User) error
}
