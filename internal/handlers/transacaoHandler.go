package handlers

import (
	"context"
	"fmt"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vitorlfaria/rinha2024q1Go/internal/db"
	"github.com/vitorlfaria/rinha2024q1Go/internal/models"
)

type idUri struct {
	Id int `uri:"id" binding:"required"`
}

func HandleTransacao(c *gin.Context) {
	conn, err := db.Db.Begin(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": "Erro ao iniciar transação no banco de dados."})
		fmt.Println("Erro ao iniciar transação no banco de dados: ", err)
		return
	}
	defer conn.Commit(context.Background())

	var transacao models.TransacaoRequest
	var id idUri
	if err := c.ShouldBindJSON(&transacao); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ValidarTransacao(&transacao, c)

	var cliente models.Cliente
	conn.QueryRow(context.Background(), `SELECT * FROM "Clientes" WHERE "Id" = $1`, id.Id).Scan(&cliente.Id, &cliente.Limite, &cliente.Saldo)

	if cliente.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cliente não encontrado."})
		return
	}

	var saldo float64
	if transacao.Tipo == "d" {
		saldo = float64(cliente.Saldo - transacao.Valor)
		if math.Abs(saldo) > float64(cliente.Limite) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Saldo insuficiente."})
			return
		}
	} else {
		saldo = float64(cliente.Saldo + transacao.Valor)
	}

	_, err = conn.Exec(context.Background(), `INSERT INTO "Transacoes" ("Valor", "Tipo", "Descricao", "Realizada_Em", "ClienteId") VALUES ($1, $2, $3, NOW(), $4)`,
		transacao.Valor, transacao.Tipo, transacao.Descricao, cliente.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao inserir transação."})
		return
	}

	_, err = conn.Exec(context.Background(), `UPDATE "Clientes" SET "Saldo" = $1 WHERE "Id" = $2`, int(saldo), cliente.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar saldo do cliente."})
		return
	}

	transacaoResponse := models.TransacaoResponse{Limite: cliente.Limite, Saldo: int(saldo)}

	c.JSON(http.StatusOK, transacaoResponse)
}

func ValidarTransacao(transacao *models.TransacaoRequest, c *gin.Context) {
	if transacao.Valor <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Valor da transação deve ser maior que zero."})
		return
	}
	if transacao.Tipo != "d" && transacao.Tipo != "c" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tipo da transação inválido."})
		return
	}
	if len(transacao.Descricao) < 1 || len(transacao.Descricao) > 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Descrição da transação deve ter entre 1 e 10 caracteres."})
		return
	}
}
