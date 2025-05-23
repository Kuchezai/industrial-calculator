package grpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	api "industrial-calculator/api/industrial-calculator.v1"
	"industrial-calculator/internal/model"
	"net"
)

type CalcExecutorServer struct {
	api.UnimplementedIndustrialCalculatorServer
	uc calcExecutorUsecase
}

type calcExecutorUsecase interface {
	ExecuteInstructions(ctx context.Context, commands []model.Command) []*model.Variable
}

func NewCalcExecutorServer(usecase calcExecutorUsecase) *CalcExecutorServer {
	return &CalcExecutorServer{uc: usecase}
}

func StartGRPCServer(handler *CalcExecutorServer) error {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}

	srv := grpc.NewServer()
	api.RegisterIndustrialCalculatorServer(srv, handler)
	return srv.Serve(lis)
}

func (s *CalcExecutorServer) Process(ctx context.Context, req *api.ProcessRequest) (*api.ProcessResponse, error) {
	commands := make([]model.Command, 0, len(req.Commands))
	vars := make(map[string]*model.Variable)

	for _, cmd := range req.Commands {
		if _, exists := vars[cmd.Var]; !exists {
			vars[cmd.Var] = model.NewVariable(cmd.Var)
		}

		if cmd.Type == api.CommandType_PRINT {
			commands = append(commands, model.Command{
				Type: model.CommandType(cmd.Type),
				Var:  vars[cmd.Var],
			})
			continue
		}

		left, err := parseArgument(cmd.GetLeft(), vars)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid left argument: %v", err)
		}

		right, err := parseArgument(cmd.GetRight(), vars)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid right argument: %v", err)
		}

		commands = append(commands, model.Command{
			Type:  model.CommandType(cmd.Type),
			Var:   vars[cmd.Var],
			Op:    model.Operation(cmd.Op),
			Left:  left,
			Right: right,
		})
	}

	result := s.uc.ExecuteInstructions(ctx, commands)
	return buildResponse(result), nil
}

func parseArgument(arg interface{}, vars map[string]*model.Variable) (model.Argument, error) {
	switch v := arg.(type) {
	case *api.Command_LeftInt:
		return model.NumericArgument(v.LeftInt), nil
	case *api.Command_LeftStr:
		if varRef, ok := vars[v.LeftStr]; ok {
			return varRef, nil
		}
		return nil, status.Errorf(codes.InvalidArgument, "variable not found: %s", v.LeftStr)
	default:
		return nil, status.Error(codes.InvalidArgument, "invalid argument type")
	}
}

func buildResponse(vars []*model.Variable) *api.ProcessResponse {
	results := make([]*api.VariableResult, len(vars))
	for i, v := range vars {
		results[i] = &api.VariableResult{
			Var:   v.GetName(),
			Value: v.GetValue(),
		}
	}
	return &api.ProcessResponse{Results: results}
}
