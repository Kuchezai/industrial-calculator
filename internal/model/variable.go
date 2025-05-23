package model

type Variable struct {
	name  string
	value int64
	done  chan struct{}
}

func NewVariable(name string) *Variable {
	return &Variable{name: name, value: 0, done: make(chan struct{})}
}

func (v *Variable) SetValue(value int64) {
	v.value = value
	close(v.done)
}

func (v *Variable) GetValue() int64 {
	select {
	case <-v.done:
		return v.value
	}
}

func (v *Variable) GetName() string {
	return v.name
}

func (v *Variable) HasDependency() bool {
	return true
}
