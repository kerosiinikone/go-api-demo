package main

import (
	"log"
	"path/filepath"

	dbx "github.com/go-ozzo/ozzo-dbx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/kerosiinikone/go-server/internal/config"
	_ "github.com/lib/pq"
)

func CreateDBConnection() *dbx.DB {
	currentDir, err := filepath.Abs(filepath.Dir("."))
	
	if err != nil {
		log.Fatal(err.Error())
	}

	targetDir := filepath.Join(currentDir)
	configFilePath := filepath.Join(targetDir, "local.yml")

	cfg, err := config.Load(configFilePath)

	if err != nil {
		log.Fatal(err.Error())
	}

	db, err := dbx.Open("postgres", cfg.DBU)

	if err != nil {
		log.Fatal(err.Error())
	}

	return db
}

