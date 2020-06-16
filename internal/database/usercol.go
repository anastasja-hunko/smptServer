package database

import (
	"context"
	"errors"
	"github.com/anastasja-hunko/smptServer/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserCol struct {
	col *mongo.Collection
}

func (db *Database) NewUserCol() *UserCol {

	return &UserCol{col: db.db.Collection(db.config.UserColName)}

}

func (uc *UserCol) Create(u *model.User) error {

	err := u.HashPass()
	if err != nil {
		return err
	}

	u.Active = true

	_, err = uc.col.InsertOne(context.TODO(), u)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UserCol) FindByLogin(login string) (*model.User, error) {

	filter := bson.D{primitive.E{Key: "_id", Value: login}}

	var user model.User

	err := uc.col.FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (uc *UserCol) UpdatePassword(u *model.User) error {

	user, _ := uc.FindByLogin(u.Login)

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

	return uc.update(update, u.Login)
}

func (uc *UserCol) UpdateActive(login string) error {

	user, _ := uc.FindByLogin(login)

	histories, logs := user.AppendToHistoryAndLogs("Active", user.Active, false)

	update := bson.D{

		primitive.E{Key: "$set", Value: bson.D{

			primitive.E{Key: "active", Value: false},

			primitive.E{Key: "history", Value: histories},

			primitive.E{Key: "log", Value: logs},
		}},
	}

	return uc.update(update, login)
}

func (uc *UserCol) UpdateUserLog(u *model.User, logMessage string) error {

	logs := u.AppendToLogs(logMessage)

	update := bson.D{

		primitive.E{Key: "$set", Value: bson.D{

			primitive.E{Key: "log", Value: logs},
		}},
	}

	return uc.update(update, u.Login)
}

func (uc *UserCol) UpdateUserMessages(u *model.User, msg *model.Message) error {

	messages := u.AppendToMessages(msg)

	update := bson.D{

		primitive.E{Key: "$set", Value: bson.D{

			primitive.E{Key: "messages", Value: messages},
		}},
	}

	return uc.update(update, u.Login)
}

func (uc *UserCol) update(update primitive.D, login string) error {

	filter := bson.D{primitive.E{Key: "_id", Value: login}}

	_, err := uc.col.UpdateOne(context.TODO(), filter, update)

	return err
}

func (uc *UserCol) FindAll() ([]*model.User, error) {

	cur, err := uc.col.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}

	var results []*model.User

	for cur.Next(context.TODO()) {

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
	cur.Close(context.TODO())

	return results, nil
}
