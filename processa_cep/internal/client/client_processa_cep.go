package client

import (
    "encoding/json"
    "fmt"
    "github.com/Douglas-Souza40/fctech-consulta-clima-log/processa_cep/internal/handler/handler_error"
    "github.com/Douglas-Souza40/fctech-consulta-clima-log/processa_cep/internal/model"
    "net/http"
)

type ClientInterface interface {
    GetLocationByCEP(string) (*model.Location, error)
    GetTemperatureByCEP(string) (*model.TemperatureResponse, error)
}

type Client struct {
    HTTPGet func(string) (*http.Response, error)
    APIBURL string
}

func NewClient(apiBURL string) *Client {
    return &Client{
        HTTPGet: http.Get,
        APIBURL: apiBURL,
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

func (c *Client) GetTemperatureByCEP(cep string) (*model.TemperatureResponse, error) {
    url := fmt.Sprintf("%s/temperatura?cep=%s", c.APIBURL, cep)
    resp, err := c.HTTPGet(url)

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
