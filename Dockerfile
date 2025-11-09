# syntax=docker/dockerfile:1.7
FROM golang:1.24.4-alpine AS build
WORKDIR /app
RUN apk add --no-cache git ca-certificates
COPY go.mod go.sum ./
RUN  go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o httpapi ./cmd/httpapi

FROM alpine:3.20
WORKDIR /app
COPY --from=build /app/httpapi /app/httpapi
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["./httpapi"]
