package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/nikolaymatrosov/sls-rosetta/internal/examples"
)

type deployItem struct {
	title       string
	description string
	value       string
	exclusive   []string
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

func constructDeployItem(value examples.Deploy) *deployItem {
	switch value.Type {
	case "terraform":
		return &deployItem{
			title:       "Terraform",
			description: "Add terraform to your project",
			value:       "terraform",
			exclusive:   value.Exclusive,
		}
	case "yccli":
		return &deployItem{
			title:       "YC CLI",
			description: "Add Makefile with YC CLI commands",
			value:       "yccli",
			exclusive:   value.Exclusive,
		}
	case "none":
		return &deployItem{
			title:       "None",
			description: "Do not add anything",
			value:       "none",
			exclusive:   []string{},
		}
	default:
		return nil
	}
}

func NewDeployList(deployOptions []examples.Deploy) list.Model {
	var deployItems []list.Item
	for _, deployOption := range deployOptions {
		di := constructDeployItem(deployOption)
		if di == nil {
			continue
		}
		deployItems = append(deployItems, *di)
	}

	dd := list.NewDefaultDelegate()

	dd.ShowDescription = false

	l := list.New(deployItems, dd, 0, 0)
	l.Title = "Select way to deploy your function"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	return l
}
