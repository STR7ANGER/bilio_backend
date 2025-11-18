package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/joho/godotenv"

	"github.com/nava1525/bilio-backend/internal/config"
	"github.com/nava1525/bilio-backend/internal/database"
)

func main() {
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env")
	_ = godotenv.Load("../../.env")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	dbClient, err := database.NewClient(cfg.Database.URL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbClient.Disconnect()

	// Get the project root
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get working directory: %v", err)
	}
	
	// Find migrations directory - check if we're in bilio_backend or cmd/server
	var migrationsDir string
	if filepath.Base(wd) == "bilio_backend" {
		migrationsDir = filepath.Join(wd, "migrations")
	} else if filepath.Base(wd) == "server" || filepath.Base(wd) == "migrate" {
		// If we're in cmd/server or cmd/migrate, go up to bilio_backend
		migrationsDir = filepath.Join(wd, "..", "..", "migrations")
	} else {
		// Try relative to current directory
		migrationsDir = filepath.Join(wd, "migrations")
	}
	
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		// Try absolute path from bilio_backend
		migrationsDir = filepath.Join(filepath.Dir(filepath.Dir(wd)), "bilio_backend", "migrations")
		if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
			log.Fatalf("migrations directory not found. Tried: %s and %s", filepath.Join(wd, "migrations"), migrationsDir)
		}
	}
	
	fmt.Printf("Using migrations directory: %s\n", migrationsDir)

	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		log.Fatalf("failed to read migrations directory: %v", err)
	}

	var migrations []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			migrations = append(migrations, file.Name())
		}
	}

	sort.Strings(migrations)

	db := dbClient.DB()
	for _, migration := range migrations {
		migrationPath := filepath.Join(migrationsDir, migration)
		fmt.Printf("Applying migration: %s\n", migration)

		sql, err := ioutil.ReadFile(migrationPath)
		if err != nil {
			log.Fatalf("failed to read migration %s: %v", migration, err)
		}

		if _, err := db.Exec(string(sql)); err != nil {
			log.Fatalf("failed to apply migration %s: %v", migration, err)
		}

		fmt.Printf("✓ Applied migration: %s\n", migration)
	}

	fmt.Println("All migrations applied successfully!")
}

