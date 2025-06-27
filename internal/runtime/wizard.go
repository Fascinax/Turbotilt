package runtime

import (
	"fmt"
	"os"
	"turbotilt/internal/config"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

// item représente un élément de la liste interactive
type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

// InteractiveConfig représente l'état de l'interface interactive
type InteractiveConfig struct {
	list         list.Model
	textInput    textinput.Model
	config       config.Config
	currentStep  int
	steps        []string
	stepComplete []bool
	quitting     bool
}

// InitModel initialise le modèle pour l'interface interactive
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

// Init initialise le modèle
func (m *InteractiveConfig) Init() tea.Cmd {
	return nil
}

// Update met à jour le modèle avec les événements utilisateur
func (m *InteractiveConfig) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			if m.currentStep == 0 {
				// Navigation dans la liste principale
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
				// Validation des entrées utilisateur
				value := m.textInput.Value()
				switch m.currentStep {
				case 1:
					// Framework
					m.config.Framework.Type = value
					m.stepComplete[0] = true
				case 2:
					// Projet
					m.config.Project.Name = value
					m.stepComplete[1] = true
				case 3:
					// Docker
					if value != "" {
						m.config.Docker.Port = value
					}
					m.stepComplete[2] = true
				case 4:
					// Dev
					m.config.Development.EnableLiveReload = (value == "y" || value == "")
					m.stepComplete[3] = true
				case 5:
					// Services
					// TODO: Analyse des services requis
					m.stepComplete[4] = true
				}
				m.currentStep = 0
				return m, nil
			}
		}

	case tea.QuitMsg:
		m.quitting = true
		return m, nil
	// Gérer les mises à jour des widgets
	case tea.WindowSizeMsg:
		h, v := lipgloss.NewStyle().Margin(1, 2).GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	if m.currentStep == 0 {
		// Mode liste
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	} else {
		// Mode saisie texte
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}
}

// generateFiles est un tea.Cmd qui génère les fichiers
func generateFiles(m *InteractiveConfig) tea.Cmd {
	return func() tea.Msg {
		p := tea.NewProgram(generateFilesModel{config: m.config})
		err := p.Start()
		return generateFilesFinishedMsg{err: err}
	}
}

// startEnvironment est un tea.Cmd qui démarre l'environnement
func startEnvironment(m *InteractiveConfig) tea.Cmd {
	return func() tea.Msg {
		p := tea.NewProgram(startEnvironmentModel{})
		err := p.Start()
		return startEnvironmentFinishedMsg{err: err}
	}
}

// generateFilesModel est un modèle pour la génération des fichiers
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

// generateFilesFinishedMsg est un message signalant la fin de la génération des fichiers
type generateFilesFinishedMsg struct {
	err error
}

// startEnvironmentModel est un modèle pour le démarrage de l'environnement
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

// startEnvironmentFinishedMsg est un message signalant la fin du démarrage
type startEnvironmentFinishedMsg struct {
	err error
}

// View renvoie la représentation textuelle du modèle
func (m *InteractiveConfig) View() string {
	if m.quitting {
		return quitTextStyle.Render("Merci d'avoir utilisé Turbotilt! Configuration générée pour " + m.config.Project.Name)
	}

	// Afficher le contenu selon l'étape actuelle
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

// Wizard lance l'assistant interactif
func Wizard() {
	p := tea.NewProgram(InitModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Erreur lors de l'exécution de l'assistant: %v\n", err)
		os.Exit(1)
	}
}
