package tools

import (
	"fmt"

	pg_local "sutjin/go-rest-template/internal/pkg/db"
)

func PopulateDatabase(test_models []interface{}) error {
	for _, model := range test_models {
		_, err := pg_local.DB.Model(model).Insert()
		if err != nil {
			fmt.Print(err.Error())
			return err
		}
	}

	return nil
}

func ClearDatabase(test_models []interface{}) error {
	for _, model := range test_models {
		_, err := pg_local.DB.Model(model).Delete()
		if err != nil {
			fmt.Print(err.Error())
			return err
		}
	}

	return nil
}
