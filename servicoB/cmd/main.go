package main

import (
	"net/http"
	"servicoB/internal/infra/web"
)

func main() {
	http.HandleFunc("/consulta", web.GetTemperaturesHandler)
	http.ListenAndServe(":8081", nil)
}
