# Guia: Docker + Zipkin — Consulta Clima (FCTech)

Este documento descreve passo-a-passo como subir a aplicação com Docker Compose (incluindo Zipkin), executar os serviços e validar traces no Zipkin.

Pré-requisitos

- Docker (Desktop/Engine) instalado
- Docker Compose v1.27+ (ou `docker compose` integrado ao Docker Desktop)
- Git

Estrutura usada

- `processa_cep` — Serviço A (porta 8080) — aceita POST JSON { "cep":"<8 digits>" } e consulta o serviço B.
- `busca_temperatura` — Serviço B (porta 8081) — consulta ViaCEP e WeatherAPI; responde cidade + temp em C/F/K.
- `zipkin` — serviço Zipkin (porta 9411) criado pelo `docker-compose.yml` para receber traces.

1) Clonar o repositório

```powershell
git clone https://github.com/Douglas-Souza40/fctech-consulta-clima-logs.git
cd fctech-consulta-clima-logs
```

2) Criar o arquivo `.env`

Crie um arquivo `.env` na raiz do repositório (o projeto já usa `godotenv`) com a chave do WeatherAPI:

```
WEATHER_API_KEY=your_weather_api_key_here
```

Observação: não versionar esse arquivo se contiver segredos.

3) Subir o stack (inclui Zipkin)

```powershell
docker compose up --build
```

O compose irá:
- construir as imagens de `processa_cep` e `busca_temperatura`;
- iniciar um container `zipkin` (openzipkin/zipkin) em `http://localhost:9411`;
- iniciar `busca_temperatura` (porta 8081) e `processa_cep` (porta 8080).

Nota: se as portas já estiverem em uso, pare o processo local que as usa ou altere os mapeamentos em `docker-compose.yml`.

4) Testar o fluxo end-to-end

- Exemplo usando PowerShell (Invoke-RestMethod) — retorna JSON:

```powershell
Invoke-RestMethod -Method Post -Uri 'http://localhost:8080/weather' -Body '{"cep":"29902555"}' -ContentType 'application/json'
```

- Exemplo usando curl:

```powershell
curl -X POST http://localhost:8080/weather -H "Content-Type: application/json" -d '{"cep":"29902555"}'
```

Você deve receber JSON com city, temp_c, temp_f e temp_k.

5) Validar traces no Zipkin (UI)

1. Abra o navegador em: http://localhost:9411
2. Em "Find Traces" selecione o Service: `processa_cep` (ou `busca_temperatura`).
3. No Span Name coloque, por exemplo, `GetTemperatureByCEP`, `GetLocationByCEP` ou `GetTemperatureByCity`.
4. Nas Tags use as chaves que instrumentamos (por exemplo `cep=29902555` ou `city=NomeDaCidade`).
5. Ajuste o Lookback (Last 1 hour) e clique em "Find Traces".

Clique num trace para ver spans encadeados entre `processa_cep` e `busca_temperatura`.

6) Validar traces via API (curl)

- Listar services:

```powershell
curl http://localhost:9411/api/v2/services
```

- Listar spans de um serviço:

```powershell
curl "http://localhost:9411/api/v2/spans?serviceName=processa_cep"
```

- Buscar traces por service + span + tag (annotationQuery):

```powershell
curl "http://localhost:9411/api/v2/traces?serviceName=processa_cep&spanName=GetTemperatureByCEP&annotationQuery=cep=29902555&limit=10"
```

- Buscar trace por trace_id (use trace_id que aparecer nos logs se disponível):

```powershell
curl http://localhost:9411/api/v2/trace/<TRACE_ID>
```

7) Logs e correlação

- Em caso de erros o serviço inclui `trace_id` e `span_id` no log para facilitar correlação. Copie o `trace_id` e cole no campo "Trace ID" da UI do Zipkin para puxar o trace completo.

8) Parar e limpar

```powershell
docker compose down
```

Extras e dicas

- Se quiser ver Zipkin junto com outros serviços em produção, considere mover para collector OTLP/Jaeger e usar um backend como Tempo/Jaeger/OTel-collector.
- Para adicionar mais atributos (HTTP status, duration, etc.) ou logs correlacionados em cada request, eu posso instrumentar mais pontos no código.

Se quiser, eu gero um script `scripts/run-with-zipkin.ps1` que:
- verifica Docker
- executa `docker compose up --build -d`
- espera Zipkin responder em /health
- faz um request de teste e abre o navegador na UI do Zipkin.

Fim.
