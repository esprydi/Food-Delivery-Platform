# Food Delivery Microservices Platform

A production-grade, event-driven Food Delivery SaaS platform built with a microservices architecture. This project serves as a comprehensive portfolio demonstrating modern backend engineering practices, robust system design, and full-stack integration.

## 🚀 Architecture Overview

The platform consists of 4 independent Go microservices and a React frontend, communicating synchronously via REST APIs and asynchronously via RabbitMQ.

- **User Service:** Manages authentication, authorization (JWT), and user profiles (Customer & Merchant roles).
- **Catalog Service:** Manages restaurant profiles and menu items. Used by merchants to add products and customers to browse.
- **Order Service:** Handles order creation, cart calculation, and order tracking.
- **Payment Service:** Integrates with Midtrans Sandbox for processing payments. Listens to order events and broadcasts payment success events.

### Tech Stack
- **Backend:** Golang (v1.21), Echo Framework, GORM
- **Frontend:** React 18, Vite, Vanilla CSS (Glassmorphism design)
- **Database:** PostgreSQL (Database-per-service logical separation)
- **Message Broker:** RabbitMQ
- **Payment Gateway:** Midtrans (Sandbox)
- **Containerization:** Docker & Docker Compose

## 📁 Repository Structure

```text
portofolio-golang-microservice/
├── user-service/       # Go Microservice (Port 8081)
├── catalog-service/    # Go Microservice (Port 8082)
├── order-service/      # Go Microservice (Port 8083)
├── payment-service/    # Go Microservice (Port 8084)
├── frontend-app/       # React/Vite Frontend (Port 5173)
├── fd_postgres_data/   # PostgreSQL Volume Mount
└── docker-compose.yml  # Orchestrates all backend services
```

*Note: All Go microservices strictly adhere to the Clean Architecture pattern (Domain, Repository, Usecase, Delivery).*

## 🛠️ Prerequisites

Ensure you have the following installed:
- [Docker](https://docs.docker.com/get-docker/) & Docker Compose
- [Node.js](https://nodejs.org/en/) (v18+) & npm
- [Go](https://golang.org/dl/) 1.21 (if running backend locally without Docker)

## ⚙️ Getting Started

### 1. Run the Backend Infrastructure

The entire backend infrastructure (Microservices, PostgreSQL, and RabbitMQ) is containerized.

```bash
# Start all services in detached mode
docker-compose up -d --build
```

You can verify the containers are running with `docker ps`. The services will be exposed on the following ports:
- User Service: `localhost:8081`
- Catalog Service: `localhost:8082`
- Order Service: `localhost:8083`
- Payment Service: `localhost:8084`

### 2. Run the Frontend App

Open a new terminal and navigate to the frontend directory:

```bash
cd frontend-app
npm install
npm run dev
```

The React app will be accessible at `http://localhost:5173`.

## 🧪 Testing the E2E Flow

1. **Merchant Registration:** Go to `/auth`, select "Register", choose "Merchant" role.
2. **Setup Restaurant:** Login as Merchant, you will be directed to the Merchant Dashboard. Create your restaurant and add some menu items with prices.
3. **Customer Registration:** Open an incognito window, register a new account as a "Customer".
4. **Order Flow:** Login as Customer. You will see the Home page with available restaurants. Click a restaurant, view menus, enter quantities, and click "Checkout".
5. **Payment Flow:** You will be redirected to the Orders page. Click "Pay Now" to open the Midtrans payment simulation page. Use Midtrans sandbox credentials to complete the payment.

## 🔗 Midtrans Integration Note
This project uses Midtrans Sandbox. Ensure that the Server Key in `payment-service/.env` is valid. All successful payments will trigger a RabbitMQ event from `payment-service` to `order-service` to automatically update the order status.

## 📝 License
MIT License
