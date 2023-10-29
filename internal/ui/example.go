package ui

import (
	"github.com/charmbracelet/bubbles/list"
)

type exampleItem struct {
	title       string
	description string
	value       string
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

func NewExampleList(item langItem) list.Model {
	examples := []list.Item{
		exampleItem{
			title:       "HTTP",
			description: "Simple HTTP handler",
			value:       "http",
		},
		exampleItem{
			title:       "SQS",
			description: "SQS event handler",
			value:       "sqs",
		},
		exampleItem{
			title:       "API Gateway",
			description: "API Gateway handler",
			value:       "apigateway",
		},
		exampleItem{
			title:       "S3",
			description: "S3 event handler",
			value:       "s3",
		},
		exampleItem{
			title:       "Billing",
			description: "Billing event handler",
			value:       "billing",
		},
	}

	l := list.New(examples, descriptedItemDelegate{}, defaultWidth, listHeight)
	l.Title = "Select function type"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	return l
}
