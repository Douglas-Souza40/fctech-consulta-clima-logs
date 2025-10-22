package handler

import (
    "encoding/json"
    "errors"
    "github.com/Douglas-Souza40/fctech-consulta-clima-log/busca_temperatura/internal/client"
    "github.com/Douglas-Souza40/fctech-consulta-clima-log/busca_temperatura/internal/handler/handler_error"
    "github.com/Douglas-Souza40/fctech-consulta-clima-log/busca_temperatura/internal/model"
    "github.com/Douglas-Souza40/fctech-consulta-clima-log/busca_temperatura/internal/service"
    "log"
    "net/http"
    "regexp"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
)

type WeatherHandlerBuscaTemp struct {
    Client client.Client
}

var cepRegex = regexp.MustCompile(`^\d{8}$`)

func (h *WeatherHandlerBuscaTemp) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    cep := r.URL.Query().Get("cep")

    if !cepRegex.MatchString(cep) {
        log.Printf("Cep invalido recebido no busca_temperatura: %v", cep)
        http.Error(w, handler_error.ErrInvalidZipcode.Error(), http.StatusUnprocessableEntity)
        return
    }

    ctx := r.Context()
    tracer := otel.Tracer("busca_temperatura")

    // span for the viacep lookup
    ctx, spanLocate := tracer.Start(ctx, "GetLocationByCEP")
    spanLocate.SetAttributes(attribute.String("cep", cep))
    location, err := h.Client.GetLocationByCEP(ctx, cep)
    sc := spanLocate.SpanContext()
    if err != nil {
        log.Printf("Erro ao buscar CEP: %v trace_id=%s span_id=%s", err, sc.TraceID().String(), sc.SpanID().String())
    }
    spanLocate.End()

    if errors.Is(err, handler_error.ErrZipcodeNotFound) {
        log.Printf("CEP nao encontrado: %v", cep)
        http.Error(w, handler_error.ErrZipcodeNotFound.Error(), http.StatusNotFound)
        return
    } else if err != nil {
        log.Printf("Erro ao buscar CEP: %v", err)
        http.Error(w, handler_error.ErrInternal.Error(), http.StatusInternalServerError)
        return
    }

    // span for the weatherapi lookup
    ctx, spanWeather := tracer.Start(ctx, "GetTemperatureByCity")
    spanWeather.SetAttributes(attribute.String("city", location.City))
    tempC, err := h.Client.GetTemperatureByCity(ctx, location.City)
    scw := spanWeather.SpanContext()
    if err != nil {
        log.Printf("Erro ao buscar temperatura: %v trace_id=%s span_id=%s", err, scw.TraceID().String(), scw.SpanID().String())
    }
    spanWeather.End()
    if err != nil {
        log.Printf("Erro ao buscar temperatura: %v", err)
        http.Error(w, handler_error.ErrInternal.Error(), http.StatusInternalServerError)
        return
    }

    resp := model.TemperatureResponse{
        City:  location.City,
        TempC: tempC,
        TempF: service.CelsiusToFahrenheit(tempC),
        TempK: service.CelsiusToKelvin(tempC),
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}
