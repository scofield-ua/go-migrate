package test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/scofield-ua/go-migrate/db"
)

const TestDbName = "testdb"
const TestDbHost = "localhost"
const TestDbUser = "postgres"
const TestPassword = "postgres"

type SetupResult struct {
	Conn          *pgx.Conn
	DeferFunc     func()
	Err           error
	MigrationsDir string
}

type SetupParams struct {
	T *testing.T
}

func Setup(params SetupParams) SetupResult {
	tmpDir, err := os.MkdirTemp("", "migrations")
	if err != nil {
		params.T.Fatal(err)
	}

	dbhost := os.Getenv("DB_HOST")
	if dbhost == "" {
		dbhost = TestDbHost
	}

	dbuser := os.Getenv("DB_USERNAME")
	if dbuser == "" {
		dbuser = TestDbUser
	}

	dbpwd := os.Getenv("DB_PASSWORD")
	if dbpwd == "" {
		dbpwd = TestPassword
	}

	dbname := os.Getenv("DB_DATABASE")
	if dbname == "" {
		dbname = TestDbName
	}

	conn, err := db.ConnectToDatabase(dbhost, dbuser, dbpwd, dbname)
	if err != nil {
		params.T.Fatal(err)
	}

	// Delete test data
	ctx := context.Background()
	conn.Exec(ctx, `drop table if exists "migrations"`)
	conn.Exec(ctx, `create database `+"\""+dbname+"\"")

	params.T.Log("Dev database has been created")
	db.CreateMigrationsTable(conn)

	return SetupResult{
		Conn: conn,
		DeferFunc: func() {
			conn.Close(ctx)
			os.RemoveAll(tmpDir)
		},
		Err:           err,
		MigrationsDir: tmpDir,
	}
}

func FillUpMigrations(dir string) {
	migrFiles, _ := os.ReadDir(dir)
	for i := range migrFiles {
		fname := migrFiles[i].Name()
		f, err := os.OpenFile(fmt.Sprintf("%s/%s", dir, fname), os.O_WRONLY|os.O_TRUNC, 06444)
		if err != nil {
			continue
		}
		defer f.Close()

		if strings.Contains(fname, "users.up.sql") {
			f.WriteString(`
				create table "users" (
					id int primary key,
					first_name varchar,
					last_name varchar
				)
			`)
		}

		if strings.Contains(fname, "users.down.sql") {
			f.WriteString("drop table users;")
		}

		if strings.Contains(fname, "messages.up.sql") {
			f.WriteString(`
				create table "messages" (
					id int primary key not null,
					user_id int,
					message text
				)
			`)
		}

		if strings.Contains(fname, "messages.down.sql") {
			f.WriteString("drop table messages;")
		}

		if strings.Contains(fname, "articles.up.sql") {
			f.WriteString(`
				create table "articles" (
					id int primary key not null,
					title varchar,
					content text
				)
			`)
		}

		if strings.Contains(fname, "articles.down.sql") {
			f.WriteString("drop table articles;")
		}
	}
}
