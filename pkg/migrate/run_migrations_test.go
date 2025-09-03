package migrate

import (
	"context"
	"testing"
	"time"

	"github.com/scofield-ua/go-migrate/test"
	"github.com/scofield-ua/go-migrate/tools"
)

func TestRunMigrations(t *testing.T) {
	dbConfig := test.TestDbConfig()

	sr := test.Setup(test.SetupParams{
		T: t,
	})
	defer sr.DeferFunc()

	sr.Conn.Exec(context.Background(), `drop table if exists "users", "messages", "articles"`)

	var err error

	for _, n := range []string{"users", "messages", "articles"} {
		time.Sleep(1 * time.Second)

		err = CreateMigration(n, sr.MigrationsDir)
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Log("Migrations filling up")
	test.FillUpMigrations(sr.MigrationsDir)

	t.Log("Running migrations")
	err = RunMigrations(tools.MigrationUp, sr.MigrationsDir, dbConfig)
	if err != nil {
		t.Fatal(err)
	}

	// Check if tables are exist
	t.Log("Check if tables are actually exist")
	var tableExist bool
	for _, table := range []string{"users", "messages", "articles"} {
		err = sr.Conn.QueryRow(context.Background(), `
			SELECT EXISTS (
				SELECT 1 FROM pg_tables WHERE schemaname = 'public' AND tablename = $1
			)
		`, table).Scan(&tableExist)

		if err != nil {
			t.Fatal(err)
		}

		if !tableExist {
			t.Fatal("Migration filled but required table was not created: ", table)
		}
	}

	// Re-run migrations to test if migrations are not going to execute multiple times
	t.Log("Re-running migrations")
	err = RunMigrations(tools.MigrationUp, sr.MigrationsDir, dbConfig)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Test 'down' migrations")
	err = RunMigrations(tools.MigrationDown, sr.MigrationsDir, dbConfig)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Check if tables are actually deleted")
	for _, table := range []string{"users", "messages", "articles"} {
		err = sr.Conn.QueryRow(context.Background(), `
			SELECT EXISTS (
				SELECT 1 FROM pg_tables WHERE schemaname = 'public' AND tablename = $1
			)
		`, table).Scan(&tableExist)

		if err != nil {
			t.Fatal(err)
		}

		if tableExist {
			t.Fatal("Down migration completed but required table is still exist: ", table)
		}
	}

	t.Log("Check if there are no migrations records inside migrations table")
	var n int
	err = sr.Conn.QueryRow(context.Background(), "select count(*) from migrations").Scan(&n)
	if err != nil {
		t.Fatal(err)
	}

	if n > 0 {
		t.Error("Migrations table is not empty after down migrations are done")
	}
}
