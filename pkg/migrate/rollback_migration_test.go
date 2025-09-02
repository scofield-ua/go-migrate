package migrate

import (
	"context"
	"testing"
	"time"

	"github.com/scofield-ua/go-migrate/test"
	"github.com/scofield-ua/go-migrate/tools"
)

func TestRollbackSingleMigration(t *testing.T) {
	sr := test.Setup(test.SetupParams{
		T: t,
	})

	if sr.Err != nil {
		t.Error(sr.Err)
		return
	}
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

	RunMigrations(tools.MigrationUp, sr.MigrationsDir, sr.Conn)

	var oldMigrCount int
	err = sr.Conn.QueryRow(context.Background(), `select count(*) from "migrations"`).Scan(&oldMigrCount)
	if err != nil {
		t.Fatal(err)
	}

	err = RollbackMigration(1, sr.MigrationsDir, sr.Conn)
	if err != nil {
		t.Fatal(err)
	}

	var tableExist bool
	err = sr.Conn.QueryRow(context.Background(), `
		SELECT EXISTS (
    	SELECT 1 FROM pg_tables WHERE schemaname = 'public' AND tablename = 'articles'
		)
	`).Scan(&tableExist)
	if tableExist {
		t.Error("Table ('articles') is still exist after rollback")
	}

	var newMigrCount int
	err = sr.Conn.QueryRow(context.Background(), `select count(*) from "migrations"`).Scan(&newMigrCount)
	if err != nil {
		t.Fatal(err)
	}

	if newMigrCount >= oldMigrCount {
		t.Error("Migration table was not rolledback")
	}
}

func TestRollbackAllMigration(t *testing.T) {
	sr := test.Setup(test.SetupParams{
		T: t,
	})

	if sr.Err != nil {
		t.Error(sr.Err)
		return
	}
	defer sr.DeferFunc()

	sr.Conn.Exec(context.Background(), `drop table if exists "users", "messages", "articles"`)

	var err error

	for _, table := range []string{"users", "messages", "articles"} {
		time.Sleep(1 * time.Second)

		err = CreateMigration(table, sr.MigrationsDir)
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Log("Migrations filling up")
	test.FillUpMigrations(sr.MigrationsDir)

	RunMigrations(tools.MigrationUp, sr.MigrationsDir, sr.Conn)

	RollbackMigration(100, sr.MigrationsDir, sr.Conn)

	t.Log("Check if tables are not exist")
	var tableExist bool
	for _, table := range []string{"users", "messages", "articles"} {
		err = sr.Conn.QueryRow(context.Background(), `
		SELECT EXISTS (
    	SELECT 1 FROM pg_tables WHERE schemaname = 'public' AND tablename = $1
		)
	`, table).Scan(&tableExist)
		if tableExist {
			t.Errorf("Table ('%s') is still exist after rollback", table)
		}
	}

	var migrRowsCount int
	err = sr.Conn.QueryRow(context.Background(), `select count(*) from "migrations"`).Scan(&migrRowsCount)
	if err != nil {
		t.Fatal(err)
	}

	if migrRowsCount != 0 {
		t.Error("Migration table was not rolledback")
	}
}
