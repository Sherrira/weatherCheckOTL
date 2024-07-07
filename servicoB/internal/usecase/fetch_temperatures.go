package usecase

import (
	"context"
	"fmt"
	"servicoB/configs"
	"servicoB/internal/infra/datasource"
	"servicoB/internal/usecase/business_errors"
	"strconv"

	"go.opentelemetry.io/otel/trace"
)

type TemperaturesUseCase struct {
	Ctx                   context.Context
	CityRepository        *datasource.CityRepository
	TemperatureRepository *datasource.TemperatureRepository
	Configs               *configs.Config
	otelTracer            trace.Tracer
}

type TemperaturesDTO struct {
	City       string  `json:"city"`
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
}

func NewTemperaturesUseCase(ctx context.Context, conf *configs.Config, tracer trace.Tracer) *TemperaturesUseCase {
	cityRepository := datasource.NewCityRepository(conf)
	temperatureRepository := datasource.NewTemperatureRepository(conf)
	return &TemperaturesUseCase{
		Ctx:                   ctx,
		CityRepository:        cityRepository,
		TemperatureRepository: temperatureRepository,
		Configs:               conf,
		otelTracer:            tracer,
	}
}

func IsValidCEP(cep string) bool {
	if len(cep) != 8 {
		return false
	}
	_, err := strconv.Atoi(cep)
	return err == nil
}

func ConvertCelsiusToFahrenheit(celsius float64) float64 {
	return celsius*1.8 + 32
}

func ConvertCelsiusToKelvin(celsius float64) float64 {
	return celsius + 273
}

func (u *TemperaturesUseCase) Execute(cep string) (TemperaturesDTO, error) {
	if !IsValidCEP(cep) {
		fmt.Printf("Invalid CEP: %s\n", cep)
		return TemperaturesDTO{}, business_errors.ErrCepValidationFailed
	}

	ctx, fetchCitySpan := u.otelTracer.Start(u.Ctx, "External Call - FetchCityByCEP")
	city, err := u.CityRepository.FetchCityByCEP(cep)
	fetchCitySpan.End()
	if err != nil {
		fmt.Printf("Error fetching city by CEP: %v\n", err)
		return TemperaturesDTO{}, business_errors.ErrCepNotFound
	}

	local := city["localidade"].(string)
	_, fetchTempSpan := u.otelTracer.Start(ctx, "External Call - FetchTemperatureByCity")
	celsius, err := u.TemperatureRepository.FetchTemperatureByCity(local)
	fetchTempSpan.End()
	if err != nil {
		fmt.Printf("Error fetching temperature by city: %v\n", err)
		return TemperaturesDTO{}, business_errors.ErrFetchTemperatureFailed
	}
	fahrenheit := ConvertCelsiusToFahrenheit(celsius)
	kelvin := ConvertCelsiusToKelvin(celsius)

	result := TemperaturesDTO{
		City:       local,
		Celsius:    celsius,
		Fahrenheit: fahrenheit,
		Kelvin:     kelvin,
	}
	fmt.Printf("Success getting temperatures: %+v\n", result)

	return result, nil
}
