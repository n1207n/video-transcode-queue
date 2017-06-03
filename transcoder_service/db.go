package main

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

// GetDatabaseConnection returns an instance
// as a database connection
func GetDatabaseConnection(user string, password string, host string, db string) *pg.DB {
	connection := pg.Connect(&pg.Options{
		User:     user,
		Password: password,
		Database: db,
		Addr:     host,
	})

	return connection
}

// CreateSchemas creates a set of database tables
// from Go struct classes
func CreateSchemas(user string, password string, host string, db string) error {
	connection := GetDatabaseConnection(user, password, host, db)

	defer func() {
		connection.Close()
	}()

	for _, model := range []interface{}{&VideoRendering{}, &Video{}} {
		err := connection.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists: true,
		})

		if err != nil {
			return err
		}
	}

	return nil
}
