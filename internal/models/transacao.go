package models

import "time"

type Transacao struct {
	Valor       int       `json:"valor"`
	Tipo        string    `json:"tipo"`
	Descricao   string    `json:"descricao"`
	RealizadaEm time.Time `json:"realizadaEm"`
}

type TransacaoResponse struct {
	Limite int `json:"limite"`
	Saldo  int `json:"saldo"`
}

type TransacaoRequest struct {
	Valor     int    `json:"valor"`
	Tipo      string `json:"tipo"`
	Descricao string `json:"descricao"`
}
