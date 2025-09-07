package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"github.com/afonso-borges/t-hub/internal/themes"
	"github.com/afonso-borges/t-hub/internal/utils"
)

const maxWidth = 80

var (
	red    = lipgloss.AdaptiveColor{Light: "#FE5F86", Dark: "#FE5F86"}
	indigo = lipgloss.AdaptiveColor{Light: "#7A34BB", Dark: "#7A34BB"}
	green  = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}
)

type Styles struct {
	Base,
	HeaderText,
	Status,
	StatusHeader,
	Highlight,
	ErrorHeaderText,
	Help lipgloss.Style
}

func NewStyles(lg *lipgloss.Renderer) *Styles {
	s := Styles{}
	s.Base = lg.NewStyle().
		Padding(1, 4, 0, 1)
	s.HeaderText = lg.NewStyle().
		Foreground(indigo).
		Bold(true).
		Padding(0, 1, 0, 2)
	s.Status = lg.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(indigo).
		PaddingLeft(1).
		MarginTop(1)
	s.StatusHeader = lg.NewStyle().
		Foreground(green).
		Bold(true)
	s.Highlight = lg.NewStyle().
		Foreground(lipgloss.Color("212"))
	s.ErrorHeaderText = s.HeaderText.
		Foreground(red)
	s.Help = lg.NewStyle().
		Foreground(lipgloss.Color("240"))
	return &s
}

type state int

const (
	stateWelcome state = iota
	stateLoading
	statePlayerRemoval
	stateResults
	stateStartOver
	stateDone
)

type Model struct {
	state           state
	lg              *lipgloss.Renderer
	styles          *Styles
	form            *huh.Form
	width           int
	height          int
	playersToRemove []string
	analyzer        string
	players         []utils.Player
	split           utils.GoldSplit
	loading         bool
	spinner         spinner.Model
}

func NewModel() Model {
	m := Model{
		width: maxWidth,
		state: stateWelcome,
	}
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)

	// Initialize spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(indigo)
	m.spinner = s

	m.createWelcomeForm()
	return m
}

func (m *Model) createWelcomeForm() {
	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Welcome to T-HUB").
				Description("Loot split calculator for your terminal").
				Next(true).
				NextLabel("Start"),
		),
	).
		WithWidth(50).
		WithShowHelp(false).
		WithShowErrors(false)
}

func (m *Model) createPlayerRemovalForm() {
	if len(m.players) == 0 {
		return
	}

	playerOptions := utils.ExtractPlayerNames(m.players)

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Remove players from loot split?").
				Description("Select players to exclude from the calculation").
				Value(&m.playersToRemove).
				Options(playerOptions...),
		),
	).
		WithWidth(50).
		WithShowHelp(false).
		WithShowErrors(false).
		WithTheme(themes.DefaultTheme())
}

func (m *Model) createResultsForm() {
	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Description(utils.FormatTransfers(m.split)).
				Next(true).
				NextLabel("Copy to clipboard"),
		),
	).
		WithWidth(50).
		WithShowHelp(false).
		WithShowErrors(false)
}

func (m *Model) createStartOverForm() {
	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Start over?").
				Description("Do you want to calculate another loot split?").
				Affirmative("Yes").
				Negative("No, exit"),
		),
	).
		WithWidth(50).
		WithShowHelp(false).
		WithShowErrors(false)
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.form.Init(), m.spinner.Tick)
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

type analyzerLoadedMsg struct {
	analyzer string
	players  []utils.Player
	err      error
}

