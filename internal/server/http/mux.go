package http

import "net/http"

func StartHTTPServer(handler http.Handler) {
	mux := http.NewServeMux()
	mux.Handle("/process", handler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
