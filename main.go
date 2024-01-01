package main

import (
	"context"
	"github.com/Vinayaks439/golang-backend/api"
	db2 "github.com/Vinayaks439/golang-backend/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

const (
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

func main() {
	conn, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store := db2.NewStore(conn)
	err = api.NewServer(store).Start("localhost:8080")
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
