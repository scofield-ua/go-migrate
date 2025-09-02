package migrate

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/scofield-ua/go-migrate/tools"
)

func RunMigrations(variant tools.MigrationVariant, migrDir string, conn *pgx.Conn) error {
	dirPath, _ := filepath.Abs(migrDir)

	migrFiles, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	if len(migrFiles) == 0 {
		log.Print("No migration files to run")
		return nil
	}

	migratedFiles := 0
	for _, migrFile := range migrFiles {
		if !strings.Contains(migrFile.Name(), string(variant)) {
			continue
		}

		fp := fmt.Sprintf("%s/%s", migrDir, migrFile.Name())

		ok, err := RunMigration(fp, conn)
		if err != nil {
			return err
		}

		if ok {
			migratedFiles++
		}
	}

	if migratedFiles == 0 {
		log.Println("Nothing to migrate")
	}

	return nil
}

// Return true only when successfully migrated
// false when migration already exist
func RunMigration(path string, conn *pgx.Conn) (bool, error) {
	base := filepath.Base(path)

	fbytes, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}

	if len(fbytes) == 0 {
		log.Printf("Empty file: %s", path)
		return false, nil
	}

	sql := string(fbytes)

	// If it's down migration
	if strings.Contains(base, fmt.Sprintf(".%s.", tools.MigrationDown.String())) {
		_, err = conn.Exec(context.Background(), sql)
		if err != nil {
			log.Fatal(err)
			return false, err
		}

		// Delete "up" record from migrations table
		upMigrName := tools.ChangeMigrationVariant(base, tools.MigrationUp)
		_, err = conn.Exec(
			context.Background(),
			`delete from "migrations" where "migration" like $1`,
			"%"+upMigrName)
		if err != nil {
			log.Fatal(err)
			return false, err
		}

		return true, nil
	}

	var migrExist int
	err = conn.QueryRow(context.Background(), `select id from migrations where migration = $1`, base).Scan(&migrExist)
	if err != nil && err != pgx.ErrNoRows {
		log.Fatal(err)
		return false, err
	}

	if migrExist > 0 {
		// log.Printf("Migration exist: %s", path)
		return false, nil
	}

	log.Printf("Running migration: %s", path)

	_, err = conn.Exec(context.Background(), sql)
	if err != nil {
		log.Fatal(err)
		return false, err
	}

	// Save succeed migration operation into database
	_, err = conn.Exec(context.Background(), `insert into "migrations" (migration) values ($1)`, base)
	if err != nil {
		log.Fatal(err)
	}

	return true, nil
}
