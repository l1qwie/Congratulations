package apptype

// AES-256
var SymKey []byte

type Employee struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Birthday string `json:"birthday"`
}

type KafkaEmployee struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Birthday string `json:"birthday"`
	WhatDo   string `json:"whatdo"`
	SecondId int    `json:"secondid"`
}
