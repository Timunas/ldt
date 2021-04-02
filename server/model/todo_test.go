package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/timunas/ldt/server/model"
)

func TestTodoCreation(t *testing.T) {
	name := "Some name"
	description := "Some description"
	userID := "uid"
	result := model.NewTodo(name, description, userID)

	assert.Equal(t, name, result.Name)
	assert.Equal(t, description, result.Description)
	assert.Equal(t, userID, result.UserID)
	assert.Empty(t, result.ID)
	assert.True(t, result.CreateAt.IsZero())
	assert.True(t, result.UpdateAt.IsZero())
}
