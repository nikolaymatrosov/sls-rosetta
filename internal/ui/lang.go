package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/nikolaymatrosov/sls-rosetta/internal/examples"
)

type langItem struct {
	title string
	value string
}

func (i langItem) Description() string {
	return ""
}

func (i langItem) Title() string {
	return i.title
}

func (i langItem) Value() string {
	return i.value
}

func (i langItem) String() string { return i.title }

func (i langItem) FilterValue() string { return i.title }

func NewLangList(languages []examples.Language) list.Model {
	var langItems []list.Item
	for _, lang := range languages {
		langItems = append(langItems, langItem{
			title: lang.Title,
			value: lang.Name,
		})
	}
	dd := list.NewDefaultDelegate()
	dd.ShowDescription = false

	l := list.New(langItems, dd, 0, 0)
	l.Title = "Select function language"
	l.SetShowStatusBar(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	return l
}
