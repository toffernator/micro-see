package model

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/toffernator/micro-see/stack_machine"
)

const (
	TABLE_WIDTH = 26

	UPPER_LEFT_CORNER_BORDER  = "╔"
	HORIZONTAL_BORDER         = "═"
	UPPER_RIGHT_CORNER_BORDER = "╗"
	VERTICAL_BORDER           = "║"
	LOWER_RIGHT_CORNER_BORDER = "╝"
	LOWER_LEFT_CORNER_BORDER  = "╚"
)

type StackMachineModel struct {
	stack_machine.StackMachine
	debugCommands chan int
	isRunning     bool
}

func New() *StackMachineModel {
	return &StackMachineModel{StackMachine: *stack_machine.New(), debugCommands: make(chan int)}
}

// Outputs the stack in a table like format:
// ╔════════════════════════╗
// ║ Addr ║ Value ║ SP ║ BP ║
// ║════════════════════════║
// ║ 0    ║ 0     ║ x  ║ x  ║
// ╚════════════════════════╝
func (s *StackMachineModel) View() (view string) {
	view += UPPER_LEFT_CORNER_BORDER + strings.Repeat(HORIZONTAL_BORDER, TABLE_WIDTH-2) + UPPER_RIGHT_CORNER_BORDER + "\n"
	view += VERTICAL_BORDER + " ADDR " + VERTICAL_BORDER + " VALUE " + VERTICAL_BORDER + " SP " + VERTICAL_BORDER + " BP " + VERTICAL_BORDER + "\n"
	view += VERTICAL_BORDER + strings.Repeat(HORIZONTAL_BORDER, TABLE_WIDTH-2) + VERTICAL_BORDER + "\n"
	for addr, value := range s.Stack {
		view += VERTICAL_BORDER + fmt.Sprintf(" %d    ", addr)
		view += VERTICAL_BORDER + fmt.Sprintf(" %d     ", value)
		if s.Sp == addr {
			view += VERTICAL_BORDER + " x  "
		} else {
			view += VERTICAL_BORDER + "    "
		}
		if s.Bp == addr {
			view += VERTICAL_BORDER + " x  "
		} else {
			view += VERTICAL_BORDER + "    "
		}
		view += VERTICAL_BORDER + "\n"
		view += VERTICAL_BORDER + strings.Repeat(HORIZONTAL_BORDER, TABLE_WIDTH-2) + VERTICAL_BORDER + "\n"
	}
	view += LOWER_LEFT_CORNER_BORDER + strings.Repeat(HORIZONTAL_BORDER, TABLE_WIDTH-2) + LOWER_RIGHT_CORNER_BORDER
	return view
}

func (s *StackMachineModel) Init() tea.Cmd {
	return nil
}

func (s *StackMachineModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return s, tea.Quit
		case "r":
			prog := []int{stack_machine.LDARGS, stack_machine.CSTI, 1, stack_machine.STOP}
			s.Load(prog)
			s.Exec(make([]int, 0))
		case "d":
			prog := []int{stack_machine.LDARGS, stack_machine.CSTI, 1, stack_machine.CSTI, 1, stack_machine.STOP}
			s.Load(prog)
			go s.ExecDebug(make([]int, 0), s.debugCommands)
		case "n":
			select {
			case s.debugCommands <- 1:
			default:
			}
		}
	}
	return s, nil
}
