package main

import (
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	if err := Cmd().Execute(); err != nil {
		os.Exit(1)
	}
}
