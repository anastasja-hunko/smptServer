package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	config  *Config
	db      *mongo.Database
	UserCol *UserCol
	LogCol  *LogCol
}

func New(config *Config) *Database {
	return &Database{config: config}
}

//connect to db and ping it
func (c *Database) Open() error {

	clientOptions := options.Client().ApplyURI(c.config.DatabaseURL)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}

	c.db = client.Database(c.config.DatabaseName)

	c.UserCol = c.NewUserCol()

	c.LogCol = c.NewLogCol()

	return nil
}

//close db connection
func (c *Database) Close() error {

	return c.db.Client().Disconnect(context.TODO())

}
