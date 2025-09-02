package migrate

import (
	"context"
	"os"
	"testing"

	"github.com/scofield-ua/go-migrate/db"
	"github.com/scofield-ua/go-migrate/test"
	"github.com/scofield-ua/go-migrate/tools"
)

func TestFreshMigration(t *testing.T) {
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
	RunMigrations(tools.MigrationUp, sr.MigrationsDir, sr.Conn)

	err = FreshMigration(sr.MigrationsDir, sr.Conn)
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
	tmpDir, err := os.MkdirTemp("", "migrations")
	if err != nil {
		b.Fatal(err)
	}

	dbhost := os.Getenv("DB_HOST")
	if dbhost == "" {
		dbhost = test.TestDbHost
	}

	dbuser := os.Getenv("DB_USERNAME")
	if dbuser == "" {
		dbuser = test.TestDbUser
	}

	dbpwd := os.Getenv("DB_PASSWORD")
	if dbpwd == "" {
		dbpwd = test.TestPassword
	}

	dbname := os.Getenv("DB_DATABASE")
	if dbname == "" {
		dbname = test.TestDbName
	}

	conn, err := db.ConnectToDatabase(dbhost, dbuser, dbpwd, dbname)
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
		FreshMigration(tmpDir, conn)
	}
}
