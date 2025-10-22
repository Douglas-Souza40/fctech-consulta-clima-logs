package model

type Location struct {
    City string
}

type TemperatureResponse struct {
    TempC float64 `json:"temp_C"`
    TempF float64 `json:"temp_F"`
    TempK float64 `json:"temp_K"`
}
