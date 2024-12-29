package usecase

import (
	"context"
	"log"

	"github.com/wrferreira1003/Deploy-Cloud-GO/internal/models"
	"github.com/wrferreira1003/Deploy-Cloud-GO/internal/services"
	"go.opentelemetry.io/otel"
)

type WeatherUsecaseIn interface {
	GetWeatherByCep(ctx context.Context, cep string) (models.TemperatureResponse, error)
}

type weatherUsecase struct {
	weatherService  services.WeatherServiceInterface
	locationService services.LocationServiceInterface
}

func NewWeatherUsecase(
	weatherService services.WeatherServiceInterface,
	locationService services.LocationServiceInterface,
) *weatherUsecase {
	return &weatherUsecase{
		weatherService:  weatherService,
		locationService: locationService,
	}
}

func (w *weatherUsecase) GetWeatherByCep(ctx context.Context, cep string) (models.TemperatureResponse, error) {
	tr := otel.Tracer("servico_b")
	_, span := tr.Start(ctx, "GetWeatherByCep")
	defer span.End()

	if len(cep) != 8 {
		return models.TemperatureResponse{}, models.ErrInvalidZipCode
	}

	//Buscar o cep
	location, err := w.locationService.GetLocationByCep(ctx, cep)
	if err != nil {
		return models.TemperatureResponse{}, models.ErrZipCodeNotFound
	}

	//Buscar o clima
	temp, err := w.weatherService.GetTemperature(ctx, location)
	if err != nil {
		return models.TemperatureResponse{}, models.ErrWeatherNotFound
	}

	log.Println(temp)

	tempK := temp.TempC + 273.15

	return models.TemperatureResponse{
		City:  temp.City,
		TempC: temp.TempC,
		TempF: temp.TempF,
		TempK: tempK,
	}, nil
}
