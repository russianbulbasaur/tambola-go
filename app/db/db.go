package db

import (
	"database/sql"
	"fmt"
	"os"
)

const poolCount = 10

func InitDB() *sql.DB {
	connectionString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		8000,
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}
	db.SetMaxIdleConns(poolCount)
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}
