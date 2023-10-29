package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func newClonePathInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "Path to clone"
	ti.CharLimit = 156
	ti.Width = defaultWidth

	return ti
}

type TextInputView struct {
	textInput textinput.Model
	Title     string
}

func (t TextInputView) Init() tea.Cmd {
	return textinput.Blink
}

func (t *TextInputView) Focus() tea.Cmd {
	return t.textInput.Focus()
}

func (t TextInputView) Blur() {
	t.textInput.Blur()
}

func (t TextInputView) Value() string {
	return t.textInput.Value()
}

func (t TextInputView) Update(msg tea.Msg) (TextInputView, tea.Cmd) {
	var cmd tea.Cmd
	t.textInput, cmd = t.textInput.Update(msg)
	return t, cmd
}

func (t TextInputView) View() string {
	return titleStyle.Render(t.Title) + "\n\n" + t.textInput.View()
}

func NewClonePathView() TextInputView {
	tv := TextInputView{
		textInput: newClonePathInput(),
		Title:     "Enter path to clone",
	}

	return tv
}
