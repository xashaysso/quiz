FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o stats_app ./main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/stats_app .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8081
CMD ["./stats_app"]