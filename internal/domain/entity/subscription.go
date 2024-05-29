package entity

type Subscription struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Plan   string `json:"plan"`
	Tokens int    `json:"tokens"`
}
