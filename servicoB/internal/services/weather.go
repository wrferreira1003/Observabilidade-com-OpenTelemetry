package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/wrferreira1003/Deploy-Cloud-GO/internal/models"
	"go.opentelemetry.io/otel"
)

type WeatherServiceInterface interface {
	GetTemperature(ctx context.Context, city string) (models.TemperatureResponse, error)
}

type WeatherAPIService struct {
	APIKey  string
	BaseURL string // Adicionando URL base configurável
}

func NewWeatherAPIService(apiKey string, baseURL string) *WeatherAPIService {
	return &WeatherAPIService{
		APIKey:  apiKey,
		BaseURL: baseURL,
	}
}

func (s *WeatherAPIService) GetTemperature(ctx context.Context, city string) (models.TemperatureResponse, error) {
	tr := otel.Tracer("servico_b")
	_, span := tr.Start(ctx, "GetTemperature")
	defer span.End()

	// Codificar o nome da cidade
	encodedCity := url.QueryEscape(city)

	url := fmt.Sprintf("%s?key=%s&q=%s", s.BaseURL, s.APIKey, encodedCity)

	resp, err := http.Get(url)
	if err != nil {
		return models.TemperatureResponse{}, fmt.Errorf("failed to fetch weather data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.TemperatureResponse{}, fmt.Errorf("weatherapi returned status: %d", resp.StatusCode)
	}

	var data struct {
		Current struct {
			TempC float64 `json:"temp_c"`
			TempF float64 `json:"temp_f"`
		} `json:"current"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return models.TemperatureResponse{}, fmt.Errorf("failed to parse weather data: %w", err)
	}

	log.Println(data)

	return models.TemperatureResponse{
		TempC: data.Current.TempC,
		TempF: data.Current.TempF,
		TempK: data.Current.TempC + 273.15,
		City:  city,
	}, nil
}
