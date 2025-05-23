package model

const (
	Plus = iota
	Minus
	Multiply
)

type Operation uint8

func IsValidOperationBySymbol(symbol string) bool {
	switch symbol {
	case "+", "-", "*":
		return true
	default:
		return false
	}
}

func GetOperationBySymbol(symbol string) Operation {
	switch symbol {
	case "+":
		return Plus
	case "-":
		return Minus
	case "*":
		return Multiply
	}

	return 0
}
