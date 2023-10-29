package main

import (
	"fmt"
	"os"

	"github.com/nikolaymatrosov/sls-rosetta/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m := ui.NewViewModel()
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
