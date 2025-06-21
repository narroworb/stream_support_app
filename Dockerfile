# syntax=docker/dockerfile:1

FROM golang:1.23.3-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download

COPY . .

RUN go build -o /app/bin/app cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/bin/app .

COPY .env .

COPY web ./web

COPY --from=build /app/migrations ./migrations

EXPOSE 8080

CMD ["./app"]
