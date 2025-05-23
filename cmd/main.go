package main

import (
	"industrial-calculator/internal/required_variables_finder"
	"industrial-calculator/internal/server/grpc"
	"industrial-calculator/internal/server/http"
	"industrial-calculator/internal/server/http/handler"
	"industrial-calculator/internal/usecase"
	"log"
	"sync"
)

func main() {
	finder := required_variables_finder.NewFinder()
	uc := usecase.NewCalcExectureUsecase(finder)
	restHandler := handler.NewCalcExecutorHandler(uc)
	grpcHandler := grpc.NewCalcExecutorServer(uc)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		log.Println("Starting HTTP server on :8080")
		http.StartHTTPServer(restHandler)
	}()

	go func() {
		defer wg.Done()
		log.Println("Starting gRPC server on :50051")
		if err := grpc.StartGRPCServer(grpcHandler); err != nil {
			log.Fatalf("gRPC server failed: %v", err)
		}
	}()

	wg.Wait()
}
