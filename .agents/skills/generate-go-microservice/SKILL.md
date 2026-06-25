---
name: generate-go-microservice
description: Menginisialisasi scaffold untuk microservice Golang production-grade (Clean Arch + Docker + Makefile)
---

# Generate Golang Microservice (Production-Grade)

Saat user memintamu untuk membuat atau menginisialisasi microservice baru, ikuti urutan eksekusi berikut secara presisi:

1. **Buat Folder Service**: Buat direktori root service (misal: `order-service`).
2. **Inisialisasi Go Mod**: Jalankan `go mod init <nama-module>`.
3. **Generate Struktur Folder Clean Architecture**:
   - `cmd/api/`
   - `internal/domain/`
   - `internal/usecase/`
   - `internal/delivery/http/`
   - `internal/repository/postgres/`
   - `internal/infrastructure/`
   - `config/`
   - `migrations/`
4. **Buat File Konfigurasi**: Buat `config/config.go` dan `.env.example`.
5. **Implementasi Domain Entity**: Buat `internal/domain/domain.go` sesuai spesifikasi Kamus Data di `PRD_Food_Delivery_Platform.md`. Pastikan struct dilengkapi tag `json:"..."` dan `db:"..."`.
6. **Implementasi Skeleton Layering**:
   - Repository dengan Context: `Create(ctx context.Context, ...)`
   - Usecase dengan Context: `Execute(ctx context.Context, ...)`
   - HTTP Handler dengan Unified JSON Response standar dari `AGENTS.md`.
7. **Generate Dockerfile (Multi-Stage Build)**:
   Buat file `Dockerfile` di root service:
   ```dockerfile
   FROM golang:1.21-alpine3.18 AS builder
   WORKDIR /app
   COPY go.mod go.sum ./
   RUN go mod download
   COPY . .
   RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

   FROM alpine:3.18
   WORKDIR /root/
   COPY --from=builder /app/main .
   COPY --from=builder /app/.env .
   EXPOSE 8080
   CMD ["./main"]
   ```
8. **Generate Makefile**:
   Buat `Makefile` di root service:
   ```makefile
   run:
   	go run cmd/api/main.go
   test:
   	go test -v ./...
   build:
   	go build -o bin/main cmd/api/main.go
   docker:
   	docker build -t $(SERVICE_NAME):latest .
   ```
9. **Wiring Entry Point (`cmd/api/main.go`)**: Hubungkan semua komponen menggunakan Dependency Injection eksplisit, jalankan server HTTP, dan tangani Graceful Shutdown.
10. **Sync Postman**: Beritahu sistem untuk memperbarui `postman_collection.json`.
