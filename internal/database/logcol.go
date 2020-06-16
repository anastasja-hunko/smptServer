package database

import (
	"context"
	"github.com/anastasja-hunko/smptServer/internal/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type logCol struct {
	col *mongo.Collection
}

func (db *Database) newLogCol() *logCol {

	return &logCol{col: db.db.Collection(db.config.LogColName)}

}

func (dc *logCol) Create(l *model.Log) error {

	l.BeforeCreate()

	_, err := dc.col.InsertOne(context.TODO(), l)
	if err != nil {

		return err

	}

	return nil
}
