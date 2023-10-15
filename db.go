package kosample

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func GetConnect() *sqlx.DB {
	dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		GetEnv("POSTGRES_HOST", "localhost"),
		GetEnv("POSTGRES_PORT", "5432"),
		GetEnv("POSTGRES_USER", "postgres"),
		GetEnv("POSTGRES_PASSWORD", "password"),
		GetEnv("POSTGRES_DB", "kosample"))
	db, err := sqlx.Connect("postgres", dbinfo)
	if err != nil {
		panic(err)
	}

	return db
}
