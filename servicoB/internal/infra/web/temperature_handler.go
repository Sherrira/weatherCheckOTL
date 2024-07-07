package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"servicoB/configs"
	"servicoB/internal/usecase"
	"servicoB/internal/usecase/business_errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type TemperatureHandler struct {
	Configs    *configs.Config
	HttpClient *http.Client
	OTELTracer trace.Tracer
}

func NewTemperatureHandler(configs *configs.Config, httpClient *http.Client, otelTracer trace.Tracer) *TemperatureHandler {
	return &TemperatureHandler{
		Configs:    configs,
		HttpClient: httpClient,
		OTELTracer: otelTracer,
	}
}

func (h *TemperatureHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	ctx, span := h.OTELTracer.Start(ctx, "Serviço B - Integrações externas de temperatura")
	defer span.End()
	if r.Method != "GET" {
		errMsg := "Method not implemented"
		log.Println(errMsg)
		http.Error(w, errMsg, http.StatusNotImplemented)
		return
	}
	cep := r.URL.Query().Get("cep")
	if cep == "" {
		msg := fmt.Sprintf("Query CEP empty: %s", cep)
		fmt.Println(msg)
		err := errors.New(msg)
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	tempUseCase := usecase.NewTemperaturesUseCase(ctx, h.Configs, h.OTELTracer)
	dto, err := tempUseCase.Execute(cep)
	if err != nil {
		if errors.Is(err, business_errors.ErrCepValidationFailed) {
			log.Printf("Error validating CEP: %v\n", err)
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		if errors.Is(err, business_errors.ErrCepNotFound) {
			log.Printf("CEP not found: %v\n", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if errors.Is(err, business_errors.ErrFetchTemperatureFailed) {
			log.Printf("Error fetching temperature: %v\n", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
	}

	err = json.NewEncoder(w).Encode(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Temperature fetched: %v\n", dto)
}
