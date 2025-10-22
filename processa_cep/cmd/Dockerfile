FROM golang:1.24.2-alpine AS builder
WORKDIR /src
RUN apk add --no-cache git
COPY go.mod go.sum /src/
RUN go env -w GOPROXY=https://proxy.golang.org
COPY . /src
RUN go build -o /out/processa_cep ./processa_cep/cmd

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=builder /out/processa_cep /usr/local/bin/processa_cep
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/processa_cep"]
