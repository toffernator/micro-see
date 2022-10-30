package stack_machine

import (
	"errors"
	"fmt"
)

const (
	STACK_SIZE = 20
)

type StackMachine struct {
	Bp         int
	Sp         int
	Pc         int
	Stack      []int
	Prog       []int
	hasStopped bool
	Res        int
}

func New() *StackMachine {
	return &StackMachine{Bp: -999, Sp: -1, Pc: 0, Stack: make([]int, STACK_SIZE), Prog: []int{STOP}, hasStopped: false, Res: 0}
}

func (s *StackMachine) Load(prog []int) {
	s.Prog = prog
}

func (s *StackMachine) reset() {
	s.Bp = -999
	s.Sp = -1
	s.Pc = 0
	s.hasStopped = false
	s.Res = 0
}

func (s *StackMachine) Exec(iargs []int) (int, error) {
	s.reset()
	for !s.hasStopped {
		if err := s.execInstruction(iargs); err != nil {
			return -1, err
		}
		s.Pc++
	}
	return s.Res, nil
}

func (s *StackMachine) ExecDebug(iargs []int, commands chan int) (int, error) {
	s.reset()
	for !s.hasStopped {
		<-commands
		if err := s.execInstruction(iargs); err != nil {
			return -1, err
		}
		s.Pc++
	}
	return s.Res, nil
}

func (s *StackMachine) execInstruction(iargs []int) error {
	switch s.Prog[s.Pc] {
	case CSTI:
		s.Pc++
		s.Stack[s.Sp+1] = s.Prog[s.Pc]
		s.Sp++
	case ADD:
		s.Stack[s.Sp-1] = s.Stack[s.Sp-1] + s.Stack[s.Sp]
		s.Sp--
	case SUB:
		s.Stack[s.Sp-1] = s.Stack[s.Sp-1] - s.Stack[s.Sp]
		s.Sp--
	case MUL:
		s.Stack[s.Sp-1] = s.Stack[s.Sp-1] * s.Stack[s.Sp]
		s.Sp--
	case DIV:
		s.Stack[s.Sp-1] = s.Stack[s.Sp-1] / s.Stack[s.Sp]
		s.Sp--
	case MOD:
		s.Stack[s.Sp-1] = s.Stack[s.Sp-1] % s.Stack[s.Sp]
		s.Sp--
	case EQ:
		if s.Stack[s.Sp-1] == s.Stack[s.Sp] {
			s.Stack[s.Sp-1] = 1
		} else {
			s.Stack[s.Sp-1] = 0
		}
		s.Sp--
	case LT:
		if s.Stack[s.Sp-1] < s.Stack[s.Sp] {
			s.Stack[s.Sp-1] = 1
		} else {
			s.Stack[s.Sp-1] = 0
		}
		s.Sp--
	case NOT:
		if s.Stack[s.Sp] == 0 {
			s.Stack[s.Sp] = 1
		} else {
			s.Stack[s.Sp] = 0
		}
	case DUP:
		s.Stack[s.Sp+1] = s.Stack[s.Sp]
		s.Sp++
	case SWAP:
		tmp := s.Stack[s.Sp]
		s.Stack[s.Sp] = s.Stack[s.Sp-1]
		s.Stack[s.Sp-1] = tmp
	case LDI:
		s.Stack[s.Sp] = s.Stack[s.Stack[s.Sp]]
	case STI:
		s.Stack[s.Stack[s.Sp-1]] = s.Stack[s.Sp]
		s.Stack[s.Sp-1] = s.Stack[s.Sp]
		s.Sp--
	case GETBP:
		s.Stack[s.Sp+1] = s.Bp
		s.Sp++
	case GETSP:
		s.Stack[s.Sp+1] = s.Sp
		s.Sp++
	case INCSP:
		s.Pc++
		s.Sp = s.Sp + s.Prog[s.Pc]
	case GOTO:
		s.Pc = s.Prog[s.Pc]
	case IFZERO:
		s.Sp--
		if s.Stack[s.Sp] == 0 {
			s.Pc = s.Prog[s.Pc]
		} else {
			s.Pc++
		}
	case IFNZRO:
		s.Sp--
		if s.Stack[s.Sp] != 0 {
			s.Pc = s.Prog[s.Pc]
		} else {
			s.Pc++
		}
	case CALL:
		s.Pc++
		argc := s.Prog[s.Pc]
		for i := 0; i < argc; i++ {
			s.Stack[s.Sp-i+2] = s.Stack[s.Sp-i]
		}
		s.Stack[s.Sp-argc+1] = s.Pc + 1
		s.Sp++
		s.Stack[s.Sp-argc+1] = s.Bp
		s.Sp++
		s.Bp = s.Sp + 1 - argc
		s.Pc = s.Prog[s.Pc]
	case TCALL:
		s.Pc++
		argc := s.Prog[s.Pc]
		s.Pc++
		pop := s.Prog[s.Pc]
		for i := argc - 1; i >= 0; i-- {
			s.Stack[s.Sp-i-pop] = s.Stack[s.Sp-i]
		}
		s.Sp = s.Sp - pop
		s.Pc = s.Prog[s.Pc]
	case RET:
		res := s.Stack[s.Sp]
		s.Sp = s.Sp - s.Prog[s.Pc]
		s.Bp = s.Stack[s.Sp]
		s.Sp--
		s.Pc = s.Stack[s.Sp]
		s.Sp--
		s.Stack[s.Sp] = res
	case PRINTI:
		fmt.Printf("%d ", s.Stack[s.Sp])
	case PRINTC:
		fmt.Printf("%c ", s.Stack[s.Sp])
	case LDARGS:
		for i := 0; i < len(iargs); i++ {
			s.Stack[s.Sp] = iargs[i]
			s.Sp++
		}
	case STOP:
		s.Res = s.Stack[s.Sp]
		s.hasStopped = true
	default:
		errString := fmt.Sprintf("illegal instruction %d at address %d", s.Prog[s.Pc-1], s.Pc-1)
		return errors.New(errString)
	}
	return nil
}
