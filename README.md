# E-Wallet API

A robust RESTful API for an E-Wallet system built with Go, Gin framework, and MySQL.

## Features

- User Authentication (JWT-based)
- Account Management
- Transaction Processing
  - Top Up
  - Payment
  - Transfer between users
- Transaction History
- Balance Management
- Secure PIN Handling

## Tech Stack

- **Language:** Go 1.19+
- **Framework:** Gin
- **Database:** MySQL 8.0+
- **ORM:** GORM
- **Authentication:** JWT
- **ID Generation:** UUID

## Prerequisites

- Go 1.19 or higher
- MySQL 8.0 or higher
- Make (optional, for using Makefile commands)

## Project Setup

### 1. Clone the Repository

```bash
git clone https://github.com/denys89/ewallet-api.git
cd ewallet-api
```

### 2. Set Up Environment Variables

Copy the example environment file and update it with your configurations:

```bash
cp .env.example .env
```

Update the following variables in `.env`:

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=ewallet_db

JWT_SECRET=your_jwt_secret
JWT_EXPIRATION_HOURS=24
REFRESH_TOKEN_SECRET=your_refresh_token_secret
REFRESH_TOKEN_EXPIRATION_DAYS=7

SERVER_PORT=8080
```

### 3. Database Setup

Create a MySQL database:

```sql
CREATE DATABASE ewallet_db;
```

Run the migrations:

```bash
mysql -u your_db_user -p ewallet_db < migrations/000001_migration.sql
```

### 4. Install Dependencies

```bash
go mod download
```

### 5. Run the Application

```bash
go run main.go
```

The API will be available at `http://localhost:8080`

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login user
- `POST /api/v1/auth/refresh-token` - Refresh JWT token

### User Management
- `GET /api/v1/user/profile` - Get user profile
- `PUT /api/v1/user/profile` - Update user profile
- `GET /api/v1/user/balance` - Get user balance

### Transactions
- `POST /api/v1/transactions/topup` - Top up wallet
- `POST /api/v1/transactions/payment` - Make payment
- `POST /api/v1/transactions/transfer` - Transfer to another user
- `GET /api/v1/transactions` - Get transaction history

## Request Examples

### Register User
```json
POST /api/v1/auth/register
{
    "first_name": "John",
    "last_name": "Doe",
    "phone_number": "1234567890",
    "address": "123 Main St",
    "pin": "123456"
}
```

### Top Up
```json
POST /api/v1/transactions/topup
{
    "amount": 100000
}
```

### Transfer
```json
POST /api/v1/transactions/transfer
{
    "amount": 50000,
    "target_user": "recipient_user_id",
    "remarks": "Payment for lunch"
}
```

### Payment
```json
POST /api/v1/transactions/payment
{
    "amount": 25000,
    "remarks": "Coffee payment"
}
```

## Security Features

- JWT-based authentication
- PIN hashing using bcrypt
- UUID for user identification
- Atomic transactions
- SQL injection protection via GORM
- Trusted proxy configuration

## Development

### Code Structure

```
.
├── config/         # Configuration files
├── middleware/     # HTTP middleware
├── migrations/     # Database migrations
├── models/         # Data models
├── repositories/   # Database operations
├── routes/         # HTTP routes
├── main.go        # Application entry point
└── .env           # Environment variables
```
