package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/scofield-ua/go-migrate/config"
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
		Sources:     cli.EnvVars("DB_HOST"),
	},
	&cli.StringFlag{
		Name:     "u",
		Usage:    "Database username",
		Required: true,
		Sources:  cli.EnvVars("DB_USERNAME"),
	},
	&cli.StringFlag{
		Name:     "p",
		Usage:    "Database password",
		Required: true,
		Sources:  cli.EnvVars("DB_PASSWORD"),
	},
	&cli.StringFlag{
		Name:     "db",
		Usage:    "Database name",
		Required: true,
		Sources:  cli.EnvVars("DB_NAME"),
	},
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

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
					_, err := commandInit(cmd)
					if err != nil {
						return err
					}

					name := cmd.String("name")
					dir := cmd.String("dir")

					err = migrate.CreateMigration(name, dir)
					if err != nil {
						return err
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
					dbConfig, err := commandInit(cmd)
					if err != nil {
						return err
					}

					dir := cmd.String("dir")

					migrate.RunMigrations(tools.MigrationUp, dir, dbConfig)

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
					dbConfig, err := commandInit(cmd)
					if err != nil {
						return err
					}

					step := cmd.Int("step")
					dir := cmd.String("dir")

					if step == 0 {
						step = 1
					}

					err = migrate.RollbackMigration(step, dir, dbConfig)
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
					dbConfig, err := commandInit(cmd)
					if err != nil {
						return err
					}

					dir := cmd.String("dir")

					err = migrate.FreshMigration(dir, dbConfig)
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

// Init database config and check if we need to create migrations table
func commandInit(cmd *cli.Command) (*config.Config, error) {
	dbConfig := config.Config{}
	dbConfig.DB.SetHost(cmd.String("h"))
	dbConfig.DB.SetUsername(cmd.String("u"))
	dbConfig.DB.SetPassword(cmd.String("p"))
	dbConfig.DB.SetDbName(cmd.String("db"))

	conn, err := db.ConnectPostgreSQL(&dbConfig)
	if err != nil {
		return nil, err
	}
	defer conn.Close(context.Background())

	db.CreateMigrationsTable(conn)

	return &dbConfig, nil
}
