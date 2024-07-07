package datasource

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"servicoB/configs"
)

type TemperatureRepository struct {
	HTTPClient HTTPClient
	conf       configs.Config
}

func NewTemperatureRepository(conf *configs.Config) *TemperatureRepository {
	return &TemperatureRepository{
		HTTPClient: &http.Client{},
		conf:       *conf,
	}
}

func NewTemperatureRepositoryForTest(client HTTPClient, conf *configs.Config) *TemperatureRepository {
	return &TemperatureRepository{
		HTTPClient: client,
		conf:       *conf,
	}
}

type TemperatureAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

func (t *TemperatureRepository) FetchTemperatureByCity(local string) (float64, error) {
	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s", t.conf.WEATHER_API_KEY, url.QueryEscape(local))
	fmt.Printf("Calling GET weatherapi: %v\n", url)
	resp, err := t.HTTPClient.Get(url)
	if err != nil {
		fmt.Printf("Error fetching weather API: %v\n", err)
		return 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading weather API response: %v\n", err)
		return 0, err
	}

	fmt.Printf("Success fetching temperature: %s\n", body)
	var weatherAPIResponse TemperatureAPIResponse
	err = json.Unmarshal(body, &weatherAPIResponse)
	if err != nil {
		fmt.Printf("Error unmarshalling weather API response: %v\n", err)
		return 0, err
	}
	return weatherAPIResponse.Current.TempC, nil
}
