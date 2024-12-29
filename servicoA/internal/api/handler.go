package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type CepRequest struct {
	Cep string `json:"cep"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

// Verificar se a string é numérica
func IsNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func Handler(w http.ResponseWriter, r *http.Request) {

	tr := otel.Tracer("servico_a")
	ctx, span := tr.Start(r.Context(), "Processar CEP")
	defer span.End()

	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var cepRequest CepRequest

	// Decodificar o corpo da requisição JSON
	err := json.NewDecoder(r.Body).Decode(&cepRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "invalid request body"})
		return
	}

	// Validar se o CEP contém exatamente 8 dígitos
	if len(cepRequest.Cep) != 8 || !IsNumeric(cepRequest.Cep) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{"message": "invalid zipcode"})
		return
	}

	// Encaminhar o CEP para o Serviço B
	resp, err := forwardToServiceB(ctx, cepRequest.Cep)
	if err != nil {
		// Caso o Serviço B não possa ser alcançado
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "error contacting service B"})
		return
	}
	defer resp.Body.Close()

	// Verificar o código de status da resposta do Serviço B
	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(resp.StatusCode) // Encaminha o status retornado pelo Serviço B
		io.Copy(w, resp.Body)          // Encaminha o corpo da resposta
		return
	}

	// Responder com o retorno do Serviço B
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.Copy(w, resp.Body)
}

// Encaminhar o CEP ao Serviço B
func forwardToServiceB(ctx context.Context, cep string) (*http.Response, error) {
	tr := otel.Tracer("servico_a")
	_, span := tr.Start(ctx, "Chamar Serviço B")
	defer span.End()

	serviceBURL := "http://servico_b:8080/weather"
	reqBody := fmt.Sprintf(`{"cep": "%s"}`, cep)

	// Criar a requisição HTTP
	req, err := http.NewRequestWithContext(ctx, "POST", serviceBURL, strings.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Propagar o contexto de tracing via cabeçalhos HTTP
	propagator := otel.GetTextMapPropagator()
	propagator.Inject(ctx, propagation.HeaderCarrier(req.Header))

	// Log para debugging
	fmt.Printf("Encaminhando para Serviço B: %s\n", reqBody)

	// Fazer a requisição ao Serviço B
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call service B: %w", err)
	}

	return resp, nil
}
