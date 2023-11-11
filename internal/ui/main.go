package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nikolaymatrosov/sls-rosetta/internal/cloner"
	"github.com/nikolaymatrosov/sls-rosetta/internal/examples"
)

// const listHeight = 14
const defaultWidth = 50
const purpleColor = "#874bfc"
const redColor = "#ff5555"

var (
	titleStyle          = lipgloss.NewStyle().MarginLeft(2)
	itemStyle           = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle   = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	descriptionStyle    = lipgloss.NewStyle().PaddingLeft(4).Foreground(lipgloss.Color("241"))
	paginationStyle     = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle           = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	selectedOptionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(purpleColor))
	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(redColor))

	checkMark = selectedOptionStyle.Render("✓")
)

type uiState uint

const (
	languageListView uiState = iota
	exampleListView
	deployListView
	pathView
	resultView
)

type model struct {
	config     *examples.Config
	activeList list.Model
	width      int
	height     int

	language   langItem
	example    exampleItem
	deployType deployItem

	globesToExclude []string

	pathInput TextInputView
	clonePath string

	state    uiState
	banner   []string
	quitting bool
}

func (m *model) selectLanguage() {
	m.language = m.activeList.SelectedItem().(langItem)
	m.activeList.Title = fmt.Sprintf("Select %s example", m.language.title)
	m.banner = append(m.banner, fmt.Sprintf("%s Language: %s", checkMark, m.language.title))
	exs, ok := m.config.Examples[m.language.value]
	if !ok {
		fmt.Printf("No examples for language %s\n", m.language.value)
		os.Exit(1)
	}
	m.activeList = NewExampleList(exs)
	m.state = exampleListView
}

func (m *model) selectExample() {
	m.example = m.activeList.SelectedItem().(exampleItem)
	m.banner = append(m.banner, fmt.Sprintf("%s Example: %s", checkMark, m.example.title))
	m.pathInput.Focus()
	// find selected example in config
	var deployOptions []examples.Deploy
	for _, ex := range m.config.Examples[m.language.value] {
		if ex.Name == m.example.value {
			deployOptions = ex.Deploy
			break
		}
	}

	m.activeList = NewDeployList(deployOptions)
	m.state = deployListView
}

func (m *model) selectDeployType() {
	m.deployType = m.activeList.SelectedItem().(deployItem)

	var globesToExclude []string
	for _, d := range m.example.DeployOptions {
		if d.Type != m.deployType.value {
			globesToExclude = append(globesToExclude, d.Exclusive...)
			break
		}
	}
	m.globesToExclude = globesToExclude

	m.banner = append(m.banner, fmt.Sprintf("%s Deploy type: %s", checkMark, m.deployType.title))
	m.state = pathView
}

func (m *model) selectClonePath() error {
	m.clonePath = m.pathInput.Value()
	if m.clonePath == "" {
		m.pathInput.Title = errorStyle.Render("Path can't be empty")
		m.pathInput.Focus()
		return fmt.Errorf("path can't be empty")
	}

	err := cloner.CheckThatPathDoesntExist(m.clonePath)
	if err != nil {
		m.pathInput.Title = errorStyle.Render(err.Error())
		m.pathInput.Focus()
		return err
	}
	m.banner = append(m.banner, fmt.Sprintf("%s Path: %s", checkMark, m.clonePath))
	m.state = resultView
	return nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateListSize()
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			switch m.state {
			case languageListView:
				m.selectLanguage()
			case exampleListView:
				m.selectExample()
			case deployListView:
				m.selectDeployType()
			case pathView:
				err := m.selectClonePath()
				if err == nil {
					c := cloner.NewCloner(
						m.config.Repo,
						fmt.Sprintf("examples/%s/%s", m.language.value, m.example.value),
						m.clonePath,
						m.globesToExclude,
					)
					c.Clone("")
					return m, tea.Quit
				}
			}
			m.updateListSize()
		}

	}
	switch m.state {
	// update whichever model is focused
	case languageListView, exampleListView, deployListView:
		m.activeList, cmd = m.activeList.Update(msg)
		cmds = append(cmds, cmd)
	case pathView:
		m.pathInput, cmd = m.pathInput.Update(msg)
		cmds = append(cmds, cmd)
	default:
		return m, nil
	}
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	s := strings.Join(m.banner, "\n") + "\n"
	switch m.state {
	case languageListView, exampleListView, deployListView:
		s += m.activeList.View()
	case pathView:
		s += m.pathInput.View()
	default:
		return s
	}

	return s
}

func (m *model) updateListSize() {
	bannerH := len(m.banner)
	m.activeList.SetSize(m.width, m.height-bannerH)
}

func NewViewModel(
	config *examples.Config,
) tea.Model {

	l := NewLangList(config.Languages)

	cp := NewClonePathView()

	m := model{
		state:      languageListView,
		config:     config,
		activeList: l,
		pathInput:  cp,
	}
	return m
}
