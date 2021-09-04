package db

import (
	"context"

	"sutjin/go-rest-template/internal/pkg/config"

	"github.com/go-pg/pg/v10"
)

var DB *pg.DB

// SetupDB opens a database and saves the reference to `Database` struct.
func SetupDB() {
	configuration := config.GetConfig()

	database := configuration.Database.Dbname
	host := configuration.Database.Host
	port := configuration.Database.Port
	username := configuration.Database.Username
	password := configuration.Database.Password

	DB = pg.Connect(&pg.Options{
		Addr:     host + ":" + port,
		User:     username,
		Password: password,
		Database: database,
	})

	// Check if database is connected
	ctx := context.Background()

	if err := DB.Ping(ctx); err != nil {
		panic(err)
	}
}
