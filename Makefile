run:
	go run cmd/web/main.go

test:
	go test ./...

generate_migration:
	migrate create -ext sql -dir db/migrations $(name)

db_migrate:
	migrate -database postgres://postgres:postgres@localhost:5432/blog?sslmode=disable  -path db/migrations up

db_rollback:
	migrate -database postgres://postgres:postgres@localhost:5432/blog?sslmode=disable  -path db/migrations down 1

db_test_migrate:
	migrate -database postgres://postgres:postgres@localhost:5433/blog_test?sslmode=disable  -path db/migrations up

db_test_rollback:
	migrate -database postgres://postgres:postgres@localhost:5433/blog_test?sslmode=disable  -path db/migrations down 1

.PHONY: test