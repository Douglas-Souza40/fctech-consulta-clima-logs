# Run with Docker Compose

This project contains two services:

- `busca_temperatura` (port 8081) — looks up city from ViaCEP and current temperature from WeatherAPI.
- `processa_cep` (port 8080) — accepts POST JSON { "cep":"<8 digits>" } and returns location + temperatures.

Prerequisites:

- Docker and Docker Compose installed.
- A `.env` file at the project root with the variable `WEATHER_API_KEY` (example provided in repo).

Start the services (builds images):

```powershell
docker compose up --build
```

Notes:

- The compose file forwards host ports 8080 -> `processa_cep` and 8081 -> `busca_temperatura`.
- `processa_cep` is configured to call `http://busca_temperatura:8081` inside the compose network.
- To override the `WEATHER_API_KEY` at runtime you can either edit `.env` or pass an env var when starting:

```powershell
$env:WEATHER_API_KEY="your_key_here"; docker compose up --build
```

Quick test (after both services are up):

```powershell
curl -X POST http://localhost:8080/weather -H "Content-Type: application/json" -d '{"cep":"29902555"}'
```
