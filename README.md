## Description

Blog source code

# Running locally

Start database

```
docker compose up -d
```

Install `migrate` CLI

```
make install_golang_migrate
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

Run migration on a custom database

```
make db_migrate db_url='DATABASE_URL'
```

Rollback 1 migration on dev database

```
make db_rollback
```

Rollback 1 migration on test database

```
make db_test_rollback
```

Rollback 1 migration on a custom database

```
make db_rollback db_url='DATABASE_URL'
```

Force migration version on dev database

In case there is a migration error and you receive the message `Dirty database version 20211109080746. Fix and force version.`. You can fix it by running the following command. Be sure to pass the previous version and not the one specified in the error.

```
make db_force version=20211108120029
```

Force migration version on test database

```
make db_test_force version=20211108120029
```

# Run with docker locally

Build the image

```
docker build -t blog .
```

Change the POSTGRES_URL to point to the docker container on the .env file

```
POSTGRES_URL=postgres://postgres:postgres@blog-postgres-1:5432/blog?sslmode=disable
```

Run using docker

```
docker run --rm  --env-file .env -p 3000:3000 --network blog_default blog
```
