package entity

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Tokens   int    `json:"tokens"`
}
