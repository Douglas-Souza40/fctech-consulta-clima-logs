FROM golang:1.24.2-alpine AS builder
WORKDIR /src
RUN apk add --no-cache git
COPY go.mod go.sum /src/
RUN go env -w GOPROXY=https://proxy.golang.org
COPY . /src
RUN go build -o /out/busca_temperatura ./busca_temperatura/cmd

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=builder /out/busca_temperatura /usr/local/bin/busca_temperatura
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser
EXPOSE 8081
ENTRYPOINT ["/usr/local/bin/busca_temperatura"]
