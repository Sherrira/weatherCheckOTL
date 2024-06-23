package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"servicoB/internal/usecase"
	"servicoB/internal/usecase/business_errors"
)

func GetTemperaturesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not implemented", http.StatusNotImplemented)
		return
	}
	cep := r.URL.Query().Get("cep")
	if cep == "" {
		msg := fmt.Sprintf("Query CEP empty: %s", cep)
		fmt.Println(msg)
		err := errors.New(msg)
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	dto, err := usecase.Execute(cep)
	if err != nil {
		if errors.Is(err, business_errors.ErrCepValidationFailed) {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		if errors.Is(err, business_errors.ErrCepNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if errors.Is(err, business_errors.ErrFetchTemperatureFailed) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
	}

	err = json.NewEncoder(w).Encode(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
