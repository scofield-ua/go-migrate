package test

import (
	"context"
	"testing"

	"github.com/scofield-ua/go-migrate/db"
)

func TestCreateMigrationsTable(t *testing.T) {
	dbConfig := TestDbConfig()

	sr := Setup(SetupParams{
		T: t,
	})
	defer sr.DeferFunc()

	sr.Conn.Exec(context.Background(), `drop table if exists "migrations"`)

	var err error

	res, err := db.CreateMigrationsTable(dbConfig)
	if res != true {
		t.Fatalf("migrations table has not created: %v", err)
	}

	var exist bool
	sr.Conn.QueryRow(context.Background(), `
		SELECT EXISTS (
			SELECT 1 FROM pg_tables WHERE schemaname = 'public' AND tablename = 'migrations'
		)
	`).Scan(&exist)

	if !exist {
		t.Fatal("migrations table has not actually created after 'CreateMigrationsTable' result was sent as 'true'")
	}

	// Re-run function to make sure that second run does not produce any errors
	_, err = db.CreateMigrationsTable(dbConfig)
	if err != nil {
		t.Fatal(err)
	}
}
