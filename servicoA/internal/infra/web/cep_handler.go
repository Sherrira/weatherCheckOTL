package web

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"servicoA/configs"
	"servicoA/internal/usecase"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type CEP struct {
	CEP string `json:"cep"`
}

type CepHandler struct {
	Config     *configs.Config
	HttpClient *http.Client
	OTELTracer trace.Tracer
}

func NewCepHandler(config *configs.Config, httpClient *http.Client, otelTracer trace.Tracer) *CepHandler {
	return &CepHandler{
		Config:     config,
		HttpClient: httpClient,
		OTELTracer: otelTracer,
	}
}

func (h *CepHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	ctx, span := h.OTELTracer.Start(ctx, "Serviço A - Validação de CEP")
	defer span.End()

	if r.Method != "POST" {
		errMsg := "Method not implemented"
		log.Println(errMsg)
		http.Error(w, errMsg, http.StatusNotImplemented)
		return
	}

	cep := CEP{}
	err := json.NewDecoder(r.Body).Decode(&cep)
	if err != nil {
		log.Printf("Error decoding request body: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !usecase.IsValidCEP(cep.CEP) {
		errMsg := "Invalid ZIP code"
		log.Printf("%s: %s\n", errMsg, cep.CEP)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	url := h.Config.BASE_URL_SERVICE_B + "/consulta?cep=" + cep.CEP
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Printf("Error creating request: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	resp, err := h.HttpClient.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Response from service B: %s\n", body)
	w.Write(body)
}
