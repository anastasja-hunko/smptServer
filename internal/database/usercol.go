package database

import (
	"context"
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

	err := u.BeforeCreate()

	if err != nil {
		return err
	}

	_, err = uc.col.InsertOne(context.TODO(), u)

	if err != nil {
		return err
	}

	return nil
}

func (uc *UserCol) FindByLogin(login string) (*model.User, error) {

	filter := bson.D{primitive.E{Key: "login", Value: login}}

	var user model.User

	err := uc.col.FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (uc *UserCol) Update(u *model.User) error {

	err := u.BeforeCreate()
	if err != nil {
		return err
	}

	filter := bson.D{primitive.E{Key: "login", Value: u.Login}}

	update := bson.D{

		primitive.E{Key: "$set", Value: bson.D{

			primitive.E{Key: "login", Value: u.Login},

			primitive.E{Key: "password", Value: u.Password},
		}},
	}
	_, err = uc.col.UpdateOne(context.TODO(), filter, update)

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

func (uc *UserCol) Delete(login string) error {

	filter := bson.D{primitive.E{Key: "login", Value: login}}

	_, err := uc.col.DeleteOne(context.TODO(), filter)

	return err
}
