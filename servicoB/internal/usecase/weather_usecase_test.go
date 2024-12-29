package usecase

import (
	"testing"

	"github.com/wrferreira1003/Deploy-Cloud-GO/internal/models"
)

type MockLocationService struct{}

func (m *MockLocationService) GetLocationByCep(cep string) (string, error) {
	if cep == "01001000" {
		return "São Paulo", nil
	}
	if cep == "00000001" {
		return "TesteCity", nil // Simula localização válida para este CEP
	}
	return "", models.ErrZipCodeNotFound
}

type MockWeatherService struct{}

func (m *MockWeatherService) GetTemperature(city string) (models.TemperatureResponse, error) {
	if city == "São Paulo" {
		return models.TemperatureResponse{
			TempC: 28.5,
			TempF: 83.3,
			TempK: 301.65,
		}, nil
	}
	return models.TemperatureResponse{}, models.ErrWeatherNotFound
}

func TestWeatherUsecase_GetWeatherByCep(t *testing.T) {
	locationService := &MockLocationService{}
	weatherService := &MockWeatherService{}
	usecase := NewWeatherUsecase(weatherService, locationService)

	t.Run("CEP válido com sucesso", func(t *testing.T) {
		cep := "01001000"
		result, err := usecase.GetWeatherByCep(cep)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		expected := models.TemperatureResponse{
			City:  "São Paulo",
			TempC: 28.5,
			TempF: 83.3,
			TempK: 301.65,
		}
		if result != expected {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("CEP inválido", func(t *testing.T) {
		cep := "123"
		_, err := usecase.GetWeatherByCep(cep)
		if err != models.ErrInvalidZipCode {
			t.Errorf("expected error %v, got %v", models.ErrInvalidZipCode, err)
		}
	})

	t.Run("CEP não encontrado", func(t *testing.T) {
		cep := "00000000"
		_, err := usecase.GetWeatherByCep(cep)
		if err != models.ErrZipCodeNotFound {
			t.Errorf("expected error %v, got %v", models.ErrZipCodeNotFound, err)
		}
	})

	t.Run("Erro ao buscar temperatura", func(t *testing.T) {
		locationService := &MockLocationService{}
		weatherService := &MockWeatherService{}
		usecase := NewWeatherUsecase(weatherService, locationService)

		cep := "00000001" // Força uma falha no WeatherService
		_, err := usecase.GetWeatherByCep(cep)
		if err != models.ErrWeatherNotFound {
			t.Errorf("expected error %v, got %v", models.ErrWeatherNotFound, err)
		}
	})

	t.Run("Cálculo de Kelvin a partir de Celsius", func(t *testing.T) {
		// Simula o resultado esperado em Celsius
		tempC := 28.5
		expectedTempK := tempC + 273.15

		// Chama diretamente o método de cálculo ou verifica o retorno do caso de uso
		result, err := usecase.GetWeatherByCep("01001000")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// Verifica se o cálculo de Kelvin foi correto
		if result.TempK != expectedTempK {
			t.Errorf("expected Kelvin temperature: %f, got: %f", expectedTempK, result.TempK)
		}
	})

}
