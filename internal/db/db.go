package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

var Db *pgxpool.Pool

func ConnectDB() {
	connStr := "postgres://postgres:postgres@db:5432/rinha"
	db, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		fmt.Println("Erro ao conectar no banco de dados: ", err)
	}
	Db = db
}
