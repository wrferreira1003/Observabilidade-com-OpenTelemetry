package services

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wrferreira1003/Deploy-Cloud-GO/configs"
)

func TestViaCepService_GetLocationByCep(t *testing.T) {
	// Servidor simulado para mockar respostas
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/01001000/json/" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"localidade": "São Paulo"}`)
		} else if r.URL.Path == "/00000000/json/" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"localidade": ""}`)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockServer.Close()

	// Configuração com URL do servidor simulado
	cfg := &configs.Config{ViaCepUrl: mockServer.URL}
	service := NewViaCepService(cfg)

	t.Run("Sucesso", func(t *testing.T) {
		location, err := service.GetLocationByCep(context.Background(), "01001000")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if location != "São Paulo" {
			t.Errorf("expected 'São Paulo', got %v", location)
		}
	})

	t.Run("Localidade Não Encontrada", func(t *testing.T) {
		_, err := service.GetLocationByCep(context.Background(), "00000000")
		if err == nil || err.Error() != "location not found for cep: 00000000" {
			t.Errorf("expected error 'location not found for cep: 00000000', got %v", err)
		}
	})

	t.Run("CEP Inválido", func(t *testing.T) {
		_, err := service.GetLocationByCep(context.Background(), "123")
		if err == nil || err.Error() != "error fetching data from ViaCep" {
			t.Errorf("expected error 'error fetching data from ViaCep', got %v", err)
		}
	})
}
