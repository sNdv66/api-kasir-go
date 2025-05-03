package services

import (
	"bytes"
	"encoding/json"
	"io"
	"api-kasir/utils"
	"net/http"
	"fmt"
)

type StockMovement struct {
	ProductID string `json:"product_id"`
	BranchID  string `json:"branch_id"`
	Type      string `json:"type"`     // "in" atau "out"
	Quantity  int    `json:"quantity"` // wajib
	Note      string `json:"note"`
}

// Kirim data ke Supabase
func AddStockMovement(data StockMovement) (*http.Response, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var body io.Reader = bytes.NewBuffer(jsonData)

	req, err := utils.NewRequest("POST", "stock_movements", body)
	if err != nil {
		return nil, err
	}

	return utils.Client.Do(req)
}


func FetchSalesReport(branchID, startDate, endDate string) (map[string]interface{}, error) {
	query := fmt.Sprintf(`transactions?branch_id=eq.%s&created_at=gte.%s&created_at=lte.%s`, branchID, startDate, endDate)
	req, err := utils.NewRequest("GET", query, nil)
	if err != nil {
		return nil, err
	}

	resp, err := utils.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var transactions []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&transactions); err != nil {
		return nil, err
	}

	var totalAmount float64
	for _, t := range transactions {
		if val, ok := t["total_amount"].(float64); ok {
			totalAmount += val
		}
	}

	return map[string]interface{}{
		"total_transactions": len(transactions),
		"total_sales":        totalAmount,
	}, nil
}


func UpdateProductStock(productID string, newStock int) error {
	updateData := map[string]interface{}{
		"stock": newStock,
	}

	jsonData, err := json.Marshal(updateData)
	if err != nil {
		return err
	}

	req, err := utils.NewRequest("PATCH", "products?id=eq."+productID, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := utils.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("failed to update stock: status code %d", resp.StatusCode)
	}

	return nil
}