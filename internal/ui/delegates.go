package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbletea"
)

type SimpleItem interface {
	Title() string
	Value() string
}

type simpleItemDelegate struct{}

func (d simpleItemDelegate) Height() int                             { return 1 }
func (d simpleItemDelegate) Spacing() int                            { return 0 }
func (d simpleItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d simpleItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(SimpleItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.Title())

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	_, err := fmt.Fprint(w, fn(str))
	if err != nil {
		return
	}
}

type DesriptedItem interface {
	Title() string
	Description() string
	Value() string
}

type descriptedItemDelegate struct{}

func (d descriptedItemDelegate) Height() int                             { return 2 }
func (d descriptedItemDelegate) Spacing() int                            { return 0 }
func (d descriptedItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d descriptedItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(DesriptedItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%2d. %s\n%s", index+1, i.Title(), descriptionStyle.Render(i.Description()))

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			res := strings.Split(strings.Join(s, " "), "\n")
			return selectedItemStyle.Render("> " + res[0] + "\n  " + res[1])
		}
	}

	_, err := fmt.Fprint(w, fn(str))
	if err != nil {
		return
	}
}
