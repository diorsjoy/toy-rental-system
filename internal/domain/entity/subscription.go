package entity

type Subscription struct {
	ID     int64  `json:"id"`
	UserID int64  `json:"user_id"`
	Plan   string `json:"plan"`
	Tokens int64  `json:"tokens"`
}
