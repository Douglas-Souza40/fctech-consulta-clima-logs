package handler

import (
    "encoding/json"
    "errors"
    "github.com/Douglas-Souza40/fctech-consulta-clima-log/processa_cep/internal/client"
    "github.com/Douglas-Souza40/fctech-consulta-clima-log/processa_cep/internal/handler/handler_error"
    "github.com/Douglas-Souza40/fctech-consulta-clima-log/processa_cep/internal/service"
    "log"
    "net/http"
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

    temperatura, err := h.Client.GetTemperatureByCEP(cep)
    if err != nil {
        log.Printf("Erro ao buscar temperatura: %v", err)
        if errors.Is(err, handler_error.ErrZipcodeNotFound) {
            http.Error(w, handler_error.ErrZipcodeNotFound.Error(), http.StatusNotFound)
            return
        }
        if errors.Is(err, handler_error.ErrInvalidZipcode) {
            http.Error(w, handler_error.ErrInvalidZipcode.Error(), http.StatusUnprocessableEntity)
            return
        }
        http.Error(w, handler_error.ErrInternal.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(temperatura)
}
