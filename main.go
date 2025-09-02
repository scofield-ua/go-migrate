package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/scofield-ua/go-migrate/db"
	"github.com/scofield-ua/go-migrate/pkg/migrate"
	"github.com/scofield-ua/go-migrate/tools"
	"github.com/urfave/cli/v3"
)

var defaultFlags []cli.Flag = []cli.Flag{
	&cli.StringFlag{
		Name:        "dir",
		Value:       "./migrations",
		Usage:       "Migrations folder path",
		DefaultText: "./migrations",
	},
	&cli.StringFlag{
		Name:        "h",
		Value:       "localhost",
		Usage:       "Database host",
		DefaultText: "localhost",
	},
	&cli.StringFlag{
		Name:     "u",
		Usage:    "Database username",
		Required: true,
	},
	&cli.StringFlag{
		Name:     "p",
		Usage:    "Database password",
		Required: true,
	},
	&cli.StringFlag{
		Name:     "db",
		Usage:    "Database name",
		Required: true,
	},
}

func main() {
	// defer conn.Close(context.Background())

	cmd := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create new migration",
				Flags: append(defaultFlags, &cli.StringFlag{
					Name:     "name",
					Usage:    "Migration file name",
					Required: true,
				}),
				Action: func(ctx context.Context, cmd *cli.Command) error {
					conn, err := commandInit(cmd)
					if err != nil {
						return err
					}
					defer conn.Close(context.Background())

					name := cmd.String("name")
					dir := cmd.String("dir")

					if name == "" {
						return fmt.Errorf("name is required")
					}

					if dir == "" {
						dir = "migrations"
					}

					err = migrate.CreateMigration(name, dir)
					if err != nil {
						return fmt.Errorf("name is required")
					}

					log.Printf("Migration file successfully created: %s", name)

					return nil
				},
			},
			{
				Name:  "run",
				Usage: "Apply migrations",
				Flags: defaultFlags,
				Action: func(ctx context.Context, cmd *cli.Command) error {
					conn, err := commandInit(cmd)
					if err != nil {
						return err
					}
					defer conn.Close(context.Background())

					dir := cmd.String("dir")

					if dir == "" {
						dir = "migrations"
					}

					migrate.RunMigrations(tools.MigrationUp, dir, conn)

					return nil
				},
			},
			{
				Name:  "rollback",
				Usage: "Rollback migrations",
				Flags: append(defaultFlags, &cli.IntFlag{
					Name:        "step",
					Value:       1,
					Usage:       "Use step to rollback more than 1 migration (default is 1)",
					DefaultText: "1",
				}),
				Action: func(ctx context.Context, cmd *cli.Command) error {
					conn, err := commandInit(cmd)
					if err != nil {
						return err
					}
					defer conn.Close(context.Background())

					step := cmd.Int("step")
					dir := cmd.String("dir")

					if step == 0 {
						step = 1
					}

					if dir == "" {
						dir = "migrations"
					}

					err = migrate.RollbackMigration(step, dir, conn)
					if err != nil {
						return fmt.Errorf("%v", err)
					}

					return nil
				},
			},
			{
				Name:  "fresh",
				Usage: "Delete all databases and re-run migrations",
				Flags: defaultFlags,
				Action: func(ctx context.Context, cmd *cli.Command) error {
					conn, err := commandInit(cmd)
					if err != nil {
						return err
					}
					defer conn.Close(context.Background())

					dir := cmd.String("dir")

					if dir == "" {
						dir = "migrations"
					}

					err = migrate.FreshMigration(dir, conn)
					if err != nil {
						return fmt.Errorf("%v", err)
					}

					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

// 1. Connect to database
// 2. Create migrations table
func commandInit(cmd *cli.Command) (*pgx.Conn, error) {
	conn, err := establishConnection(cmd.String("h"), cmd.String("u"), cmd.String("p"), cmd.String("db"))
	if err != nil {
		return nil, err
	}

	db.CreateMigrationsTable(conn)

	return conn, nil
}

func establishConnection(h string, u string, p string, d string) (*pgx.Conn, error) {
	conn, err := db.ConnectToDatabase(h, u, p, d)

	if err != nil {
		log.Printf("Error during connection to the database: %v\n", err)
		return nil, err
	}

	return conn, err
}
