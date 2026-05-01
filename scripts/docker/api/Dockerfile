## build
FROM golang:1.26.1 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY cmd/main.go .
RUN go build -o /ledger-config

## deploy
FROM ubuntu AS base

WORKDIR /

COPY --from=build /ledger-config /ledger-config

#AWS 
ENV AWS_ADDRESS="http://192.168.49.2:4566"
ENV AWS_REGION="us-east-1"

EXPOSE 5000

ENTRYPOINT [ "/ledger-config" ]