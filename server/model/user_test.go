package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/timunas/ldt/server/model"
)

func TestUserCreation(t *testing.T) {
	name := "Some name"
	email := "some@email.com"
	result := model.NewUser(name, email)

	assert.Equal(t, name, result.Name)
	assert.Equal(t, email, result.Email)
	assert.Empty(t, result.ID)
	assert.True(t, result.CreateAt.IsZero())
	assert.True(t, result.UpdateAt.IsZero())
}
