package migrate

import (
	"context"
	"os"
	"testing"

	"github.com/scofield-ua/go-migrate/test"
	"github.com/scofield-ua/go-migrate/tools"
)

func TestFreshMigration(t *testing.T) {
	dbConfig := test.TestDbConfig()

	sr := test.Setup(test.SetupParams{
		T: t,
	})
	defer sr.DeferFunc()

	sr.Conn.Exec(context.Background(), `drop table if exists "users", "messages"`)

	var err error

	for _, n := range []string{"users", "messages"} {
		err = CreateMigration(n, sr.MigrationsDir)
		if err != nil {
			t.Fatal(err)
		}
	}

	test.FillUpMigrations(sr.MigrationsDir)
	RunMigrations(tools.MigrationUp, sr.MigrationsDir, dbConfig)

	err = FreshMigration(sr.MigrationsDir, dbConfig)
	if err != nil {
		t.Fatal(err)
	}

	var tableExist bool
	for _, n := range []string{"users", "messages"} {
		err = sr.Conn.QueryRow(context.Background(), `
			SELECT EXISTS (
    		SELECT 1 FROM pg_tables WHERE schemaname = 'public' AND tablename = $1
			)
		`, n).Scan(&tableExist)
		if !tableExist {
			t.Errorf("Table ('%s') does not exist after fresh migration", n)
		}
	}
}

func BenchmarkFreshMigration(b *testing.B) {
	dbConfig := test.TestDbConfig()

	tmpDir, err := os.MkdirTemp("", "migrations")
	if err != nil {
		b.Fatal(err)
	}

	for _, n := range []string{"users", "messages"} {
		err = CreateMigration(n, tmpDir)
		if err != nil {
			b.Fatal(err)
		}
	}

	test.FillUpMigrations(tmpDir)

	for b.Loop() {
		FreshMigration(tmpDir, dbConfig)
	}
}
