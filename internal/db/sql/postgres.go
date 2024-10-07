package sql

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *sql.DB

func InitPostgres(dsn string) (*sql.DB, error) {
	fmt.Println("Connecting to Postgres", dsn)
	dbConnection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	database, err := dbConnection.DB()
	if err != nil {
		return nil, err
	}
	database.SetConnMaxIdleTime(20)
	database.SetMaxOpenConns(200)
	DB = database
	return DB, nil
}

func ClosePostgres() {
	if DB != nil {
		DB.Close()
	}
}
