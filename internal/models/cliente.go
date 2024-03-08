package models

type Cliente struct {
	Id     int `json:"id"`
	Limite int `json:"limite"`
	Saldo  int `json:"saldo"`
}
