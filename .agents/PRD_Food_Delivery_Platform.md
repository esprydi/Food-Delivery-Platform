# Product Requirements Document (PRD): Food Delivery Platform (Production-Grade)

## 1. Product Overview
Food Delivery Platform adalah sistem terdistribusi berbasis Microservices yang menghubungkan Pelanggan (Customer), Restoran (Merchant), dan Kurir (Driver). Sistem ini dirancang untuk menangani konkurensi tinggi, komunikasi asinkron berbasis event, serta pemisahan domain data yang ketat.

## 2. Tech Stack & Tools
* **Backend Development:** Golang (v1.21)
* **Frontend Development:** React JS 18 (Vite 5 + Tailwind CSS v3)
* **Database:** PostgreSQL 15 (Database-per-service pattern)
* **Message Broker:** RabbitMQ 3.12 (Event-Driven Asynchronous Messaging)
* **Architecture Style:** Microservices & Clean Architecture
* **API Design:** RESTful API (JSON)
* **API Documentation & Testing:** Postman
* **Containerization & Orchestration:** Docker & Docker Compose

---

## 3. Microservices Domain & Data Dictionary

Setiap microservice memiliki database PostgreSQL terpisah. Berikut adalah spesifikasi entitas datanya:

### 3.1 User Service (`user_db`)
Mengelola data pengguna, autentikasi (JWT), dan hak akses (Role).
* **Table `users`:**
  * `id` (UUID, Primary Key)
  * `name` (VARCHAR)
  * `email` (VARCHAR, Unique)
  * `password_hash` (VARCHAR)
  * `role` (ENUM: `CUSTOMER`, `MERCHANT`, `DRIVER`, `ADMIN`)
  * `phone` (VARCHAR)
  * `created_at` (TIMESTAMP), `updated_at` (TIMESTAMP)

### 3.2 Catalog Service (`catalog_db`)
Mengelola profil restoran dan menu makanan.
* **Table `restaurants`:**
  * `id` (UUID, Primary Key)
  * `owner_id` (UUID, Ref -> User ID)
  * `name` (VARCHAR)
  * `address` (TEXT)
  * `is_open` (BOOLEAN)
* **Table `menu_items`:**
  * `id` (UUID, Primary Key)
  * `restaurant_id` (UUID, Foreign Key)
  * `name` (VARCHAR)
  * `description` (TEXT)
  * `price` (DECIMAL(12,2))
  * `is_available` (BOOLEAN)

### 3.3 Order Service (`order_db`)
Core State Machine pesanan.
* **Table `orders`:**
  * `id` (UUID, Primary Key)
  * `customer_id` (UUID)
  * `restaurant_id` (UUID)
  * `status` (VARCHAR: `PENDING`, `PAID`, `PREPARING`, `READY_FOR_PICKUP`, `DELIVERING`, `COMPLETED`, `CANCELED`)
  * `total_amount` (DECIMAL(12,2))
  * `delivery_address` (TEXT)
  * `created_at` (TIMESTAMP), `updated_at` (TIMESTAMP)
* **Table `order_items`:**
  * `id` (UUID, Primary Key)
  * `order_id` (UUID, Foreign Key)
  * `menu_item_id` (UUID)
  * `menu_item_name` (VARCHAR) -- Snapshot nama menu saat pesanan dibuat
  * `quantity` (INTEGER)
  * `unit_price` (DECIMAL(12,2))

---

## 4. Kontrak Event Message Broker (RabbitMQ)

Semua komunikasi asinkron menggunakan format JSON dengan exchange type `topic`.

### 4.1 Event: `order.created`
Diterbitkan oleh **Order Service** saat checkout berhasil.
```json
{
  "event_id": "uuid-v4",
  "event_type": "ORDER_CREATED",
  "timestamp": "2026-06-26T10:00:00Z",
  "payload": {
    "order_id": "ord-123-uuid",
    "customer_id": "cust-456-uuid",
    "restaurant_id": "rest-789-uuid",
    "total_amount": 150000.00
  }
}
```

### 4.2 Event: `payment.success`
Diterbitkan oleh **Payment Service** setelah pembayaran berhasil diverifikasi.
```json
{
  "event_id": "uuid-v4",
  "event_type": "PAYMENT_SUCCESS",
  "timestamp": "2026-06-26T10:01:00Z",
  "payload": {
    "order_id": "ord-123-uuid",
    "payment_id": "pay-999-uuid",
    "paid_amount": 150000.00
  }
}
```

---

## 5. Spesifikasi Endpoint REST API

Semua endpoint menggunakan prefix `/api/v1`.

### 5.1 User Service (`:8081`)
* `POST /api/v1/auth/register` - Registrasi user baru
* `POST /api/v1/auth/login` - Login & mengembalikan JWT Token
* `GET /api/v1/users/me` - Ambil profil user login (Requires Bearer Token)

### 5.2 Catalog Service (`:8082`)
* `GET /api/v1/restaurants` - Daftar restoran aktif
* `GET /api/v1/restaurants/{id}/menus` - Daftar menu dari restoran tertentu
* `POST /api/v1/merchant/menus` - Tambah menu baru (Requires Role: MERCHANT)

### 5.3 Order Service (`:8083`)
* `POST /api/v1/orders` - Buat pesanan baru
* `GET /api/v1/orders/{id}` - Detail pesanan & status tracking
* `GET /api/v1/orders/customer` - Riwayat pesanan milik customer login
* `PATCH /api/v1/orders/{id}/status` - Update status pesanan (Requires Role: MERCHANT/DRIVER)

---

## 6. Persyaratan Kualitas & Pengujian (Postman)
1. **Automated Testing:** Setiap endpoint harus dilengkapi test Postman Script (`pm.test`) untuk memvalidasi HTTP status code dan JSON Schema.
2. **Postman Collection Sync:** Seluruh fungsionalitas harus diekspor ke dalam file `postman_collection.json` di root repositori.
