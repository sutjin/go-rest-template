package test

import (
	"testing"

	"sutjin/go-rest-template/internal/pkg/models"

	unitTest "github.com/Valiben/gin_unit_test"
	"github.com/Valiben/gin_unit_test/utils"
)

func TestIsAcive(t *testing.T) {
	resp := models.AppInfo{}

	err := unitTest.TestHandlerUnMarshalResp(utils.GET, "/isActive", "json", nil, &resp)
	if err != nil {
		t.Errorf("TestIsAcive: %v\n", err)
		return
	}

	t.Log("Success")
}
