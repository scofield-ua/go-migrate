package migrate

import (
	"os"
	"strings"
	"testing"

	"github.com/scofield-ua/go-migrate/test"
)

func TestCreateMigration(t *testing.T) {
	sr := test.Setup(test.SetupParams{
		T: t,
	})

	defer sr.DeferFunc()

	var err error

	err = CreateMigration("test_migration", sr.MigrationsDir)
	if err != nil {
		t.Fatal(err)
	}

	items, err := os.ReadDir(sr.MigrationsDir)
	if err != nil {
		t.Fatal(err)
	}

	if len(items) == 0 {
		t.Error("Migrations folder is empty")
	}

	var upMigrExist bool
	var downMigrExist bool
	for _, item := range items {
		if strings.Contains(item.Name(), ".up.") {
			upMigrExist = true
		}

		if strings.Contains(item.Name(), ".down.") {
			downMigrExist = true
		}
	}

	if !upMigrExist {
		t.Fatal("UP migration file does not exist")
	}

	if !downMigrExist {
		t.Fatal("DOWN migration file does not exist")
	}
}

func BenchmarkCreateMigration(b *testing.B) {
	tmpDir, err := os.MkdirTemp("", "migrations")
	if err != nil {
		b.Fatal(err)
	}

	for b.Loop() {
		CreateMigration("abc", tmpDir)
	}
}
