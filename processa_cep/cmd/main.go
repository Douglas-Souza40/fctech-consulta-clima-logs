package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/Douglas-Souza40/fctech-consulta-clima-log/processa_cep/internal/client"
    "github.com/Douglas-Souza40/fctech-consulta-clima-log/processa_cep/internal/handler"
    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
    "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/zipkin"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Println("Warning: .env file not found")
    }

    // init otel zipkin exporter (optional ZIPKIN_URL env)
    zipkinURL := os.Getenv("ZIPKIN_URL")
    if zipkinURL == "" {
        zipkinURL = "http://localhost:9411/api/v2/spans"
    }

    exporter, err := zipkin.New(zipkinURL)
    if err != nil {
        log.Printf("warning: could not initialize zipkin exporter: %v", err)
    } else {
        tp := sdktrace.NewTracerProvider(
            sdktrace.WithBatcher(exporter),
        )
        otel.SetTracerProvider(tp)
        defer func() {
            _ = tp.Shutdown(context.Background())
        }()
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
    r.Handle("/weather", otelhttp.NewHandler(weatherHandler, "ProcessarCEP")).Methods("POST")

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    log.Printf("Listening on port %s", port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}
