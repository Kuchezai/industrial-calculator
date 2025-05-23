package usecase

import (
	"context"
	"industrial-calculator/internal/model"
	"industrial-calculator/internal/sentence"
	"sync"
)

type requiredVariablesFinder interface {
	FindRequiredVariables(calcCommandByVariable map[*model.Variable]model.Command, targets []*model.Variable) map[*model.Variable]struct{}
}

type CalcExecutorUsecase struct {
	finder requiredVariablesFinder
}

func NewCalcExectureUsecase(finder requiredVariablesFinder) *CalcExecutorUsecase {
	return &CalcExecutorUsecase{finder: finder}
}

func (c *CalcExecutorUsecase) ExecuteInstructions(ctx context.Context, commands []model.Command) []*model.Variable {
	calcCommandsByVariable := make(map[*model.Variable]model.Command)
	printTargets := make([]*model.Variable, 0)

	for i := range commands {
		if commands[i].IsPrint() {
			printTargets = append(printTargets, commands[i].Var)
		} else if commands[i].IsCalc() {
			calcCommandsByVariable[commands[i].Var] = commands[i]
		}
	}

	requiredVariables := c.finder.FindRequiredVariables(calcCommandsByVariable, printTargets)

	var wg sync.WaitGroup
	for variable := range requiredVariables {
		cmd := calcCommandsByVariable[variable]

		sentence := sentence.NewSentence(variable, cmd.Op, cmd.Left, cmd.Right)

		wg.Add(1)
		go func() {
			sentence.Calc(ctx)
			wg.Done()
		}()
	}

	wg.Wait()

	return printTargets
}
