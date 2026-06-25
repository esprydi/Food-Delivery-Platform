# Food Delivery Microservices Workspace Rules

## 1. General Architectural Guidelines
* **Versi Stack Stabil:** Gunakan versi 1-2 generasi di belakang yang terbaru (Misal: Golang 1.21, React 18, PostgreSQL 15, RabbitMQ 3.12).
* **Clean Architecture:** Semua service Golang HARUS mengimplementasikan pola Clean Architecture secara ketat.
* **Database-per-service:** Setiap microservice harus memiliki skema/database sendiri. Jangan pernah melakukan query lintas database.
* **Komunikasi Antar Service:**
  * Komunikasi Sinkron: Gunakan REST API.
  * Komunikasi Asinkron: Gunakan RabbitMQ untuk event-driven (mengikuti spesifikasi event di PRD).

## 2. Standar Engineering Golang (Wajib)

### 2.1 Unified API Response Format
Semua HTTP Handler wajib mengembalikan format JSON standar berikut tanpa terkecuali:
```json
{
  "success": true,
  "message": "Pesan deskriptif hasil operasi",
  "data": { ... object atau array ... },
  "error": null
}
```
*Jika terjadi error:* `"success": false`, `"data": null`, dan `"error": "Pesan error teknis/validasi"`.

### 2.2 Golang Context Propagation
Semua method di layer `Usecase` dan `Repository` wajib menerima `ctx context.Context` sebagai parameter pertama:
```go
// Contoh Interface Repository yang Benar
type OrderRepository interface {
    Create(ctx context.Context, order *domain.Order) error
    GetByID(ctx context.Context, id string) (*domain.Order, error)
}
```

### 2.3 Structured Logging
Gunakan package `log/slog` bawaan Golang (1.21+) untuk pencatatan log terstruktur berbasis JSON. Jangan gunakan `fmt.Println` di kode production.

## 3. Aturan Dokumentasi & Sinkronisasi Postman
* Setiap kali agent selesai membuat atau mengubah rute endpoint HTTP, agent **WAJIB** memperbarui file `postman_collection.json` di root repository agar sinkron dengan implementasi kode terbaru.
* Sertakan contoh payload request dan response di dalam dokumentasi Postman tersebut.

## 4. Standar Struktur Folder Golang
* `cmd/api/main.go`: Entry point (Router setup, DB Connect, Dependency Injection).
* `internal/domain/`: Struct Entity, Interface Repository, dan Interface Usecase.
* `internal/usecase/`: Business Logic murni.
* `internal/repository/postgres/`: Implementasi query PostgreSQL.
* `internal/delivery/http/`: HTTP Controller / Handler.
* `migrations/`: File SQL migrasi database (`up.sql` & `down.sql`).
