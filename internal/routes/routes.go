package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vitorlfaria/rinha2024q1Go/internal/handlers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/clientes/:id/transacoes", handlers.HandleTransacao)
	r.GET("/clientes/:id/extrato", handlers.HandleExtrato)
	return r
}
