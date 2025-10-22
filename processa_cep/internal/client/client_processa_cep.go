package client

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"

    "github.com/Douglas-Souza40/fctech-consulta-clima-log/processa_cep/internal/handler/handler_error"
    "github.com/Douglas-Souza40/fctech-consulta-clima-log/processa_cep/internal/model"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/propagation"
)

type ClientInterface interface {
    GetLocationByCEP(ctx context.Context, cep string) (*model.Location, error)
    GetTemperatureByCEP(ctx context.Context, cep string) (*model.TemperatureResponse, error)
}

type Client struct {
    HttpClient *http.Client
    APIBURL    string
}

func NewClient(apiBURL string) *Client {
    return &Client{
        HttpClient: &http.Client{},
        APIBURL:    apiBURL,
    }
}

var _ ClientInterface = (*Client)(nil)

func (c *Client) GetLocationByCEP(ctx context.Context, cep string) (*model.Location, error) {
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep), nil)
    if err != nil {
        return nil, err
    }
    // propagate trace context
    otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
    resp, err := c.HttpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var data viacepResponse
    if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
        return nil, err
    }
    if bool(data.Erro) || data.Localidade == "" {
        return nil, handler_error.ErrZipcodeNotFound
    }
    return &model.Location{City: data.Localidade}, nil
}

func (c *Client) GetTemperatureByCEP(ctx context.Context, cep string) (*model.TemperatureResponse, error) {
    url := fmt.Sprintf("%s/temperatura?cep=%s", c.APIBURL, cep)
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        return nil, handler_error.ErrInternal
    }
    // propagate trace context
    otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
    resp, err := c.HttpClient.Do(req)

    if err != nil {
        return nil, handler_error.ErrInternal
    }
    defer resp.Body.Close()

    switch resp.StatusCode {
    case http.StatusOK:
    case http.StatusNotFound:
        return nil, handler_error.ErrZipcodeNotFound
    case http.StatusUnprocessableEntity:
        return nil, handler_error.ErrInvalidZipcode
    default:
        return nil, handler_error.ErrInternal
    }

    var temp model.TemperatureResponse
    if err := json.NewDecoder(resp.Body).Decode(&temp); err != nil {
        return nil, handler_error.ErrInternal
    }
    return &temp, nil
}

type viacepResponse struct {
    Localidade string     `json:"localidade"`
    Erro       erroViaCEP `json:"erro"`
}

type erroViaCEP bool

func (e *erroViaCEP) UnmarshalJSON(data []byte) error {
    var b bool
    if err := json.Unmarshal(data, &b); err == nil {
        *e = erroViaCEP(b)
        return nil
    }
    var s string
    if err := json.Unmarshal(data, &s); err == nil {
        *e = erroViaCEP(s == "true")
        return nil
    }
    *e = false
    return nil
}
