package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/toffernator/micro-see/model"
)

func main() {
	p := tea.NewProgram(model.New())
	p.Start()
}
