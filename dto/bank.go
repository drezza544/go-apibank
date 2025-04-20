package dto

type BankRequest struct {
	Code            *string `json:"code"`
	Name            *string `json:"name"`
	Address         *string `json:"address"`
	CostTransaction *int64  `json:"cost_transaction"`
}
