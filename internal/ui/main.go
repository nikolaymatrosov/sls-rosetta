package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 14
const defaultWidth = 50
const purpleColor = "#874bfc"

var (
	titleStyle          = lipgloss.NewStyle().MarginLeft(2)
	itemStyle           = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle   = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	descriptionStyle    = lipgloss.NewStyle().PaddingLeft(4).Foreground(lipgloss.Color("241"))
	paginationStyle     = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle           = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle       = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	selectedOptionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(purpleColor))

	checkMark = selectedOptionStyle.Render("✓")
)

type uiState uint

const (
	languageView uiState = iota
	exampleView
	pathView
	deployView
	resultView
)

type model struct {
	languages list.Model
	language  langItem

	examples list.Model
	example  exampleItem

	deploySelector list.Model
	deployType     deployItem

	pathInput TextInputView
	clonePath string

	state    uiState
	banner   []string
	quitting bool
}

func (m *model) selectLanguage() {
	m.language = m.languages.SelectedItem().(langItem)
	m.examples.Title = fmt.Sprintf("Select %s example", m.language.title)
	m.banner = append(m.banner, fmt.Sprintf("%s Language: %s", checkMark, m.language.title))
	m.state = exampleView
}

func (m *model) selectExample() {
	m.example = m.examples.SelectedItem().(exampleItem)
	m.banner = append(m.banner, fmt.Sprintf("%s Example: %s", checkMark, m.example.title))
	m.pathInput.Focus()
	m.state = deployView
}

func (m *model) selectDeployType() {
	m.deployType = m.deploySelector.SelectedItem().(deployItem)
	m.banner = append(m.banner, fmt.Sprintf("%s Deploy type: %s", checkMark, m.deployType.title))
	m.state = pathView
}

func (m *model) selectClonePath() {
	m.clonePath = m.pathInput.Value()
	m.banner = append(m.banner, fmt.Sprintf("%s Path: %s", checkMark, m.clonePath))
	m.state = resultView
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.languages.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			switch m.state {
			case languageView:
				m.selectLanguage()
			case exampleView:
				m.selectExample()
			case deployView:
				m.selectDeployType()
			case pathView:
				m.selectClonePath()
				return m, tea.Quit
			}
		}
	}
	switch m.state {
	// update whichever model is focused
	case languageView:
		m.languages, cmd = m.languages.Update(msg)
		cmds = append(cmds, cmd)
	case exampleView:
		m.examples, cmd = m.examples.Update(msg)
		cmds = append(cmds, cmd)
	case pathView:
		m.pathInput, cmd = m.pathInput.Update(msg)
		cmds = append(cmds, cmd)
	case deployView:
		m.deploySelector, cmd = m.deploySelector.Update(msg)
		cmds = append(cmds, cmd)
	default:
		return m, nil
	}
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	s := strings.Join(m.banner, "\n") + "\n"
	switch m.state {
	case languageView:
		s += m.languages.View()
	case exampleView:
		s += m.examples.View()
	case pathView:
		s += m.pathInput.View()
	case deployView:
		s += m.deploySelector.View()
	default:
		return s
	}

	return s
}

func NewViewModel() tea.Model {
	l := NewLangList()

	e := NewExampleList(langItem{})
	cp := NewClonePathView()

	m := model{
		state:          languageView,
		languages:      l,
		examples:       e,
		deploySelector: NewDeployList(langItem{}, exampleItem{}),
		pathInput:      cp,
	}
	return m
}
