package ui

import (
	"github.com/charmbracelet/bubbles/list"
)

type deployItem struct {
	title       string
	description string
	value       string
}

func (i deployItem) Title() string {
	return i.title
}

func (i deployItem) Description() string {
	return i.description
}

func (i deployItem) Value() string {
	return i.value
}

func (i deployItem) String() string { return i.title }

func (i deployItem) FilterValue() string { return i.title }

func NewDeployList(item langItem, ex exampleItem) list.Model {
	deployOptions := []list.Item{
		deployItem{
			title:       "Terraform",
			description: "Add terraform to your project",
			value:       "terraform",
		},
		deployItem{
			title:       "YC CLI",
			description: "Add Makefile with YC CLI commands",
			value:       "yccli",
		},
		deployItem{
			title:       "None",
			description: "Do not add anything",
			value:       "none",
		},
	}

	l := list.New(deployOptions, descriptedItemDelegate{}, defaultWidth, listHeight)
	l.Title = "Select way to deploy your function"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	return l
}
