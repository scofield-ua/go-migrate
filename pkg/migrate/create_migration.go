package migrate

import (
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Create migration file
func CreateMigration(name string, path string) error {
	var err error

	path, err = filepath.Abs(path)
	if err != nil {
		return err
	}
	if err = os.Mkdir(path, 0755); err != nil && !os.IsExist(err) {
		return err
	}

	for i, ext := range []string{"up", "down"} {
		unixSeconds := time.Now().Unix()

		fileName := strconv.FormatInt(unixSeconds, 10) + strconv.Itoa(i) + "_" + name + "." + ext + ".sql"
		// fileName := fmt.Sprintf("%d%d_%s.%s.sql", unixSeconds, i, name, ext)
		fullPath := filepath.Join(path, fileName)

		err = os.WriteFile(fullPath, []byte("-- Replace with SQL"), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
