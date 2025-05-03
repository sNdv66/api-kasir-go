package models

type Transaction struct {
	ID            string  `json:"id,omitempty"`
	BranchID      string  `json:"branch_id"`
	UserID        string  `json:"user_id"`
	Total         float64 `json:"total"`
	PaymentMethod string  `json:"payment_method"`
	PaidAmount    float64 `json:"paid_amount"`
	Change        float64 `json:"change"`
	CreatedAt     string  `json:"created_at,omitempty"`
}