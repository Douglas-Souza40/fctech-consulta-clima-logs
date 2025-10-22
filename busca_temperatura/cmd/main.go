package main

import (
    "context"
    "log"
    "net/http"
    "os"

    "github.com/Douglas-Souza40/fctech-consulta-clima-log/busca_temperatura/internal/client"
    "github.com/Douglas-Souza40/fctech-consulta-clima-log/busca_temperatura/internal/handler"
    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
    "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/zipkin"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func main() {
    _ = godotenv.Load()

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
        // ensure flush on exit
        defer func() {
            _ = tp.Shutdown(context.Background())
        }()
    }

    r := mux.NewRouter()
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("busca_temperatura no ar!"))
    }).Methods("GET")

    weatherHandler := &handler.WeatherHandlerBuscaTemp{
        Client: *client.NewClient(),
    }
    // instrument handler with otelhttp
    r.Handle("/temperatura", otelhttp.NewHandler(weatherHandler, "GetTemperatureByCEP")).Methods("GET")

    port := os.Getenv("PORT")
    if port == "" {
        port = "8081"
    }
    log.Printf("Listening busca_temperatura on port %s", port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}
