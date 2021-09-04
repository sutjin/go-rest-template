package test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"sutjin/go-rest-template/internal/api/controllers"
	pg_local "sutjin/go-rest-template/internal/pkg/db"
	"sutjin/go-rest-template/internal/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"

	unitTest "github.com/Valiben/gin_unit_test"
)

func init() {
	// initialize the router
	router := gin.Default()

	// Handlers for testing
	router.GET("/isActive", controllers.GetVersion)
	router.GET("/actions", controllers.GetActions)
	router.POST("/actions", controllers.PostAction)

	// Setup the router
	unitTest.SetRouter(router)
	newLog := log.New(os.Stdout, "", log.Llongfile|log.Ldate|log.Ltime)
	unitTest.SetLog(newLog)
}

func createSchema(db *pg.DB) error {
	models := []interface{}{
		(*models.Action)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp:        true,
			IfNotExists: true,
		})
		if err != nil {
			fmt.Print(err.Error())
			return err
		}
	}

	return nil
}

func deleteSchema(db *pg.DB) error {
	models := []interface{}{
		(*models.Action)(nil),
	}

	for _, model := range models {
		err := db.Model(model).DropTable(&orm.DropTableOptions{
			IfExists: true,
		})
		if err != nil {
			fmt.Print(err.Error())
			return err
		}
	}

	return nil
}

func TestMain(m *testing.M) {
	// user and password will need to match running postgres instance
	pgdb := pg.Connect(&pg.Options{
		Addr:     os.Getenv("DATABASE_URL") + ":5432",
		User:     "postgres",
		Password: "1234",
	})
	defer pgdb.Close()

	// Check if DB is connected
	ctx := context.Background()
	if err := pgdb.Ping(ctx); err != nil {
		fmt.Print(err.Error())
		panic(err)
	}

	createSchema(pgdb)

	// replaced package DB to our mock DB
	pg_local.DB = pgdb

	log.Println("Database setup for test")
	exitVal := m.Run()
	log.Println("Database dropped after test")

	deleteSchema(pgdb)

	os.Exit(exitVal)
}
