package model

type Argument interface {
	GetValue() int64
	HasDependency() bool
}
