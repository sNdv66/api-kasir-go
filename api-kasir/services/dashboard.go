package services

import (
	"api-kasir/utils"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"time"
)

type TodaySalesResult struct {
	TotalSales       float64 `json:"total_sales"`
	TotalTransactions int     `json:"total_transactions"`
}

func FetchTodaySales(branchID string) (*TodaySalesResult, error) {
	today := time.Now().Format("2006-01-02")
	query := fmt.Sprintf("branch_id=eq.%s&created_at=gte.%sT00:00:00&created_at=lt.%sT23:59:59", branchID, today, today)

	req, err := utils.NewRequest("GET", "transactions?select=total,created_at&"+query, nil)
	if err != nil {
		return nil, err
	}

	resp, err := utils.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var transactions []struct {
		Total     float64   `json:"total"`
		CreatedAt time.Time `json:"created_at"`
	}

	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &transactions); err != nil {
		return nil, err
	}

	var totalSales float64
	for _, tx := range transactions {
		totalSales += tx.Total
	}

	return &TodaySalesResult{
		TotalSales:       totalSales,
		TotalTransactions: len(transactions),
	}, nil
}

type TopProduct struct {
	ProductID    string `json:"product_id"`
	Name         string `json:"name"`
	QuantitySold int    `json:"quantity_sold"`
}

func FetchTopProductsToday(branchID string) ([]TopProduct, error) {
	today := time.Now().Format("2006-01-02")

	// Ambil transaksi hari ini untuk cabang
	query := fmt.Sprintf(`
		transaction_items?select=product_id,quantity,transaction_id,transactions(branch_id,created_at),products(name)&
		transactions.branch_id=eq.%s&
		transactions.created_at=gte.%sT00:00:00&transactions.created_at=lt.%sT23:59:59
	`, branchID, today, today)

	req, err := utils.NewRequest("GET", query, nil)
	if err != nil {
		return nil, err
	}

	resp, err := utils.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var items []struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
		Products  struct {
			Name string `json:"name"`
		} `json:"products"`
	}

	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, err
	}

	// Rekap
	counter := map[string]TopProduct{}
	for _, item := range items {
		if p, ok := counter[item.ProductID]; ok {
			p.QuantitySold += item.Quantity
			counter[item.ProductID] = p
		} else {
			counter[item.ProductID] = TopProduct{
				ProductID:    item.ProductID,
				Name:         item.Products.Name,
				QuantitySold: item.Quantity,
			}
		}
	}

	// Ubah ke slice & sort
	var result []TopProduct
	for _, p := range counter {
		result = append(result, p)
	}

	// Optional: sort by QuantitySold desc
	sort.Slice(result, func(i, j int) bool {
		return result[i].QuantitySold > result[j].QuantitySold
	})

	// Batasi misalnya 5 produk
	if len(result) > 5 {
		result = result[:5]
	}

	return result, nil
}

type DailySales struct {
	Date       string  `json:"date"`
	TotalSales float64 `json:"total_sales"`
}

func FetchWeeklySales(branchID string) ([]DailySales, error) {
	end := time.Now()
	start := end.AddDate(0, 0, -6) // 7 hari terakhir
	startStr := start.Format("2006-01-02")
	endStr := end.Format("2006-01-02")

	query := fmt.Sprintf(`
		transactions?select=created_at,total,branch_id&branch_id=eq.%s&created_at=gte.%sT00:00:00&created_at=lte.%sT23:59:59
	`, branchID, startStr, endStr)

	req, err := utils.NewRequest("GET", query, nil)
	if err != nil {
		return nil, err
	}
	resp, err := utils.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var transactions []struct {
		CreatedAt time.Time `json:"created_at"`
		Total     float64   `json:"total"`
	}

	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &transactions); err != nil {
		return nil, err
	}

	// Inisialisasi map tanggal
	salesMap := make(map[string]float64)
	for i := 0; i < 7; i++ {
		date := start.AddDate(0, 0, i).Format("2006-01-02")
		salesMap[date] = 0
	}

	for _, t := range transactions {
		date := t.CreatedAt.Format("2006-01-02")
		salesMap[date] += t.Total
	}

	var result []DailySales
	for i := 0; i < 7; i++ {
		date := start.AddDate(0, 0, i).Format("2006-01-02")
		result = append(result, DailySales{
			Date:       date,
			TotalSales: salesMap[date],
		})
	}

	return result, nil
}

// TransactionDaily defines the daily transaction count structure
type TransactionDaily struct {
	Date              string `json:"date"`
	TotalTransactions int    `json:"total_transactions"`
}

// FetchTransactionsDaily retrieves daily transaction counts for a branch
func FetchTransactionsDaily(branchID string) ([]TransactionDaily, error) {
	url := fmt.Sprintf("transactions?select=created_at&branch_id=eq.%s", branchID)
	req, err := utils.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	resp, err := utils.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var raw []struct {
		CreatedAt string `json:"created_at"`
	}
	
	err = json.NewDecoder(resp.Body).Decode(&raw)
	if err != nil {
		return nil, err
	}
	
	countMap := make(map[string]int)
	for _, r := range raw {
		date := r.CreatedAt[:10] // ambil yyyy-mm-dd
		countMap[date]++
	}
	
	result := make([]TransactionDaily, 0)
	for date, count := range countMap {
		result = append(result, TransactionDaily{
			Date: date,
			TotalTransactions: count,
		})
	}
	
	return result, nil
}

// FetchAverageTransactionValue calculates the average transaction value for a branch
func FetchAverageTransactionValue(branchID string) (float64, error) {
	url := fmt.Sprintf("transactions?branch_id=eq.%s&select=total", branchID)
	req, err := utils.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	
	resp, err := utils.Client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	
	var raw []struct {
		Total float64 `json:"total"`
	}
	
	err = json.NewDecoder(resp.Body).Decode(&raw)
	if err != nil {
		return 0, err
	}
	
	var sum float64
	for _, r := range raw {
		sum += r.Total
	}
	
	if len(raw) == 0 {
		return 0, nil
	}
	
	return sum / float64(len(raw)), nil
}

// FetchLowStockAlert retrieves products with low stock for a branch
func FetchLowStockAlert(branchID string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("products?branch_id=eq.%s&stock=lt.10", branchID)
	req, err := utils.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	resp, err := utils.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var result []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	return result, err
}


func FetchTopProducts(branchID string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("top_products?branch_id=eq.%s&limit=5&order=total_quantity.desc", branchID)

	req, err := utils.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := utils.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}


func FetchSalesChart(branchID string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("sales_chart?branch_id=eq.%s&order=date.asc", branchID)

	req, err := utils.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := utils.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}


func FetchLowStock(branchID string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("products?branch_id=eq.%s&stock=lt.10&order=stock.asc", branchID)

	req, err := utils.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := utils.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}


