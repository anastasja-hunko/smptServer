package database

import (
	"context"
	"errors"
	"github.com/anastasja-hunko/smptServer/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type userCol struct {
	col *mongo.Collection
}

func (db *Database) newUserCol() *userCol {

	return &userCol{col: db.db.Collection(db.config.UserColName)}

}

func (uc *userCol) Create(ctx context.Context, u *model.User) error {

	err := u.HashPass()
	if err != nil {
		return err
	}

	u.Active = true

	ctx, cancel := context.WithTimeout(ctx, 2*time.Millisecond)

	defer cancel()

	_, err = uc.col.InsertOne(ctx, u)
	if err != nil {
		return err
	}

	return nil
}

func (uc *userCol) FindByLogin(ctx context.Context, login string) (*model.User, error) {

	filter := bson.D{primitive.E{Key: "_id", Value: login}}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)

	defer cancel()

	var user model.User

	err := uc.col.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (uc *userCol) UpdatePassword(ctx context.Context, u *model.User) error {

	user, _ := uc.FindByLogin(ctx, u.Login)

	equal := user.ComparePasswords(u.Password)

	if equal {

		return errors.New("passwords are equals")

	}

	err := u.HashPass()
	if err != nil {
		return err
	}

	histories, logs := user.AppendToHistoryAndLogs("Password", user.Password, u.Password)

	update := bson.D{

		primitive.E{Key: "$set", Value: bson.D{

			primitive.E{Key: "password", Value: u.Password},

			primitive.E{Key: "history", Value: histories},

			primitive.E{Key: "log", Value: logs},
		}},
	}

	return uc.update(ctx, update, u.Login)
}

func (uc *userCol) UpdateActive(ctx context.Context, login string) error {

	user, _ := uc.FindByLogin(ctx, login)

	histories, logs := user.AppendToHistoryAndLogs("Active", user.Active, false)

	update := bson.D{

		primitive.E{Key: "$set", Value: bson.D{

			primitive.E{Key: "active", Value: false},

			primitive.E{Key: "history", Value: histories},

			primitive.E{Key: "log", Value: logs},
		}},
	}

	return uc.update(ctx, update, login)
}

func (uc *userCol) UpdateUserLog(ctx context.Context, u *model.User, logMessage string) error {

	logs := u.AppendToLogs(logMessage)

	update := bson.D{

		primitive.E{Key: "$set", Value: bson.D{

			primitive.E{Key: "log", Value: logs},
		}},
	}

	return uc.update(ctx, update, u.Login)
}

func (uc *userCol) UpdateUserMessages(ctx context.Context, u *model.User, msg *model.Message) error {

	messages := u.AppendToMessages(msg)

	update := bson.D{

		primitive.E{Key: "$set", Value: bson.D{

			primitive.E{Key: "messages", Value: messages},
		}},
	}

	return uc.update(ctx, update, u.Login)
}

func (uc *userCol) update(ctx context.Context, update primitive.D, login string) error {

	ctx, cancel := context.WithTimeout(ctx, 2*time.Millisecond)

	defer cancel()

	filter := bson.D{primitive.E{Key: "_id", Value: login}}

	_, err := uc.col.UpdateOne(ctx, filter, update)

	return err
}

func (uc *userCol) FindAll(ctx context.Context) ([]*model.User, error) {

	ctx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)

	defer cancel()

	cur, err := uc.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var results []*model.User

	for cur.Next(ctx) {

		var elem model.User

		err := cur.Decode(&elem)
		if err != nil {
			return results, err
		}

		results = append(results, &elem)
	}

	err = cur.Err()
	if err != nil {
		return results, err
	}

	// Close the cursor once finished
	cur.Close(ctx)

	return results, nil
}
