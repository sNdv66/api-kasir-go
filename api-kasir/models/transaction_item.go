package models

type TransactionItem struct {
	ID            string  `json:"id,omitempty"`
	TransactionID string  `json:"transaction_id"`
	ProductID     string  `json:"product_id"`
	Quantity      int     `json:"quantity"`
	Subtotal      float64 `json:"subtotal"`
}