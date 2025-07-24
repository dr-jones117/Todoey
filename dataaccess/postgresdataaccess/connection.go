package postgresdataaccess

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type PostgresTodoDataAccess struct {
	db *sql.DB
}

func (da *PostgresTodoDataAccess) ConnectDataAccess() error {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	fmt.Println(dsn)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("sql.Open failed: %w", err)
	}

	da.db = db

	log.Println("✅ Connected to Postgres")
	return nil
}

func (da *PostgresTodoDataAccess) DisconnectDataAccess() error {
	if err := da.db.Close(); err != nil {
		return fmt.Errorf("error")
	}

	log.Println("Successfully disconnected from the PostgreSQL server...")
	return nil
}
