package ui

import "github.com/charmbracelet/bubbles/list"

type langItem struct {
	title string
	value string
}

func (i langItem) Title() string {
	return i.title
}

func (i langItem) Value() string {
	return i.value
}

func (i langItem) String() string { return i.title }

func (i langItem) FilterValue() string { return i.title }

func NewLangList() list.Model {
	languages := []list.Item{
		langItem{
			title: "TypeScript",
			value: "typescript",
		},
		langItem{
			title: "Go",
			value: "go",
		},
		langItem{
			title: "C#",
			value: "csharp",
		},
		langItem{
			title: "Java",
			value: "java",
		},
		langItem{
			title: "Python",
			value: "python",
		},
		langItem{
			title: "PHP",
			value: "php",
		},
	}

	l := list.New(languages, simpleItemDelegate{}, defaultWidth, listHeight)
	l.Title = "Select function language"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	return l
}
