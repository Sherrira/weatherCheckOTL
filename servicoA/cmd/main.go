package main

import (
	"net/http"
	"servicoA/internal/infra/web"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cep", web.CepHandler)
	http.ListenAndServe(":8080", mux)
}
