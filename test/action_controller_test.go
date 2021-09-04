package test

import (
	"testing"

	"sutjin/go-rest-template/internal/pkg/models"
	helper "sutjin/go-rest-template/pkg/tools"

	unitTest "github.com/Valiben/gin_unit_test"
	"github.com/Valiben/gin_unit_test/utils"
)

var test_models = []interface{}{
	(*models.Action)(&models.Action{
		Action: "test_action",
	}),
}

func TestGetActionsEmptyResult(t *testing.T) {
	var resp []models.Action

	err := unitTest.TestHandlerUnMarshalResp(utils.GET, "/actions", "json", nil, &resp)
	if err != nil {
		t.Errorf("TestGetActions: %v\n", err)
		return
	}

	if len(resp) != 0 {
		t.Errorf("Expected %v but got %v", 0, len(resp))
		return
	}

	t.Log("Success")
}

func TestGetActions(t *testing.T) {
	var resp []models.Action

	helper.PopulateDatabase(test_models)

	err := unitTest.TestHandlerUnMarshalResp(utils.GET, "/actions", "json", nil, &resp)
	if err != nil {
		t.Errorf("TestGetActions: %v\n", err)
		return
	}

	if len(resp) != 1 {
		t.Errorf("Expected %v but got %v", 1, len(resp))
		return
	}

	helper.ClearDatabase(test_models)

	t.Log("Success")
}

func TestPostActions(t *testing.T) {
	var action models.Action
	var input = models.Action{
		Action: "test_action",
	}

	err := unitTest.TestHandlerUnMarshalResp(utils.POST, "/actions", "json", input, &action)
	if err != nil {
		t.Errorf("TestPostMappingClientNotFound: %v\n", err)
		return
	}

	if action.Action != input.Action {
		t.Errorf("Expected %v but got %v", input.Action, action.Action)
		return
	}
}
