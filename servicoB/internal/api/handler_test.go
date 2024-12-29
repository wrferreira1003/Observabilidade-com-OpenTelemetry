package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/wrferreira1003/Deploy-Cloud-GO/internal/models"
)

// Mock do caso de uso
type MockWeatherUseCase struct{}

func (m *MockWeatherUseCase) GetWeatherByCep(cep string) (models.TemperatureResponse, error) {
	switch cep {
	case "01001000":
		return models.TemperatureResponse{
			TempC: 28.5,
			TempF: 83.3,
			TempK: 301.65,
		}, nil
	case "00000000":
		return models.TemperatureResponse{}, models.ErrZipCodeNotFound
	default:
		return models.TemperatureResponse{}, models.ErrInvalidZipCode
	}
}

func TestGetWeatherHandler(t *testing.T) {
	mockUseCase := &MockWeatherUseCase{}
	handler := NewWeatherHandler(mockUseCase)

	// Teste de sucesso
	t.Run("Sucesso", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/weather?cep=01001000", nil)
		rr := httptest.NewRecorder()

		handler.GetWeatherHandler(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		expected := `{"temp_C":28.5,"temp_F":83.3,"temp_K":301.65}`
		if strings.TrimSpace(rr.Body.String()) != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})

	// Teste de CEP inválido
	t.Run("CEP Inválido", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/weather?cep=123", nil)
		rr := httptest.NewRecorder()

		handler.GetWeatherHandler(rr, req)

		if status := rr.Code; status != http.StatusUnprocessableEntity {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnprocessableEntity)
		}

		expected := `{"message": "invalid zipcode"}`
		if strings.TrimSpace(rr.Body.String()) != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})

	t.Run("CEP Não Encontrado", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/weather?cep=00000000", nil)
		rr := httptest.NewRecorder()

		handler.GetWeatherHandler(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}

		expected := `{"message": "can not find zipcode"}`
		if strings.TrimSpace(rr.Body.String()) != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})

	//Teste do retorno dos dados
	t.Run("Retorno dos dados", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/weather?cep=01001000", nil)
		rr := httptest.NewRecorder()

		handler.GetWeatherHandler(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		expected := `{"temp_C":28.5,"temp_F":83.3,"temp_K":301.65}`
		if strings.TrimSpace(rr.Body.String()) != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})
}
