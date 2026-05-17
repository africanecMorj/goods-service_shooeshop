package config

import (
	"github.com/jackc/pgx/v5/pgxpool"
	
	"os"
	"context"
	"log"
)

func PostgresInit () *pgxpool.Pool{
	db, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	log.Println(err)
	if err != nil {
		log.Fatal("Postgres init error:", err)
	}

	if err := db.Ping(context.Background()); err != nil {
		log.Fatal("Postgres ping error:", err)
	}

	log.Println("Postgres connected")


	return db
}

