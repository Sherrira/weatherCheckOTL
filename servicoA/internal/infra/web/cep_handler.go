package web

import (
	"encoding/json"
	"io"
	"net/http"
	"servicoA/internal/usecase"
)

type CEP struct {
	CEP string `json:"cep"`
}

func CepHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not implemented", http.StatusNotImplemented)
		return
	}

	cep := CEP{}
	err := json.NewDecoder(r.Body).Decode(&cep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !usecase.IsValidCEP(cep.CEP) {
		http.Error(w, "Invalid ZIP code", http.StatusBadRequest)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:8081/consulta?cep="+cep.CEP, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(body)
}
