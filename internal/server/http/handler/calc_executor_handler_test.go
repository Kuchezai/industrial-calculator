package handler_test

import (
	"bytes"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"industrial-calculator/internal/model"
	"industrial-calculator/internal/server/http/handler"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidateAndTransformRequest(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		requestBody    string
		expectedStatus int
		expectedError  error
	}{
		{
			name:           "invalid method",
			method:         http.MethodGet,
			requestBody:    `[]`,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedError:  errors.New("method not allowed"),
		},
		{
			name:           "invalid json",
			method:         http.MethodPost,
			requestBody:    `invalid json`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  errors.New("invalid request body"),
		},
		{
			name:           "invalid command type",
			method:         http.MethodPost,
			requestBody:    `[{"type": "invalid", "var": "x"}]`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  errors.New("invalid request body"),
		},
		{
			name:           "invalid operation",
			method:         http.MethodPost,
			requestBody:    `[{"type": "calc", "op": "invalid", "var": "x", "left": 1, "right": 2}]`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  errors.New("invalid request body"),
		},
		{
			name:           "invalid left argument type",
			method:         http.MethodPost,
			requestBody:    `[{"type": "calc", "op": "+", "var": "x", "left": true, "right": 2}]`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  errors.New("invalid request body"),
		},
		{
			name:           "invalid right argument type",
			method:         http.MethodPost,
			requestBody:    `[{"type": "calc", "op": "+", "var": "x", "left": 1, "right": false}]`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  errors.New("invalid request body"),
		},
		{
			name:        "valid print command",
			method:      http.MethodPost,
			requestBody: `[{"type": "print", "var": "x"}]`,
		},
		{
			name:        "valid calc command with numbers",
			method:      http.MethodPost,
			requestBody: `[{"type": "calc", "op": "+", "var": "x", "left": 1, "right": 2}]`,
		},
		{
			name:        "valid calc command with variables",
			method:      http.MethodPost,
			requestBody: `[{"type": "calc", "op": "*", "var": "y", "left": "x", "right": "z"}]`,
		},
		{
			name:   "mixed valid commands",
			method: http.MethodPost,
			requestBody: `[
				{"type": "calc", "op": "+", "var": "x", "left": 1, "right": 2},
				{"type": "calc", "op": "*", "var": "y", "left": "x", "right": 3},
				{"type": "print", "var": "y"}
			]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := &mockCalcExecutorUsecase{}
			h := handler.NewCalcExecutorHandler(mockUsecase)

			req := httptest.NewRequest(tt.method, "/", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			commands, err := h.ValidateAndTransformRequest(w, req)

			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedStatus, w.Code)
				assert.Equal(t, tt.expectedError.Error()+"\n", w.Body.String())
				assert.Nil(t, commands)
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, commands)

			if tt.name == "valid print command" {
				assert.Len(t, commands, 1)
				assert.Equal(t, model.Print, commands[0].Type)
				assert.Equal(t, "x", commands[0].Var.GetName())
			}

			if tt.name == "mixed valid commands" {
				assert.Len(t, commands, 3)
				assert.Equal(t, model.Calc, commands[0].Type)
				assert.Equal(t, model.Print, commands[2].Type)
			}
		})
	}
}

type mockCalcExecutorUsecase struct{}

func (m *mockCalcExecutorUsecase) ExecuteInstructions(ctx context.Context, commands []model.Command) []*model.Variable {
	return nil
}
