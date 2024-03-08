package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vitorlfaria/rinha2024q1Go/internal/db"
	"github.com/vitorlfaria/rinha2024q1Go/internal/models"
)

func HandleExtrato(c *gin.Context) {
	conn, err := db.Db.Begin(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao iniciar transação no banco de dados."})
		return
	}
	defer conn.Commit(context.Background())

	var id idUri
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var cliente models.Cliente
	conn.QueryRow(context.Background(), `SELECT * FROM "Clientes" WHERE "Id" = $1`, id.Id).Scan(&cliente.Id, &cliente.Limite, &cliente.Saldo)

	if cliente.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cliente não encontrado."})
		return
	}

	rows, err := conn.Query(
		context.Background(),
		`SELECT "Valor", "Tipo", "Descricao", "Realizada_Em"
		FROM "Transacoes"
		WHERE "ClienteId" = $1
		ORDER BY "Realizada_Em"
		LIMIT 10`, cliente.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar transações."})
		return
	}

	var transacoes []models.Transacao
	for rows.Next() {
		var transacao models.Transacao
		if err := rows.Scan(&transacao.Valor, &transacao.Tipo, &transacao.Descricao, &transacao.RealizadaEm); err != nil {
			fmt.Println("Erro ao buscar transações: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar transação."})
			return
		}

		transacoes = append(transacoes, transacao)
	}

	extrato := models.Extrato{
		Saldo: models.Saldo{
			Total:       cliente.Saldo,
			DataExtrato: time.Now(),
			Limite:      cliente.Limite,
		},
		UltimasTransacoes: transacoes,
	}

	c.JSON(http.StatusOK, extrato)
}
