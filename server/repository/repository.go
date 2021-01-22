package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/qiniu/qmgo"
	"github.com/timunas/ldt/server/model"
	"go.mongodb.org/mongo-driver/bson"
)

type TodoRepo struct {
	collection *qmgo.Collection
	context    *context.Context
}

func NewTodoRepository(collection *qmgo.Collection, context *context.Context) *TodoRepo {
	return &TodoRepo{
		collection,
		context,
	}
}

func (r *TodoRepo) FindAll() ([]*model.Todo, error) {
	todos := []*model.Todo{}
	err := r.collection.Find(*r.context, bson.M{}).All(&todos)

	if err != nil {
		return nil, err
	}

	return todos, nil
}

func (r *TodoRepo) FindByID(id string) (*model.Todo, error) {
	todo := model.Todo{}
	err := r.collection.Find(*r.context, bson.M{"_id": id}).One(&todo)

	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (r *TodoRepo) FindByUserID(id string) ([]*model.Todo, error) {
	todos := []*model.Todo{}
	err := r.collection.Find(*r.context, bson.M{"user_id": id}).All(&todos)

	if err != nil {
		return nil, err
	}

	return todos, nil
}

func (r *TodoRepo) Save(t *model.Todo) (*model.Todo, error) {
	if len(t.ID) == 0 {
		t.ID = uuid.New().String()
		t.CreateAt = time.Now().Local()
		t.UpdateAt = time.Now().Local()
	} else {
		t.UpdateAt = time.Now().Local()
	}

	_, err := r.collection.Upsert(*r.context, bson.M{"_id": t.ID}, t)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func (r *TodoRepo) Delete(t *model.Todo) error {
	return r.collection.RemoveId(*r.context, t.ID)
}
