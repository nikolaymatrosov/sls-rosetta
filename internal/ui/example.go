package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/nikolaymatrosov/sls-rosetta/internal/examples"
)

type exampleItem struct {
	title         string
	description   string
	value         string
	DeployOptions []examples.Deploy
}

func (i exampleItem) Title() string {
	return i.title
}

func (i exampleItem) Description() string {
	return i.description
}

func (i exampleItem) Value() string {
	return i.value
}

func (i exampleItem) String() string { return i.title }

func (i exampleItem) FilterValue() string { return i.title }

func NewExampleList(exs []examples.Example) list.Model {
	var exampleItems []list.Item
	for _, example := range exs {
		exampleItems = append(exampleItems, exampleItem{
			title:         example.Title,
			description:   example.Description,
			value:         example.Name,
			DeployOptions: example.Deploy,
		})
	}

	dd := list.NewDefaultDelegate()

	l := list.New(exampleItems, dd, 0, 0)
	l.Title = "Select function type"
	l.SetShowStatusBar(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	return l
}
