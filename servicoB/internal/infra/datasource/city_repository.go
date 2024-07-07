package datasource

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"servicoB/configs"
)

type CityRepository struct {
	HTTPClient HTTPClient
	Configs    *configs.Config
}

func NewCityRepository(conf *configs.Config) *CityRepository {
	// TODO: Remove InsecureSkipVerify
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transport}
	return &CityRepository{
		HTTPClient: client,
		Configs:    conf,
	}
}

func NewCityRepositoryForTest(client HTTPClient) *CityRepository {
	return &CityRepository{HTTPClient: client}
}

func (c *CityRepository) FetchCityByCEP(cep string) (map[string]interface{}, error) {
	url := "https://viacep.com.br/ws/" + cep + "/json/"
	fmt.Printf("Calling GET viacep: %v\n", url)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		fmt.Printf("Error fetching city by CEP: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading city by CEP response: %v\n", err)
		return nil, err
	}

	fmt.Printf("Success fetching city: %s\n", body)
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Printf("Error unmarshalling city by CEP response: %v\n", err)
		return nil, err
	}
	if result["erro"] != nil {
		fmt.Printf("CEP not found: %s\n", cep)
		return nil, errors.New("CEP not found")
	}
	return result, nil
}
