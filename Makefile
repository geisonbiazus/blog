run:
	go run cmd/web/main.go

test:
	go test ./...

db_url ?= postgres://postgres:postgres@localhost:5432/blog?sslmode=disable

install_golang_migrate:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

generate_migration:
	migrate create -ext sql -dir db/migrations $(name)

db_migrate:
	migrate -database $(db_url) -path db/migrations up

db_rollback:
	migrate -database $(db_url) -path db/migrations down 1

db_force:
	migrate -database $(db_url) -path db/migrations force $(version)

db_test_migrate:
	migrate -database postgres://postgres:postgres@localhost:5433/blog_test?sslmode=disable  -path db/migrations up

db_test_rollback:
	migrate -database postgres://postgres:postgres@localhost:5433/blog_test?sslmode=disable  -path db/migrations down 1

db_test_force:
	migrate -database postgres://postgres:postgres@localhost:5433/blog_test?sslmode=disable -path db/migrations force $(version)

send:
	echo $(m)

.PHONY: test