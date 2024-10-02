package main

import (
	"os"
	migration "test-no-migrate/migrate/script"
)

func main() {
	migration.DBMigration(os.Args[1:]...)
}
