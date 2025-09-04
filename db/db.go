package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/scofield-ua/go-migrate/pkg/config"
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

func CreateMigrationsTable(c *config.Config) (bool, error) {
	var tableExist bool
	var err error

	conn, err := ConnectPostgreSQL(c)
	if err != nil {
		return false, err
	}

	err = conn.QueryRow(context.Background(), `
		SELECT EXISTS (
    	SELECT 1 FROM pg_tables WHERE schemaname = 'public' AND tablename = 'migrations'
		)
	`).Scan(&tableExist)

	if err != nil {
		return false, err
	}

	if !tableExist {
		_, err = conn.Exec(context.Background(), migrationsTableSql)
		if err != nil {
			return false, err
		}

		log.Println("Migrations table has been created")
		return true, nil
	} else {
		log.Println("Migrations table already exist")
	}

	return true, nil
}
