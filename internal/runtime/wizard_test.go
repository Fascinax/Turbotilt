package runtime

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestItemInterface(t *testing.T) {
	// Créer un item de test
	testItem := item{
		title: "Test Title",
		desc:  "Test Description",
	}

	// Tester la méthode Title
	if title := testItem.Title(); title != "Test Title" {
		t.Errorf("Title() devrait retourner 'Test Title', mais a retourné %s", title)
	}

	// Tester la méthode Description
	if desc := testItem.Description(); desc != "Test Description" {
		t.Errorf("Description() devrait retourner 'Test Description', mais a retourné %s", desc)
	}

	// Tester la méthode FilterValue
	if filter := testItem.FilterValue(); filter != "Test Title" {
		t.Errorf("FilterValue() devrait retourner 'Test Title', mais a retourné %s", filter)
	}
}

func TestInitModel(t *testing.T) {
	// Créer le modèle
	model := InitModel()

	// Vérifier que le modèle a été créé
	if model == nil {
		t.Fatal("InitModel() ne devrait pas retourner nil")
	}

	// Vérifier que le modèle est du bon type
	if _, ok := model.(*InteractiveConfig); !ok {
		t.Errorf("InitModel() devrait retourner un *InteractiveConfig, mais a retourné %T", model)
	}
}

func TestInteractiveConfigInitialState(t *testing.T) {
	// Créer un modèle
	model := InitModel()
	ic, ok := model.(*InteractiveConfig)
	if !ok {
		t.Fatal("InitModel() devrait retourner un *InteractiveConfig")
	}

	// Vérifier l'état initial
	if ic.currentStep != 0 {
		t.Errorf("currentStep devrait être 0, mais est %d", ic.currentStep)
	}

	if ic.quitting {
		t.Error("quitting devrait être false initialement")
	}
}

func TestInteractiveConfigView(t *testing.T) {
	// Créer un modèle
	model := InitModel()
	ic, ok := model.(*InteractiveConfig)
	if !ok {
		t.Fatal("InitModel() devrait retourner un *InteractiveConfig")
	}

	// Tester la méthode View
	view := ic.View()

	// La vue ne devrait pas être vide
	if view == "" {
		t.Error("View() ne devrait pas retourner une chaîne vide")
	}

	// Simuler une condition de sortie
	ic.quitting = true
	view = ic.View()

	// La vue devrait contenir un message de sortie
	if !strings.Contains(strings.ToLower(view), "config") {
		t.Error("View() devrait contenir une référence à la configuration quand quitting=true")
	}
}

func TestInteractiveConfigUpdate(t *testing.T) {
	// Créer un modèle
	model := InitModel()
	ic, ok := model.(*InteractiveConfig)
	if !ok {
		t.Fatal("InitModel() devrait retourner un *InteractiveConfig")
	}

	// Tester la méthode Update avec différents messages

	// Message de sortie
	model, cmd := ic.Update(tea.QuitMsg{})
	if cmd != nil {
		t.Error("Update(QuitMsg) devrait retourner un cmd nil")
	}

	// Vérifier que le modèle a été mis à jour
	updatedIc, ok := model.(*InteractiveConfig)
	if !ok {
		t.Fatal("Update() devrait retourner un *InteractiveConfig")
	}

	if !updatedIc.quitting {
		t.Error("Update(QuitMsg) devrait définir quitting=true")
	}
}
