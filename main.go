package main

import (
	"context"
	"github.com/Vinayaks439/golang-backend/api"
	db2 "github.com/Vinayaks439/golang-backend/db/sqlc"
	"github.com/Vinayaks439/golang-backend/util"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn, err := pgxpool.New(context.Background(), config.DbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store := db2.NewStore(conn)
	err = api.NewServer(store).Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