func loadAnalyzer() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(1000 * time.Millisecond)

		analyzer, err := utils.CopyFromClipboard()
		if err != nil {
			return analyzerLoadedMsg{err: err}
		}

		_, players, err := utils.ParseAnalyzer(analyzer)
		if err != nil {
			return analyzerLoadedMsg{err: err}
		}

		return analyzerLoadedMsg{
			analyzer: analyzer,
			players:  players,
		}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = min(msg.Width, maxWidth) - m.styles.Base.GetHorizontalFrameSize()
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Interrupt
		case "esc", "q":
			return m, tea.Quit
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case analyzerLoadedMsg:
		if msg.err != nil {
			log.Fatal(msg.err)
		}
		m.analyzer = msg.analyzer
		m.players = msg.players
		m.state = statePlayerRemoval
		m.loading = false
		m.createPlayerRemovalForm()
		return m, m.form.Init()
	}

	var cmds []tea.Cmd

	// Process the form
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	if m.form.State == huh.StateCompleted {
		switch m.state {
		case stateWelcome:
			m.state = stateLoading
			m.loading = true
			return m, loadAnalyzer()
		case statePlayerRemoval:
			remainingPlayers := utils.FilterRemainingPlayers(m.players, m.playersToRemove)
			m.split = utils.CalculateGoldSplit(remainingPlayers)
			m.state = stateResults
			m.createResultsForm()
			return m, m.form.Init()
		case stateResults:
			utils.SaveToClipboard(m.split)
			m.state = stateStartOver
			m.createStartOverForm()
			return m, m.form.Init()
		case stateStartOver:
			if m.form.GetBool("") { // User confirmed "Yes"
				// Reset state and start over
				m.playersToRemove = []string{}
				m.analyzer = ""
				m.players = []utils.Player{}
				m.split = utils.GoldSplit{}
				m.state = stateLoading
				m.loading = true
				return m, loadAnalyzer()
			} else {
				return m, tea.Quit
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	s := m.styles

	// Get header and footer content
	var headerText string
	var footerText string

	if m.loading {
		headerText = "T-HUB - Loot Split Calculator"
		footerText = ""
	} else {
		errors := m.form.Errors()
		if len(errors) > 0 {
			headerText = m.errorView()
			footerText = ""
		} else {
			switch m.state {
			case stateWelcome:
				headerText = "T-HUB - Loot Split Calculator"
			case statePlayerRemoval:
				headerText = "T-HUB - Player Selection"
			case stateResults:
				headerText = "T-HUB - Results"
			case stateStartOver:
				headerText = "T-HUB - Start Over"
			default:
				headerText = "T-HUB - Loot Split Calculator"
			}
			footerText = m.form.Help().ShortHelpView(m.form.KeyBinds())
		}
	}

	// Create header and footer
	var header, footer string
	if len(m.form.Errors()) > 0 {
		header = m.appErrorBoundaryView(headerText)
		footer = m.appErrorBoundaryView(footerText)
	} else {
		header = m.appBoundaryView(headerText)
		footer = m.appBoundaryView(footerText)
	}

	// Calculate available height for content (subtract header and footer lines)
	headerHeight := lipgloss.Height(header)
	footerHeight := lipgloss.Height(footer)
	contentHeight := m.height - headerHeight - footerHeight

	// Create main content
	var content string
	if m.loading {
		spinnerText := fmt.Sprintf("%s Processing analyzer...", m.spinner.View())
		centeredLoading := s.Status.
			Width(50).
			Padding(2).
			Render(spinnerText)
		content = lipgloss.Place(m.width, contentHeight, lipgloss.Center, lipgloss.Center, centeredLoading)
	} else {
		// Form (centered)
		v := strings.TrimSuffix(m.form.View(), "\n\n")
		form := s.Status.
			Width(60).
			Padding(2).
			Render(v)
		content = lipgloss.Place(m.width, contentHeight, lipgloss.Center, lipgloss.Center, form)
	}

	// Combine all parts with proper spacing
	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		content,
		footer,
	)
}

func (m Model) errorView() string {
	var s string
	for _, err := range m.form.Errors() {
		s += err.Error()
	}
	return s
}

func (m Model) appBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		m.styles.HeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(indigo),
	)
}

func (m Model) appErrorBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		m.styles.ErrorHeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(red),
	)
}

func main() {
	_, err := tea.NewProgram(NewModel(), tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Println("Oh no:", err)
		os.Exit(1)
	}
}
