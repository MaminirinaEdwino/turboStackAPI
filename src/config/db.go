package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "user=postgres password=root dbname=turbostack_db sslmode=disable host=localhost port=5432"
	}
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}

	// Initialisation des tables de base si elles n'existent pas
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		role VARCHAR(20) NOT NULL DEFAULT 'user' CHECK (role IN ('admin', 'user'))
	);
	CREATE TABLE IF NOT EXISTS projects (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		nom VARCHAR(255) UNIQUE NOT NULL,
		type_project VARCHAR(100),
		description TEXT,
		data JSONB,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err = DB.Exec(schema)
	if err != nil {
		log.Fatalf("Échec de la création des tables initiales : %v", err)
	}

	fmt.Println("Successfully connected to PostgreSQL database")
}
