package client

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "net/url"
    "os"

    "github.com/Douglas-Souza40/fctech-consulta-clima-log/busca_temperatura/internal/handler/handler_error"
    "github.com/Douglas-Souza40/fctech-consulta-clima-log/busca_temperatura/internal/model"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/propagation"
)

type ClientInterface interface {
    GetTemperatureByCity(ctx context.Context, city string) (float64, error)
    GetLocationByCEP(ctx context.Context, cep string) (*model.Location, error)
}

type Client struct {
    HttpClient *http.Client
}

func NewClient() *Client {
    return &Client{
        HttpClient: &http.Client{},
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

func (c *Client) GetTemperatureByCity(ctx context.Context, city string) (float64, error) {
    apiKey := os.Getenv("WEATHER_API_KEY")
    if apiKey == "" {
        return 0, errors.New("weather api key not set")
    }
    cityEscaped := url.QueryEscape(city)
    url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, cityEscaped)
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        return 0, err
    }
    // propagate trace context
    otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
    resp, err := c.HttpClient.Do(req)
    if err != nil {
        return 0, err
    }
    defer resp.Body.Close()

    var data struct {
        Current struct {
            TempC float64 `json:"temp_c"`
        } `json:"current"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
        return 0, err
    }
    return data.Current.TempC, nil
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
