package main

import (
	"os"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

var (
	pgDb, pgUser, pgPassword, pgHost string
)

func main() {
	loadEnvironmentVariables()

	var db *pg.DB = getDatabaseConnection()
	createSchemas(db)
}

// loadEnvironmentVariables loads PostgreSQL
// information from dotenv
func loadEnvironmentVariables() {
	pgDb = os.Getenv("PGDB")
	if len(pgDb) == 0 {
		panic("No PGDB environment variable")
	}

	pgUser = os.Getenv("PGUSER")
	if len(pgUser) == 0 {
		panic("No PGUSER environment variable")
	}

	pgPassword = os.Getenv("PGPASSWORD")
	if len(pgPassword) == 0 {
		panic("No PGPASSWORD environment variable")
	}

	pgHost = os.Getenv("PGHOST")
	if len(pgHost) == 0 {
		panic("No PGHOST environment variable")
	}
}

// getDatabaseConnection returns an instance
// as a database connection
func getDatabaseConnection() *pg.DB {
	db := pg.Connect(&pg.Options{
		User:     pgUser,
		Password: pgPassword,
		Database: pgDb,
		Addr:     pgHost,
	})

	return db
}

// createSchemas creates a set of database tables
// from Go struct classes
func createSchemas(db *pg.DB) error {
	for _, model := range []interface{}{&VideoRendering{}, &Video{}} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists: true,
		})

		if err != nil {
			return err
		}
	}

	return nil
}
