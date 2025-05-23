package model

type Command struct {
	Type  CommandType
	Var   *Variable
	Op    Operation
	Left  Argument
	Right Argument
}

const (
	Print CommandType = "print"
	Calc  CommandType = "calc"
)

type CommandType string

func IsValidCommand(command CommandType) bool {
	switch command {
	case Print, Calc:
		return true
	default:
		return false
	}
}

func (c *Command) IsCalc() bool {
	return c.Type == Calc
}

func (c *Command) IsPrint() bool {
	return c.Type == Print
}
