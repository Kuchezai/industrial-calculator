package sentence

import (
	"context"
	"industrial-calculator/internal/model"
)

type Sentence struct {
	vr    *model.Variable
	op    model.Operation
	left  model.Argument
	right model.Argument
}

func NewSentence(vr *model.Variable, op model.Operation, left model.Argument, right model.Argument) *Sentence {
	return &Sentence{vr: vr, op: op, left: left, right: right}
}

func (s *Sentence) Calc(ctx context.Context) {
	done := make(chan struct{})
	defer close(done)
	go func() {
		s.vr.SetValue(mustCalcTwoValuesByOperation(s.left.GetValue(), s.right.GetValue(), s.op))

		done <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return
	case <-done:
		return
	}
}

func mustCalcTwoValuesByOperation(a, b int64, op model.Operation) int64 {
	switch op {
	case model.Plus:
		return a + b
	case model.Minus:
		return a - b
	case model.Multiply:
		return a * b
	}

	return 0
}
