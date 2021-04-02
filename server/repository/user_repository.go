package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/qiniu/qmgo"
	"github.com/timunas/ldt/server/model"
	"go.mongodb.org/mongo-driver/bson"
)

type UserRepo struct {
	collection *qmgo.Collection
	context    *context.Context
}

func NewUserRepository(collection *qmgo.Collection, context *context.Context) *UserRepo {
	return &UserRepo{
		collection,
		context,
	}
}

func (r *UserRepo) FindByID(id string) (*model.User, error) {
	user := model.User{}
	err := r.collection.Find(*r.context, bson.M{"_id": id}).One(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) FindByEmail(email string) (*model.User, error) {
	user := model.User{}
	err := r.collection.Find(*r.context, bson.M{"email": email}).One(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) Save(u *model.User) (*model.User, error) {
	if len(u.ID) == 0 {
		u.ID = uuid.New().String()
		u.CreateAt = time.Now().Local()
		u.UpdateAt = time.Now().Local()
	} else {
		u.UpdateAt = time.Now().Local()
	}

	_, err := r.collection.Upsert(*r.context, bson.M{"_id": u.ID}, u)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *UserRepo) Delete(u *model.User) error {
	return r.collection.RemoveId(*r.context, u.ID)
}
