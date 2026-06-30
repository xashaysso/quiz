FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o auth_service ./main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/auth_service .
COPY --from=builder /app/migrations ./migrations

EXPOSE 50051
CMD ["./auth_service"]