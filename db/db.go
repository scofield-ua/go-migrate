package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

const migrationsTableSql string = `
	create table if not exists migrations (
  	id int primary key generated always as identity,
  	migration varchar not null,
  	created_at timestamp not null default current_timestamp
	)
`

func ConnectToDatabase(h string, u string, p string, d string) (*pgx.Conn, error) {
	ctx := context.Background()

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s", h, u, p, d)

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func CreateMigrationsTable(conn *pgx.Conn) {
	var tableExist bool

	err := conn.QueryRow(context.Background(), `
		SELECT EXISTS (
    	SELECT 1 FROM pg_tables WHERE schemaname = 'public' AND tablename = 'migrations'
		)
	`).Scan(&tableExist)

	if err != nil {
		log.Fatal(err)
	}

	if !tableExist {
		log.Println("Creating migrations table...")

		conn.Exec(context.Background(), migrationsTableSql)
	}
}
