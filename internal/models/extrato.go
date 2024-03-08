package models

import "time"

type Extrato struct {
	Saldo             Saldo       `json:"saldo"`
	UltimasTransacoes []Transacao `json:"ultimasTransacoes"`
}

type Saldo struct {
	Total       int       `json:"total"`
	DataExtrato time.Time `json:"dataExtrato"`
	Limite      int       `json:"limite"`
}
