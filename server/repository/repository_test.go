package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/qiniu/qmgo"
	"github.com/stretchr/testify/assert"
	"github.com/timunas/ldt/server/model"
	"github.com/timunas/ldt/server/repository"
)

var database = "test"
var coll = "todo"

func TestCreationAndFetching(t *testing.T) {
	client, ctx, err := initDB()
	assert.NoError(t, err)
	defer client.Close(*ctx)
	db := client.Database(database)
	todoCollection := db.Collection(coll)
	repo := repository.NewTodoRepository(todoCollection, ctx)

	beforeTestTime := time.Now()
	todo := model.NewTodo("Some name", "Some description")

	savedTodo, err := repo.Save(todo)

	assert.NoError(t, err)
	assert.NotEmpty(t, savedTodo.ID)
	assert.True(t, savedTodo.CreateAt.After(beforeTestTime))
	assert.True(t, savedTodo.UpdateAt.After(beforeTestTime))

	fetchedTodo, err := repo.FindByID(savedTodo.ID)
	assert.NoError(t, err)
	assert.Equal(t, savedTodo.ID, fetchedTodo.ID)
	assert.Equal(t, savedTodo.Name, fetchedTodo.Name)
	assert.Equal(t, savedTodo.Description, fetchedTodo.Description)
	assert.Equal(t, savedTodo.CreateAt.Unix(), fetchedTodo.CreateAt.Unix())
	assert.Equal(t, savedTodo.UpdateAt.Unix(), fetchedTodo.UpdateAt.Unix())
}

func TestSaveExistingTodo(t *testing.T) {
	client, ctx, err := initDB()
	assert.NoError(t, err)
	defer client.Close(*ctx)
	db := client.Database(database)
	todoCollection := db.Collection(coll)
	repo := repository.NewTodoRepository(todoCollection, ctx)

	savedTodo, err := repo.Save(model.NewTodo("Some name", "Some description"))
	assert.NoError(t, err)
	assert.NotEmpty(t, savedTodo.ID)

	newDescription := "new description"
	savedTodo.Description = newDescription
	beforeTestTime := savedTodo.UpdateAt

	time.Sleep(time.Second)

	_, err = repo.Save(savedTodo)
	assert.NoError(t, err)

	updatedTodo, err := repo.FindByID(savedTodo.ID)
	assert.NoError(t, err)
	assert.Equal(t, savedTodo.ID, updatedTodo.ID)
	assert.Equal(t, savedTodo.Name, updatedTodo.Name)
	assert.Equal(t, newDescription, updatedTodo.Description)
	assert.Equal(t, savedTodo.CreateAt.Unix(), updatedTodo.CreateAt.Unix())
	assert.Less(t, beforeTestTime.Unix(), updatedTodo.UpdateAt.Unix())
}

func TestFindAll(t *testing.T) {
	client, ctx, err := initDB()
	assert.NoError(t, err)
	defer client.Close(*ctx)
	db := client.Database(database)
	todoCollection := db.Collection(coll)
	repo := repository.NewTodoRepository(todoCollection, ctx)

	name := "Some name"
	description := "Some description"
	firstTodo, err := repo.Save(model.NewTodo(name, description))
	assert.NoError(t, err)
	secondTodo, err := repo.Save(model.NewTodo(name, description))
	assert.NoError(t, err)

	todos, err := repo.FindAll()
	assert.NoError(t, err)
	assert.Equal(t, firstTodo.ID, todos[0].ID)
	assert.Equal(t, firstTodo.Name, todos[0].Name)
	assert.Equal(t, firstTodo.Description, todos[0].Description)
	assert.Equal(t, firstTodo.CreateAt.Unix(), todos[0].CreateAt.Unix())
	assert.Equal(t, firstTodo.UpdateAt.Unix(), todos[0].UpdateAt.Unix())
	assert.Equal(t, secondTodo.ID, todos[1].ID)
	assert.Equal(t, secondTodo.Name, todos[1].Name)
	assert.Equal(t, secondTodo.Description, todos[1].Description)
	assert.Equal(t, secondTodo.CreateAt.Unix(), todos[1].CreateAt.Unix())
	assert.Equal(t, secondTodo.UpdateAt.Unix(), todos[1].UpdateAt.Unix())
}

func TestFindByUserID(t *testing.T) {
	client, ctx, err := initDB()
	assert.NoError(t, err)
	defer client.Close(*ctx)
	db := client.Database(database)
	todoCollection := db.Collection(coll)
	repo := repository.NewTodoRepository(todoCollection, ctx)

	firstTodo := model.NewTodo("Some name", "Some description")
	secondTodo := model.NewTodo("Some name", "Some description")
	firstTodo.UserID = uuid.New().String()
	secondTodo.UserID = firstTodo.UserID

	firstTodo, err = repo.Save(firstTodo)
	assert.NoError(t, err)
	secondTodo, err = repo.Save(secondTodo)
	assert.NoError(t, err)

	todos, err := repo.FindByUserID(firstTodo.UserID)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(todos))
	assert.Equal(t, firstTodo.ID, todos[0].ID)
	assert.Equal(t, secondTodo.ID, todos[1].ID)
}

func TestDelete(t *testing.T) {
	client, ctx, err := initDB()
	assert.NoError(t, err)
	defer client.Close(*ctx)
	db := client.Database(database)
	todoCollection := db.Collection(coll)
	repo := repository.NewTodoRepository(todoCollection, ctx)

	todo, err := repo.Save(model.NewTodo("Some name", "Some description"))
	assert.NoError(t, err)
	_, err = repo.FindByID(todo.ID)
	assert.NoError(t, err)

	err = repo.Delete(todo)
	assert.NoError(t, err)
	_, err = repo.FindByID(todo.ID)
	assert.Error(t, err)
}

func initDB() (*qmgo.Client, *context.Context, error) {
	ctx := context.Background()
	client, err := qmgo.NewClient(ctx, &qmgo.Config{Uri: "mongodb://localhost:27017"})

	if err != nil {
		return nil, nil, err
	}

	err = client.Database(database).DropDatabase(ctx)

	return client, &ctx, err
}
