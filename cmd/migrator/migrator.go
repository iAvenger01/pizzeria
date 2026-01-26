package main

import (
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"os"
)

func main() {
	m, err := migrate.New(
		"file://migrations",
		"pgx5://pizzeria:pizzeria14256@192.168.0.127:5432/pizzeria?sslmode=enable",
	)

	if err != nil {
		panic(err)
	}

	steps := flag.Int("step", 0, "Сколько миграций обработать?")
	flag.Parse()

	if os.Args[1] == "up" {
		if *steps == 0 {
			if err := m.Up(); err != nil {
				fmt.Println("Migration failed:", err)
				return
			}
		} else if *steps > 0 {
			if err := m.Steps(*steps); err != nil {
				fmt.Println("Migration failed:", err)
				return
			}
		}
		fmt.Println("Migrations applied successfully!")
	} else if os.Args[1] == "down" {
		if *steps == 0 {
			if err := m.Down(); err != nil {
				fmt.Println("Migration failed:", err)
				return
			}
		} else if *steps < 0 {
			if err := m.Steps(*steps); err != nil {
				fmt.Println("Migration failed:", err)
				return
			}
		}
		fmt.Println("Migrations rollback successfully!")
	}

	return
}
