package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/avran02/kode/config"
)

func MustGetPostgresConnection(confing config.DB) *sql.DB {
	dsn := getDsn(confing)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err.Error())
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %s", err.Error())
	}

	return db
}

func getDsn(config config.DB) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.Database,
	)
}
