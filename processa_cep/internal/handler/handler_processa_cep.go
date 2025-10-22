package handler

import (
    "encoding/json"
    "errors"
    "github.com/Douglas-Souza40/fctech-consulta-clima-log/processa_cep/internal/client"
    "github.com/Douglas-Souza40/fctech-consulta-clima-log/processa_cep/internal/handler/handler_error"
    "github.com/Douglas-Souza40/fctech-consulta-clima-log/processa_cep/internal/service"
    "log"
    "net/http"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
)

type WeatherHandlerProcessaCep struct {
    Client client.ClientInterface
}

type requestBody struct {
    Cep string `json:"cep"`
}

func (h *WeatherHandlerProcessaCep) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    var req requestBody
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        log.Printf("Erro ao decodificar body: %v", err)
        http.Error(w, handler_error.ErrInvalidZipcode.Error(), http.StatusUnprocessableEntity)
        return
    }

    cep := req.Cep
    if !service.IsValidCEP(cep) {
        log.Printf("Cep invalido: %v", cep)
        http.Error(w, handler_error.ErrInvalidZipcode.Error(), http.StatusUnprocessableEntity)
        return
    }

    // create a span for the external temperature lookup
    ctx := r.Context()
    tracer := otel.Tracer("processa_cep")
    ctx, span := tracer.Start(ctx, "GetTemperatureByCEP")
    // attach cep as attribute
    span.SetAttributes(attribute.String("cep", cep))
    defer span.End()

    temperatura, err := h.Client.GetTemperatureByCEP(ctx, cep)
    if err != nil {
        sc := span.SpanContext()
        log.Printf("Erro ao buscar temperatura: %v trace_id=%s span_id=%s", err, sc.TraceID().String(), sc.SpanID().String())
        if errors.Is(err, handler_error.ErrZipcodeNotFound) {
            http.Error(w, handler_error.ErrZipcodeNotFound.Error(), http.StatusNotFound)
            return
        }
        if errors.Is(err, handler_error.ErrInvalidZipcode) {
            http.Error(w, handler_error.ErrInvalidZipcode.Error(), http.StatusUnprocessableEntity)
            return
        }
        span.SetAttributes(attribute.String("error", err.Error()))
        http.Error(w, handler_error.ErrInternal.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(temperatura)
}
