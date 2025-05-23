package required_variables_finder

import "industrial-calculator/internal/model"

type finder struct {
}

func NewFinder() *finder {
	return &finder{}
}

func (f *finder) FindRequiredVariables(calcCommandByVariable map[*model.Variable]model.Command, targets []*model.Variable,
) map[*model.Variable]struct{} {
	required := make(map[*model.Variable]struct{})
	visited := make(map[*model.Variable]struct{})

	var dfs func(argument *model.Variable)
	dfs = func(argument *model.Variable) {
		if _, ok := visited[argument]; ok {
			return
		}

		visited[argument] = struct{}{}

		if cmd, ok := calcCommandByVariable[argument]; ok {
			required[argument] = struct{}{}

			if cmd.Left.HasDependency() {
				dfs(f.mustGetVariableByArgument(cmd.Left))
			}

			if cmd.Right.HasDependency() {
				dfs(f.mustGetVariableByArgument(cmd.Right))
			}
		}
	}

	for _, target := range targets {
		dfs(target)
	}

	return required
}

func (f *finder) mustGetVariableByArgument(argument model.Argument) *model.Variable {
	if v, ok := argument.(*model.Variable); ok {
		return v
	}

	return nil
}
