FROM golang:1.21-alpine3.18 AS builder
WORKDIR /app
COPY go.mod ./
# Copy go.sum if it exists, otherwise skip (touch avoids error)
RUN touch go.sum
COPY go.sum ./
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

FROM alpine:3.18
WORKDIR /root/
COPY --from=builder /app/main .
# Fallback for .env
RUN touch .env
COPY --from=builder /app/.env .
EXPOSE 8085
CMD ["./main"]
