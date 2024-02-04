package entities

type TransactionRequest struct {
	Payload Transaction `json:"payload"`
}

type Transaction struct {
	Header Header  `json:"header"`
	ID     string  `json:"id"`
	Origin string  `json:"origin"`
	UserID string  `json:"user_id"`
	Amount float64 `json:"amount"`
	Kind   string  `json:"kind"`
}
