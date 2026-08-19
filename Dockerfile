# syntax=docker/dockerfile:1

FROM golang:1.23-alpine AS builder

WORKDIR /src

COPY main.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/http-server-projeto-korp main.go

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /out/http-server-projeto-korp /app/http-server-projeto-korp

EXPOSE 8080

USER 65532:65532

ENTRYPOINT ["/app/http-server-projeto-korp"]

