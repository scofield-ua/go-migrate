package tools

import (
	"strings"
)

type MigrationVariant string

const (
	MigrationUp   MigrationVariant = "up"
	MigrationDown MigrationVariant = "down"
)

func (mv MigrationVariant) String() string {
	return string(mv)
}

// Change migration variant and get partial of it
// Example: 12345_create_users_table.down.sql -> create_users_table.up.sql
func ChangeMigrationVariant(current string, reqVariant MigrationVariant) string {
	_, after, _ := strings.Cut(current, "_")

	var currentVariant MigrationVariant
	if strings.Contains(current, ".up.") {
		currentVariant = MigrationUp
	} else {
		currentVariant = MigrationDown
	}

	return strings.Replace(after, "."+currentVariant.String()+".", "."+reqVariant.String()+".", 1)
}
