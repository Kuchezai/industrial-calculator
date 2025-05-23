package required_variables_finder_test

import (
	"github.com/stretchr/testify/assert"
	"industrial-calculator/internal/model"
	"industrial-calculator/internal/required_variables_finder"
	"testing"
)

func TestFindRequiredVariables(t *testing.T) {
	f := required_variables_finder.NewFinder()

	vars := map[string]*model.Variable{
		"x": model.NewVariable("x"),
		"a": model.NewVariable("a"),
		"b": model.NewVariable("b"),
		"c": model.NewVariable("c"),
		"y": model.NewVariable("y"),
		"z": model.NewVariable("z"),
	}

	testCases := []struct {
		name     string
		commands map[*model.Variable]model.Command
		targets  []*model.Variable
		expected []string
	}{
		{
			name:     "empty input",
			commands: make(map[*model.Variable]model.Command),
			targets:  []*model.Variable{},
			expected: []string{},
		},
		{
			name: "single independent variable",
			commands: map[*model.Variable]model.Command{
				vars["x"]: {
					Type:  model.Calc,
					Var:   vars["x"],
					Op:    model.Plus,
					Left:  model.NumericArgument(1),
					Right: model.NumericArgument(2),
				},
			},
			targets:  []*model.Variable{vars["x"]},
			expected: []string{"x"},
		},
		{
			name: "dependency chain",
			commands: map[*model.Variable]model.Command{
				vars["a"]: {
					Type:  model.Calc,
					Var:   vars["a"],
					Op:    model.Plus,
					Left:  model.NumericArgument(1),
					Right: model.NumericArgument(2),
				},
				vars["b"]: {
					Type:  model.Calc,
					Var:   vars["b"],
					Op:    model.Multiply,
					Left:  vars["a"],
					Right: model.NumericArgument(3),
				},
				vars["c"]: {
					Type:  model.Calc,
					Var:   vars["c"],
					Op:    model.Minus,
					Left:  vars["b"],
					Right: model.NumericArgument(4),
				},
			},
			targets:  []*model.Variable{vars["c"]},
			expected: []string{"a", "b", "c"},
		},
		{
			name: "multiple targets with shared dependencies",
			commands: map[*model.Variable]model.Command{
				vars["x"]: {
					Type:  model.Calc,
					Var:   vars["x"],
					Op:    model.Plus,
					Left:  model.NumericArgument(1),
					Right: model.NumericArgument(2),
				},
				vars["y"]: {
					Type:  model.Calc,
					Var:   vars["y"],
					Op:    model.Multiply,
					Left:  vars["x"],
					Right: model.NumericArgument(3),
				},
				vars["z"]: {
					Type:  model.Calc,
					Var:   vars["z"],
					Op:    model.Minus,
					Left:  vars["x"],
					Right: model.NumericArgument(4),
				},
			},
			targets:  []*model.Variable{vars["y"], vars["z"]},
			expected: []string{"x", "y", "z"},
		},
		{
			name: "circular dependencies",
			commands: map[*model.Variable]model.Command{
				vars["a"]: {
					Type:  model.Calc,
					Var:   vars["a"],
					Op:    model.Plus,
					Left:  vars["b"],
					Right: model.NumericArgument(1),
				},
				vars["b"]: {
					Type:  model.Calc,
					Var:   vars["b"],
					Op:    model.Minus,
					Left:  vars["a"],
					Right: model.NumericArgument(2),
				},
			},
			targets:  []*model.Variable{vars["a"]},
			expected: []string{"a", "b"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := f.FindRequiredVariables(tc.commands, tc.targets)

			var resultNames []string
			for v := range result {
				resultNames = append(resultNames, v.GetName())
			}

			assert.Equal(t, len(tc.expected), len(resultNames))

			for _, name := range tc.expected {
				found := false
				for _, resName := range resultNames {
					if resName == name {
						found = true
						break
					}
				}
				assert.True(t, found, "variable %s not found in result", name)
			}
		})
	}
}
