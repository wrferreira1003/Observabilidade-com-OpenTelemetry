package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/wrferreira1003/Deploy-Cloud-GO/configs"
	"github.com/wrferreira1003/Deploy-Cloud-GO/internal/api"
	"github.com/wrferreira1003/Deploy-Cloud-GO/internal/services"
	"github.com/wrferreira1003/Deploy-Cloud-GO/internal/tracing"
	"github.com/wrferreira1003/Deploy-Cloud-GO/internal/usecase"
)

func main() {
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}
	fmt.Println("Config loaded successfully", cfg)

	tracer := tracing.InitTracer("servico_b")
	defer tracer()

	weatherService := services.NewWeatherAPIService(cfg.WeatherApiKey, cfg.WeatherBaseURL)
	locationService := services.NewViaCepService(cfg)

	weatherUsecase := usecase.NewWeatherUsecase(weatherService, locationService)
	weatherHandler := api.NewWeatherHandler(weatherUsecase)

	http.HandleFunc("/weather", weatherHandler.GetWeatherHandler)

	// Obter a porta da variável de ambiente PORT
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Porta padrão para testes locais
		log.Printf("Variável de ambiente PORT não definida. Usando porta padrão %s", port)
	}

	log.Printf("Iniciando servidor na porta %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Falha ao iniciar o servidor: %v", err)
	}
}
