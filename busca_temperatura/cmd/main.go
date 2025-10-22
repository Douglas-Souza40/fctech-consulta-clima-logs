package main

import (
    "github.com/Douglas-Souza40/fctech-consulta-clima-log/busca_temperatura/internal/client"
    "github.com/Douglas-Souza40/fctech-consulta-clima-log/busca_temperatura/internal/handler"
    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
    "log"
    "net/http"
    "os"
)

func main() {
    _ = godotenv.Load()

    r := mux.NewRouter()
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("busca_temperatura no ar!"))
    }).Methods("GET")

    weatherHandler := &handler.WeatherHandlerBuscaTemp{
        Client: *client.NewClient(),
    }
    r.Handle("/temperatura", weatherHandler).Methods("GET")

    port := os.Getenv("PORT")
    if port == "" {
        port = "8081"
    }
    log.Printf("Listening busca_temperatura on port %s", port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}
