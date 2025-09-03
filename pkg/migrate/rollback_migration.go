package migrate

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/scofield-ua/go-migrate/config"
	"github.com/scofield-ua/go-migrate/db"
	"github.com/scofield-ua/go-migrate/tools"
)

func RollbackMigration(step int, dir string, config *config.Config) error {
	conn, err := db.ConnectPostgreSQL(config)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	rolledback := 0
	dirPath, _ := filepath.Abs(dir)

	rows, err := conn.Query(context.Background(), `select migration from migrations order by id desc limit $1`, step)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer rows.Close()

	records := make([]string, 0)
	for rows.Next() {
		var upMigrName string
		err := rows.Scan(&upMigrName)
		if err != nil {
			log.Fatal(err)
			return err
		}

		records = append(records, upMigrName)
	}

	if len(records) == 0 {
		log.Println("Nothing to rollback")
	}

	for _, upMigrName := range records {
		downMigr := tools.ChangeMigrationVariant(upMigrName, tools.MigrationDown)

		items, err := os.ReadDir(dirPath)
		if err != nil {
			log.Fatal(err)
			return err
		}

		var fbytes []byte

		// Find down migration file
		for _, f := range items {
			f.Name()

			if strings.Contains(f.Name(), downMigr) {
				fbytes, _ = os.ReadFile(fmt.Sprintf("%s/%s", dirPath, f.Name()))
				break
			}
		}

		if len(fbytes) == 0 {
			log.Println("Down migration file is empty")
			continue
		}

		sql := string(fbytes)
		log.Printf("Rolling back: %s", downMigr)

		_, err = conn.Exec(context.Background(), sql)
		if err != nil {
			log.Fatal(err)
			return err
		}

		// Delete "up" record from migrations table
		_, err = conn.Exec(
			context.Background(),
			`delete from "migrations" where "migration" = $1`,
			upMigrName)
		if err != nil {
			log.Fatal(err)
			return err
		}

		rolledback++
	}

	return nil
}
