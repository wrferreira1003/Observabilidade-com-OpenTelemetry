package services

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWeatherAPIService_GetTemperature(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if query.Get("key") != "mock-api-key" {
			http.Error(w, "missing or incorrect API key", http.StatusUnauthorized)
			return
		}
		if query.Get("q") != "São Paulo" {
			http.Error(w, "city not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"current": {"temp_c": 28.5, "temp_f": 83.3}, "location": {"name": "São Paulo"}}`)
	}))
	defer mockServer.Close()

	// Configurar o serviço para usar o mockServer
	mockService := &WeatherAPIService{
		APIKey:  "mock-api-key",
		BaseURL: mockServer.URL, // Redirecionar chamadas para o mockServer
	}

	t.Run("Sucesso", func(t *testing.T) {
		temp, err := mockService.GetTemperature(context.Background(), "São Paulo")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if temp.TempC != 28.5 || temp.TempF != 83.3 {
			t.Errorf("expected temps: 28.5°C, 83.3°F; got: %v°C, %v°F", temp.TempC, temp.TempF)
		}

		if temp.City != "São Paulo" {
			t.Errorf("expected city: São Paulo; got: %v", temp.City)
		}
	})

	t.Run("Erro - Cidade não encontrada", func(t *testing.T) {
		_, err := mockService.GetTemperature(context.Background(), "Cidade Inexistente")
		if err == nil || err.Error() != "weatherapi returned status: 404" {
			t.Errorf("expected error 'weatherapi returned status: 404', got %v", err)
		}
	})

	t.Run("Erro - Chave de API inválida", func(t *testing.T) {
		mockService.APIKey = "invalid-api-key"

		_, err := mockService.GetTemperature(context.Background(), "São Paulo")
		if err == nil || err.Error() != "weatherapi returned status: 401" {
			t.Errorf("expected error 'weatherapi returned status: 401', got %v", err)
		}
	})
}
