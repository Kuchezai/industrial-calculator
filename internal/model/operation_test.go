package model_test

import (
	"github.com/stretchr/testify/assert"
	"industrial-calculator/internal/model"
	"testing"
)

func TestOperations(t *testing.T) {
	tests := []struct {
		name          string
		symbol        string
		expectedValid bool
		expectedOp    model.Operation
	}{
		{
			name:          "Plus operation",
			symbol:        "+",
			expectedValid: true,
			expectedOp:    model.Plus,
		},
		{
			name:          "Minus operation",
			symbol:        "-",
			expectedValid: true,
			expectedOp:    model.Minus,
		},
		{
			name:          "Multiply operation",
			symbol:        "*",
			expectedValid: true,
			expectedOp:    model.Multiply,
		},
		{
			name:          "Invalid operation",
			symbol:        "/",
			expectedValid: false,
		},
		{
			name:          "Empty symbol",
			symbol:        "",
			expectedValid: false,
		},
		{
			name:          "Unknown symbol",
			symbol:        "?",
			expectedValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedValid, model.IsValidOperationBySymbol(tt.symbol))

			op := model.GetOperationBySymbol(tt.symbol)
			assert.Equal(t, tt.expectedOp, op)
		})
	}
}
