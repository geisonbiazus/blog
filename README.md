## Description

Blog source code

# Running locally

Start database

```
docker compose up -d
```

Install `migrate` CLI

```
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Migrate database

```
make db_migrate
make db_test_migrate
```

Run tests

```
make test
```

Start the project locally

```
make run
```

## Migrations

Create migration

```
make generate_migration name=migration_name
```

Run migrations on dev database

```
make db_migrate
```

Run migrations on test database

```
make db_test_migrate
```

Rollback 1 migration on dev database

```
make db_rollback
```

Rollback 1 migration on test database

```
make db_test_rollback
```
