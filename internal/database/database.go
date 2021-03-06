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
	UserCol *userCol
	LogCol  *logCol
}

//New returns initialized database
func New(config *Config) *Database {

	return &Database{config: config}

}

//Open connects to db
func (c *Database) Open() error {

	clientOptions := options.Client().ApplyURI(c.config.DatabaseURL)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}

	c.db = client.Database(c.config.DatabaseName)

	c.UserCol = c.newUserCol()

	c.LogCol = c.newLogCol()

	return nil
}

//Close closes db connection
func (c *Database) Close() error {

	return c.db.Client().Disconnect(context.TODO())

}
