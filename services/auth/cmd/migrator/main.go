package migrator

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
)

func main() {
	var storagePath, migrationsPath string

	flag.StringVar(&storagePath, "storage", "", "storage path")
	flag.StringVar(&migrationsPath, "migrations", "", "migrations path")
	flag.Parse()

	if storagePath == "" {
		panic("storage path is empty")
	}

	if migrationsPath == "" {
		panic("migrations path is empty")
	}

	m, err := migrate.New(migrationsPath, storagePath)

	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("ok; no migrations")
			return
		}

		panic(err)
	}

	fmt.Println("migrations completed")
}
