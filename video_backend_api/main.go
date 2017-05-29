package main

import (
	"log"
	"os"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/joho/godotenv"
)

var (
	pgDb, pgUser, pgPassword, pgHost string
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	loadEnvironmentVariables()

	var db *pg.DB = getDatabaseConnection()
	createSchemas(db)
}

// loadEnvironmentVariables loads PostgreSQL
// information from dotenv
func loadEnvironmentVariables() {
	pgDb = os.Getenv("PGDB")
	pgUser = os.Getenv("PGUSER")
	pgPassword = os.Getenv("PGPASSWORD")
	pgHost = os.Getenv("PGHOST")
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
	for _, model := range []interface{}{&Video{}} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists: true,
		})

		if err != nil {
			return err
		}
	}

	return nil
}
