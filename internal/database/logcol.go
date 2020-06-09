package database

import (
	"context"
	"github.com/anastasja-hunko/smptServer/internal/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type LogCol struct {
	col *mongo.Collection
}

func (db *Database) NewLogCol() *LogCol {

	return &LogCol{col: db.db.Collection(db.config.LogColName)}

}

func (dc *LogCol) Create(l *model.Log) error {

	l.BeforeCreate()

	_, err := dc.col.InsertOne(context.TODO(), l)

	if err != nil {

		return err

	}

	return nil
}
