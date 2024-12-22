package main

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	log "github.com/sirupsen/logrus"
)

//nolint: gochecknoglobals
var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	if buildVersion == "" {
		buildVersion = "N/A"
	}

	if buildDate == "" {
		buildDate = "N/A"
	}

	if buildCommit == "" {
		buildCommit = "N/A"
	}

	log.Printf("Build version: %s\n", buildVersion)
	log.Printf("Build date: %s\n", buildDate)
	log.Printf("Build commit: %s\n", buildCommit)

	if err := Cmd().Execute(); err != nil {
		panic(err)
	}
}
