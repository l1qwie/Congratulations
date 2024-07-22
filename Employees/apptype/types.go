package apptype

var SymKey []byte

type Employee struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Birthday string `json:"birthday"`
}
