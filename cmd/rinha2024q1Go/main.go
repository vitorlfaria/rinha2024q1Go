package main

import (
	"github.com/vitorlfaria/rinha2024q1Go/internal/db"
	"github.com/vitorlfaria/rinha2024q1Go/internal/routes"
)

func main() {
	db.ConnectDB()
	defer db.Db.Close()
	r := routes.SetupRouter()
	r.Run()
}
