package runtime

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestItemInterface(t *testing.T) {
	// Create a test item
	testItem := item{
		title: "Test Title",
		desc:  "Test Description",
	}

	// Test the Title method
	if title := testItem.Title(); title != "Test Title" {
		t.Errorf("Title() should return 'Test Title', but returned %s", title)
	}

	// Test the Description method
	if desc := testItem.Description(); desc != "Test Description" {
		t.Errorf("Description() should return 'Test Description', but returned %s", desc)
	}

	// Test the FilterValue method
	if filter := testItem.FilterValue(); filter != "Test Title" {
		t.Errorf("FilterValue() should return 'Test Title', but returned %s", filter)
	}
}

func TestInitModel(t *testing.T) {
	// Create the model
	model := InitModel()

	// Verify that the model was created
	if model == nil {
		t.Fatal("InitModel() should not return nil")
	}

	// Verify that the model is of the correct type
	if _, ok := model.(*InteractiveConfig); !ok {
		t.Errorf("InitModel() should return a *InteractiveConfig, but returned %T", model)
	}
}

func TestInteractiveConfigInitialState(t *testing.T) {
	// Create a model
	model := InitModel()
	ic, ok := model.(*InteractiveConfig)
	if !ok {
		t.Fatal("InitModel() should return a *InteractiveConfig")
	}

	// Verify the initial state
	if ic.currentStep != 0 {
		t.Errorf("currentStep should be 0, but is %d", ic.currentStep)
	}

	if ic.quitting {
		t.Error("quitting should be false initially")
	}
}

func TestInteractiveConfigView(t *testing.T) {
	// Create a model
	model := InitModel()
	ic, ok := model.(*InteractiveConfig)
	if !ok {
		t.Fatal("InitModel() should return a *InteractiveConfig")
	}

	// Test the View method
	view := ic.View()

	// The view should not be empty
	if view == "" {
		t.Error("View() should not return an empty string")
	}

	// Simulate an exit condition
	ic.quitting = true
	view = ic.View()

	// The view should contain an exit message
	if !strings.Contains(strings.ToLower(view), "config") {
		t.Error("View() should contain a reference to the configuration when quitting=true")
	}
}

func TestInteractiveConfigUpdate(t *testing.T) {
	// Create a model
	model := InitModel()
	ic, ok := model.(*InteractiveConfig)
	if !ok {
		t.Fatal("InitModel() should return a *InteractiveConfig")
	}

	// Test the Update method with different messages

	// Exit message
	model, cmd := ic.Update(tea.QuitMsg{})
	if cmd != nil {
		t.Error("Update(QuitMsg) should return a nil cmd")
	}

	// Verify that the model was updated
	updatedIc, ok := model.(*InteractiveConfig)
	if !ok {
		t.Fatal("Update() should return a *InteractiveConfig")
	}

	if !updatedIc.quitting {
		t.Error("Update(QuitMsg) should set quitting=true")
	}
}
