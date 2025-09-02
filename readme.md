### Go migration tool for PostgreSQL

Includes `migrations` table for migrations history. Create migration, add SQL code for migration, run it. Rollback if needed or re-run migrations from scratch. CLI and library interface.

#### CLI Commands:

- **`migrate create --name=create_users_table`**

  This command will create 2 migration files with name `create_users_table.up.sql` and `create_users_table.down.sql`. `.down.sql` one going to be used for rollbacks

  **Parameters:**

  - `name` - name of the migration
  - `dir` - directory where all migration file are store (default path is `./migrations`)

- **`migrate run`**

  Run all migrations.

  **Parameters:**

  - `dir` - directory where all migration file are store (default path is `./migrations`)

- **`migrate rollback`**

  Rollback a single (or more) migration.

  **Parameters:**

  - `step` - how many migrations we need to rollback (default is 1)
  - `dir` - directory where all migration file are store (default path is `./migrations`)

- **`migrate fresh`**

  Delete all database tables and re-run migrations.

  **Parameters:**

  - `dir` - directory where all migration file are store (default path is `./migrations`)
