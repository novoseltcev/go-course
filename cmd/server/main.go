package main

import _ "github.com/jackc/pgx/v5/stdlib"

func main() {
	if err := Cmd().Execute(); err != nil {
		panic(err)
	}
}
