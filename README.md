cat > README.md <<'EOF'
# Task 3 - Go Backend API

A Go-based backend API implementing asynchronous CSV/Excel data import, JWT authentication with role-based access control, product/inventory management, and database health monitoring.

## Tech Stack

- Go
- Gin
- GORM
- SQLite
- JWT
- bcrypt
- Excelize
- godotenv
- Goroutines
- Background scheduler

## Project Structure

```text
task3-go/
├── cmd/server/main.go
├── internal/
│   ├── database/
│   ├── handlers/
│   ├── middleware/
│   ├── models/
│   └── services/
├── sample_products.csv
├── sample_products.xlsx
├── .env
├── .gitignore
├── go.mod
├── go.sum
└── README.md