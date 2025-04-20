FROM golang:1.23 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . ./

# Собираем единый бинарь, который запускает HTTP и gRPC сервера
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /server ./cmd/service

# ---
FROM gcr.io/distroless/static
COPY --from=builder /server /server
ENTRYPOINT ["/server"]
