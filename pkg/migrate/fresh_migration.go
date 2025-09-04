package migrate

import (
	"context"
	"log"

	"github.com/scofield-ua/go-migrate/db"
	"github.com/scofield-ua/go-migrate/pkg/config"
	"github.com/scofield-ua/go-migrate/tools"
)

const deleteAllTablesSql = `
  DO $$
  DECLARE
		r RECORD;
  BEGIN
		FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = current_schema()) LOOP
			EXECUTE 'DROP TABLE IF EXISTS ' || quote_ident(r.tablename) || ' CASCADE';
		END LOOP;
	END
	$$;
`

func FreshMigration(dir string, c *config.Config) error {
	conn, err := db.ConnectPostgreSQL(c)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(), deleteAllTablesSql)
	if err != nil {
		log.Fatal(err)
		return err
	}

	db.CreateMigrationsTable(c)

	RunMigrations(tools.MigrationUp, dir, c)

	return nil
}
