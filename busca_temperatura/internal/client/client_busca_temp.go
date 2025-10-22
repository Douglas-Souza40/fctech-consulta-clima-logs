package client

import (
    "encoding/json"
    "errors"
    "fmt"
    "github.com/Douglas-Souza40/fctech-consulta-clima-log/busca_temperatura/internal/handler/handler_error"
    "github.com/Douglas-Souza40/fctech-consulta-clima-log/busca_temperatura/internal/model"
    "net/http"
    "net/url"
    "os"
)

type ClientInterface interface {
    GetTemperatureByCity(string) (float64, error)
    GetLocationByCEP(string) (*model.Location, error)
}

type Client struct {
    HTTPGet func(string) (*http.Response, error)
}

func NewClient() *Client {
    return &Client{
        HTTPGet: http.Get,
    }
}

var _ ClientInterface = (*Client)(nil)

func (c *Client) GetLocationByCEP(cep string) (*model.Location, error) {
    resp, err := c.HTTPGet(fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep))
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

func (c *Client) GetTemperatureByCity(city string) (float64, error) {
    apiKey := os.Getenv("WEATHER_API_KEY")
    if apiKey == "" {
        return 0, errors.New("weather api key not set")
    }
    cityEscaped := url.QueryEscape(city)
    url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, cityEscaped)
    resp, err := c.HTTPGet(url)
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
