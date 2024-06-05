package entity

type Subscription struct {

	ID       int64  `json:"id"`
	UserID   int64  `json:"user_id"`
	Tokens   int64  `json:"tokens"`
	Price    int64  `json:"price"`
	Currency string `json:"currency"`

}
