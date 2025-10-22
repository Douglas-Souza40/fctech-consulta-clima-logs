package main

import (
    "fmt"
    "github.com/Douglas-Souza40/fctech-consulta-clima-log/processa_cep/internal/client"
    "github.com/Douglas-Souza40/fctech-consulta-clima-log/processa_cep/internal/handler"
    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
    "log"
    "net/http"
    "os"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Println("Warning: .env file not found")
    }

    r := mux.NewRouter()
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Aplicacao no ar!")
    }).Methods("GET")

    buscaTempURL := os.Getenv("BUSCA_TEMP_URL")
    if buscaTempURL == "" {
        buscaTempURL = "http://localhost:8081"
    }

    weatherHandler := &handler.WeatherHandlerProcessaCep{
        Client: client.NewClient(buscaTempURL),
    }
    r.Handle("/weather", weatherHandler).Methods("POST")

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    log.Printf("Listening on port %s", port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}
