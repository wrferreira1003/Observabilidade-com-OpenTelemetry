package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wrferreira1003/Deploy-Cloud-GO/internal/models"
	"github.com/wrferreira1003/Deploy-Cloud-GO/internal/usecase"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type WeatherHandler struct {
	weatherUsecase usecase.WeatherUsecaseIn
}

func NewWeatherHandler(weatherUsecase usecase.WeatherUsecaseIn) *WeatherHandler {
	return &WeatherHandler{
		weatherUsecase: weatherUsecase,
	}
}

func (h *WeatherHandler) GetWeatherHandler(w http.ResponseWriter, r *http.Request) {
	// Extrair o contexto do tracing dos cabeçalhos
	propagator := otel.GetTextMapPropagator()
	ctx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))

	tr := otel.Tracer("servico_b")
	ctx, span := tr.Start(ctx, "GetWeatherHandler")
	defer span.End()

	// Validar se o método é POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Estrutura para decodificar o corpo da requisição
	var req struct {
		Cep string `json:"cep"`
	}

	// Decodificar o corpo da requisição
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "invalid request body"})
		return
	}

	// Validar se o CEP contém exatamente 8 dígitos e é numérico
	if len(req.Cep) != 8 || !isNumeric(req.Cep) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{"message": "invalid zipcode"})
		return
	}

	// Buscar o clima
	weather, err := h.weatherUsecase.GetWeatherByCep(ctx, req.Cep)
	if err != nil {
		fmt.Println("Error:", err)
		h.handleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(weather)
}

func (h *WeatherHandler) handleError(w http.ResponseWriter, err error) {
	var statusCode int
	var message string

	switch err {
	case models.ErrInvalidZipCode:
		statusCode = http.StatusUnprocessableEntity
		message = `{"message": "invalid zipcode"}`
	case models.ErrZipCodeNotFound:
		statusCode = http.StatusNotFound
		message = `{"message": "can not find zipcode"}`
	case models.ErrWeatherNotFound:
		statusCode = http.StatusNotFound
		message = `{"message": "can not find weather"}`
	default:
		statusCode = http.StatusInternalServerError
		message = `{"message": "internal server error"}`
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
