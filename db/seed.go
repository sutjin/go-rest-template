package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-pg/pg/v10"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"sutjin/go-rest-template/internal/pkg/config"
	"sutjin/go-rest-template/internal/pkg/models"
)

func main() {
	config.Setup("data/config.yml")

	database_config := config.Config.Database

	db, err := sql.Open("postgres",
		"postgres://"+database_config.Username+":"+database_config.Password+"@"+database_config.Host+":5432/"+database_config.Dbname+"?sslmode=disable")
	if err != nil {
		fmt.Print(err.Error())
		panic(err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		fmt.Print(err.Error())
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:./db/migrations",
		"postgres", driver)

	if err != nil {
		fmt.Print(err.Error())
		panic(err)
	}

	m.Steps(2)

	fmt.Println("Migration completed")
	fmt.Println("Seeding Database")

	pgdb := pg.Connect(&pg.Options{
		User:     database_config.Username,
		Password: database_config.Password,
		Database: database_config.Dbname,
	})
	defer pgdb.Close()

	// Check if DB is connected
	ctx := context.Background()
	if err := pgdb.Ping(ctx); err != nil {
		fmt.Print(err.Error())
		panic(err)
	}

	models := []interface{}{
		(*models.Action)(&models.Action{
			Action: "PostRegistration",
		}),
	}

	for _, model := range models {
		_, err := pgdb.Model(model).SelectOrInsert()
		if err != nil {
			fmt.Print(err.Error())
			panic(err)
		}
	}

	fmt.Print("Database Seeded")
}
