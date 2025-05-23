package handler

import (
	"context"
	"encoding/json"
	"errors"
	"industrial-calculator/internal/model"
	"net/http"
	"time"
)

type CalcExecutorHandler struct {
	uc calcExecutorUsecase
}

type calcExecutorUsecase interface {
	ExecuteInstructions(ctx context.Context, commands []model.Command) []*model.Variable
}

func NewCalcExecutorHandler(usecase calcExecutorUsecase) *CalcExecutorHandler {
	return &CalcExecutorHandler{uc: usecase}
}

type Request []struct {
	Type  string      `json:"type"`
	Op    string      `json:"op,omitempty"`
	Var   string      `json:"var"`
	Left  interface{} `json:"left,omitempty"`
	Right interface{} `json:"right,omitempty"`
}

type Response struct {
	Items []Item `json:"items"`
}

type Item struct {
	Var   string `json:"var"`
	Value int64  `json:"value"`
}

func (h *CalcExecutorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	commands, err := h.ValidateAndTransformRequest(w, r)
	if err != nil {
		return
	}

	result := h.uc.ExecuteInstructions(ctx, commands)

	h.writeResponse(result, w)

	return
}

func (h *CalcExecutorHandler) ValidateAndTransformRequest(w http.ResponseWriter, r *http.Request) ([]model.Command, error) {
	if r.Method != http.MethodPost {
		err := errors.New("method not allowed")
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)

		return nil, err
	}

	var req Request
	errRequestBody := errors.New("invalid request body")

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, errRequestBody.Error(), http.StatusBadRequest)
		return nil, errRequestBody
	}

	commands := make([]model.Command, len(req))
	vars := make(map[string]*model.Variable)

	for _, cmd := range req {
		if _, ok := vars[cmd.Var]; !ok {
			vars[cmd.Var] = model.NewVariable(cmd.Var)
		}
	}

	for i, cmd := range req {
		if !model.IsValidCommand(model.CommandType(cmd.Type)) {
			http.Error(w, errRequestBody.Error(), http.StatusBadRequest)
			return nil, errRequestBody
		}

		if model.CommandType(cmd.Type) == model.Print {
			command := model.Command{
				Type: model.CommandType(cmd.Type),
				Var:  vars[cmd.Var],
			}
			commands[i] = command

			continue
		}

		if !model.IsValidOperationBySymbol(cmd.Op) {
			http.Error(w, errRequestBody.Error(), http.StatusBadRequest)
			return nil, errRequestBody
		}

		var left model.Argument
		switch v := cmd.Left.(type) {
		case string:
			left = vars[v]
		case float64:
			left = model.NumericArgument(v)
		default:
			http.Error(w, errRequestBody.Error(), http.StatusBadRequest)
			return nil, errRequestBody
		}

		var right model.Argument
		switch v := cmd.Right.(type) {
		case string:
			right = vars[v]
		case float64:
			right = model.NumericArgument(v)
		default:
			http.Error(w, errRequestBody.Error(), http.StatusBadRequest)
			return nil, errRequestBody
		}

		command := model.Command{
			Type:  model.CommandType(cmd.Type),
			Var:   vars[cmd.Var],
			Op:    model.GetOperationBySymbol(cmd.Op),
			Left:  left,
			Right: right,
		}

		commands[i] = command
	}

	return commands, nil
}

func (h *CalcExecutorHandler) writeResponse(result []*model.Variable, w http.ResponseWriter) {
	items := make([]Item, len(result))
	for i := range items {
		items[i] = Item{
			Var:   result[i].GetName(),
			Value: result[i].GetValue(),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(Response{Items: items}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}

	return
}
