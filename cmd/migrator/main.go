package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var postgresUrl, migrationsPath string

func main() {
	flag.StringVar(&postgresUrl, "pg-url", "", "postgres url")
	flag.StringVar(&migrationsPath, "mg-path", "", "path to migrations folder")
	flag.Parse()

	if postgresUrl == "" {
		panic("postgres url is required")
	}

	if migrationsPath == "" {
		panic("migrations path is required")
	}

	m, err := migrate.New("file://"+migrationsPath, postgresUrl)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("Nothing to migrate")
			return
		}

		panic(err)
	}

	fmt.Println("Migrations migrated")
}
