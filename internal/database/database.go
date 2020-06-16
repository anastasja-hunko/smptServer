package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Database struct
type Database struct {
	config  *Config
	db      *mongo.Database
	userCol *userCol
	logCol  *logCol
}

//New returns initialized database
func New(config *Config) *Database {
	return &Database{config: config}
}

//Open connects to db and ping it
func (c *Database) Open() error {

	clientOptions := options.Client().ApplyURI(c.config.DatabaseURL)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}

	c.db = client.Database(c.config.DatabaseName)

	c.userCol = c.newUserCol()

	c.logCol = c.newLogCol()

	return nil
}

//Close closes db connection
func (c *Database) Close() error {

	return c.db.Client().Disconnect(context.TODO())

}
