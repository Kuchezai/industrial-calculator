package model

type NumericArgument int64

func (n NumericArgument) GetValue() int64 {
	return int64(n)
}

func (n NumericArgument) HasDependency() bool {
	return false
}
