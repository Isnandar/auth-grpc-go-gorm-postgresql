Auth GRPC Go Gorm PostgreSQL

Tutorial Run
- copy and paste .env and setup it (make sure about both db and redis connection is on)
- go mod tidy (for package installation)
- go run ./cmd/migrate.go (for migration and seeding db)
- go run ./cmd/main.go (for running system)