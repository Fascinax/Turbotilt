package runtime

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"turbotilt/internal/config"
)

var (
	titleStyle = lipgloss.NewStyle().MarginLeft(2)
	// Remove unused style declarations that are commented out
	paginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle       = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle   = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

// item represents an interactive list element
type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

// InteractiveConfig represents the state of the interactive interface
type InteractiveConfig struct {
	list         list.Model
	textInput    textinput.Model
	config       config.Config
	currentStep  int
	steps        []string
	stepComplete []bool
	quitting     bool
}

// InitModel initializes the model for the interactive interface
func InitModel() tea.Model {
	items := []list.Item{
		item{title: "🔄 Détecter le framework", desc: "Scanner le projet pour identifier le framework"},
		item{title: "⚙️  Configurer le projet", desc: "Définir les paramètres du projet"},
		item{title: "🐳 Configurer Docker", desc: "Paramétrer les options Docker"},
		item{title: "🛠️  Configurer le développement", desc: "Options de développement et live reload"},
		item{title: "🔌 Configurer les services", desc: "Ajouter/configurer les services dépendants"},
		item{title: "📄 Générer les fichiers", desc: "Créer les fichiers de configuration"},
		item{title: "🚀 Démarrer l'environnement", desc: "Lancer l'environnement de développement"},
	}

	l := list.New(items, list.NewDefaultDelegate(), 40, 14)
	l.Title = "🌟 Assistant Turbotilt"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	ti := textinput.New()
	ti.Placeholder = "Entrez votre réponse..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50

	m := &InteractiveConfig{
		list:         l,
		textInput:    ti,
		config:       config.DefaultConfig("unknown"),
		steps:        []string{"select", "project", "docker", "dev", "services", "generate", "start"},
		currentStep:  0,
		stepComplete: make([]bool, 7),
	}

	return m
}

// Init initializes the model
func (m *InteractiveConfig) Init() tea.Cmd {
	return nil
}

// Update updates the model with user events
func (m *InteractiveConfig) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			if m.currentStep == 0 {
				// Main list navigation
				i, ok := m.list.SelectedItem().(item)
				if ok {
					switch i.title {
					case "🔄 Détecter le framework":
						m.currentStep = 1
						m.textInput.SetValue("")
						m.textInput.Placeholder = "Type de framework (spring, quarkus, java, ...)"
						return m, nil
					case "⚙️  Configurer le projet":
						m.currentStep = 2
						m.textInput.SetValue("")
						m.textInput.Placeholder = "Nom du projet"
						return m, nil
					case "🐳 Configurer Docker":
						m.currentStep = 3
						m.textInput.SetValue("")
						m.textInput.Placeholder = "Port à exposer (8080)"
						return m, nil
					case "🛠️  Configurer le développement":
						m.currentStep = 4
						m.textInput.SetValue("")
						m.textInput.Placeholder = "Activer le live reload? (y/n)"
						return m, nil
					case "🔌 Configurer les services":
						m.currentStep = 5
						m.textInput.SetValue("")
						m.textInput.Placeholder = "Services requis (mysql,postgres,redis,kafka) séparés par virgule"
						return m, nil
					case "📄 Générer les fichiers":
						m.currentStep = 6
						m.config.Project.Name = "myapp"
						return m, generateFiles(m)
					case "🚀 Démarrer l'environnement":
						m.currentStep = 7
						return m, startEnvironment(m)
					}
				}
			} else {
				// User input validation
				value := m.textInput.Value()
				switch m.currentStep {
				case 1:
					// Framework
					m.config.Framework.Type = value
					m.stepComplete[0] = true
				case 2:
					// Project
					m.config.Project.Name = value
					m.stepComplete[1] = true
				case 3:
					// Docker
					if value != "" {
						m.config.Docker.Port = value
					}
					m.stepComplete[2] = true
				case 4:
					// Dev options
					m.config.Development.EnableLiveReload = (value == "y" || value == "")
					m.stepComplete[3] = true
				case 5:
					// Services
					// Analyze required services and dependencies
					m.stepComplete[4] = true
				}
				m.currentStep = 0
				return m, nil
			}
		}

	case tea.QuitMsg:
		m.quitting = true
		return m, nil
	// Handle widget updates
	case tea.WindowSizeMsg:
		h, v := lipgloss.NewStyle().Margin(1, 2).GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	if m.currentStep == 0 {
		// List mode
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	} else {
		// Text input mode
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}
}

// generateFiles is a tea.Cmd that generates configuration files
func generateFiles(m *InteractiveConfig) tea.Cmd {
	return func() tea.Msg {
		p := tea.NewProgram(generateFilesModel{config: m.config})
		_, err := p.Run()
		return generateFilesFinishedMsg{err: err}
	}
}

// startEnvironment is a tea.Cmd that starts the development environment
func startEnvironment(m *InteractiveConfig) tea.Cmd {
	return func() tea.Msg {
		p := tea.NewProgram(startEnvironmentModel{})
		_, err := p.Run()
		return startEnvironmentFinishedMsg{err: err}
	}
}

// generateFilesModel is a model for file generation
type generateFilesModel struct {
	config config.Config
}

func (m generateFilesModel) Init() tea.Cmd {
	return nil
}

func (m generateFilesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m generateFilesModel) View() string {
	return "Génération des fichiers..."
}

// generateFilesFinishedMsg is a message indicating file generation completion
type generateFilesFinishedMsg struct {
	err error
}

// startEnvironmentModel is a model for environment startup
type startEnvironmentModel struct{}

func (m startEnvironmentModel) Init() tea.Cmd {
	return nil
}

func (m startEnvironmentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m startEnvironmentModel) View() string {
	return "Démarrage de l'environnement..."
}

// startEnvironmentFinishedMsg is a message indicating environment startup completion
type startEnvironmentFinishedMsg struct {
	err error
}

// View returns the text representation of the model
func (m *InteractiveConfig) View() string {
	if m.quitting {
		return quitTextStyle.Render("Thanks for using Turbotilt! Configuration generated for " + m.config.Project.Name)
	}

	// Display content based on current step
	if m.currentStep == 0 {
		return "\n" + m.list.View()
	} else {
		var title string

		switch m.currentStep {
		case 1:
			title = "🔄 Configuration du Framework"
		case 2:
			title = "⚙️ Configuration du Projet"
		case 3:
			title = "🐳 Configuration Docker"
		case 4:
			title = "🛠️ Configuration du Développement"
		case 5:
			title = "🔌 Configuration des Services"
		}

		return fmt.Sprintf("\n%s\n\n%s\n\n%s",
			titleStyle.Render(title),
			m.textInput.View(),
			"(Entrer pour valider, ESC pour annuler)")
	}
}

// Wizard launches the interactive setup assistant
func Wizard() {
	p := tea.NewProgram(InitModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running the assistant: %v\n", err)
		os.Exit(1)
	}
}
