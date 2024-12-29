package main

import (
	"net/http"

	"github.com/wrferreira1003/servicoA/internal/api"
	"github.com/wrferreira1003/servicoA/internal/tracing"
)

func main() {
	tracer := tracing.InitTracer("servico_a")
	defer tracer()

	http.HandleFunc("/process-cep", api.Handler)
	http.ListenAndServe(":8081", nil)
}
