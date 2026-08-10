## Goose

These and the sqlc.md docs are written to reference the chirpy projects

_<https://github.com/pressly/goose>_

  Goose is a database migration cli tool written in Go
  [Installing](https://github.com/pressly/goose#install) using go

```bash
  go install github.com/pressly/goose/v3/cmd/goose@latest
```

### Migrations

A "migration" in Goose is just a .sql file with some SQL queries and some special comments.
he simplest format for these files is:

```bash
  number_name.sql
  # example
  # 001_users.sql
```

Where the contents of 001_users.sql schema is

```sql
  -- +goose Up
  CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    email TEXT NOT NULL UNIQUE
  );

  -- +goose Down
  DROP TABLE users;
```

### To run migrations

cd into the `sql/schema` directory and run:

- Migrate Up, to bring db up to date

  ```bash
    goose postgres "db_connection_string" up

    # example:
    # goose postgres "postgres://postgres:postgres@localhost:5432/chirpy" up
    # db connection string format is: protocol://username:password@host:port/database
  ```

- Migrate Down, to delete schema changes

  Running multiple migrate downs will go one schema down

  ```bash
    goose postgres "db_connection_string" down

    # example:
    # goose postgres "postgres://postgres:postgres@localhost:5432/chirpy" down
  ```
