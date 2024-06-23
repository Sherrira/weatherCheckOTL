package usecase

import (
	"fmt"
	"servicoB/internal/infra/datasource"
	"servicoB/internal/usecase/business_errors"
	"strconv"
)

type TemperaturesDTO struct {
	City       string  `json:"city"`
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
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

func Execute(cep string) (TemperaturesDTO, error) {
	if !IsValidCEP(cep) {
		fmt.Printf("Invalid CEP: %s\n", cep)
		return TemperaturesDTO{}, business_errors.ErrCepValidationFailed
	}

	cityRepository := datasource.NewCityRepository()
	city, err := cityRepository.FetchCityByCEP(cep)
	if err != nil {
		fmt.Printf("Error fetching city by CEP: %v\n", err)
		return TemperaturesDTO{}, business_errors.ErrCepNotFound
	}

	temperatureRepository := datasource.NewTemperatureRepository()
	local := city["localidade"].(string)
	celsius, err := temperatureRepository.FetchTemperatureByCity(local)
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
