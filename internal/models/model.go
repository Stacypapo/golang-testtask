package models

type Wallet struct {
	Address string  `json:"address"`
	Balance float64 `json:"balance"`
}

type Transaction struct {
	ID     int     `json:"id"`
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

type CreateTransactionRequest struct {
	From   string  `json:"from" example:"e240d825d255af751f5f55af8d9671be"`
	To     string  `json:"to" example:"abdf2236c0a3b4e2639b3e182d994c88e"`
	Amount float64 `json:"amount" example:"10" minimum:"0.01"`
}

type StatusResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Transaction completed"`
}
