package main

import (
	"auth_wallet/api"
	db "auth_wallet/db/sqlc"
	"auth_wallet/token"
	"auth_wallet/util"
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(connPool)
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		log.Fatal("cannot create token maker:", err)
	}

	server, err := api.NewServer(config, store, tokenMaker)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
