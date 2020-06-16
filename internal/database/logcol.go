package database

import (
	"context"
	"github.com/anastasja-hunko/smptServer/internal/model"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type logCol struct {
	col *mongo.Collection
}

func (db *Database) newLogCol() *logCol {

	return &logCol{col: db.db.Collection(db.config.LogColName)}

}

func (dc *logCol) Create(ctx context.Context, l *model.Log) error {

	ctx, cancel := context.WithTimeout(ctx, 2*time.Millisecond)

	defer cancel()

	l.BeforeCreate()

	_, err := dc.col.InsertOne(ctx, l)
	if err != nil {

		return err

	}

	return nil
}
