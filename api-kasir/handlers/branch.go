package handlers

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"api-kasir/models"
	"api-kasir/supabase"
	"api-kasir/utils"
	"api-kasir/services"
	"io"
	"bytes"
	"fmt"
	"net/http"
	"time"
	"strings"
)

type PaymentRequest struct {
	Items []struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	} `json:"items"`
	PaymentMethod string  `json:"payment_method"`
	PaidAmount    float64 `json:"paid_amount"`
}

type PendingOrder struct {
	ID        string          `json:"id"`
	BranchID  string          `json:"branch_id"`
	UserID    string          `json:"user_id"`
	Orders    json.RawMessage `json:"orders"` // JSON array
	Note      string          `json:"note,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
}



func HandlePayment(c *fiber.Ctx) error {
	var req PaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	total := 0.0
	transactionItems := []map[string]interface{}{}

	for _, item := range req.Items {
		// Ambil produk dari Supabase
		path := fmt.Sprintf("products?id=eq.%s", item.ProductID)
		httpReq, err := utils.NewRequest(http.MethodGet, path, nil)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "Gagal request produk"})
		}
		res, err := utils.Client.Do(httpReq)
		if err != nil || res.StatusCode != 200 {
			return c.Status(400).JSON(fiber.Map{"message": "Produk tidak ditemukan"})
		}
		defer res.Body.Close()

		var products []struct {
			ID    string  `json:"id"`
			Price float64 `json:"price"`
			Stock int     `json:"stock"`
		}
		body, _ := io.ReadAll(res.Body)
		_ = json.Unmarshal(body, &products)

		if len(products) == 0 {
			return c.Status(400).JSON(fiber.Map{"message": "Produk tidak valid"})
		}

		product := products[0]
		if product.Stock < item.Quantity {
			return c.Status(400).JSON(fiber.Map{
				"message": "Stok tidak cukup untuk produk " + product.ID,
			})
		}

		subtotal := product.Price * float64(item.Quantity)
		total += subtotal

		transactionItems = append(transactionItems, map[string]interface{}{
			"product_id": item.ProductID,
			"quantity":   item.Quantity,
			"subtotal":   subtotal,
		})
	}

	if req.PaidAmount < total {
		return c.Status(400).JSON(fiber.Map{"message": "Jumlah bayar kurang"})
	}

	// Ambil data dari JWT middleware
    user := c.Locals("user").(*utils.Claims)
    userID := user.UserID
    branchID := user.BranchID
	// Simpan transaksi
	trxPayload := map[string]interface{}{
		"branch_id":      branchID,
		"user_id":        userID,
		"total":          total,
		"payment_method": req.PaymentMethod,
		"paid_amount":    req.PaidAmount,
		"change":         req.PaidAmount - total,
	}

	trxJSON,err:= json.Marshal(trxPayload)
	/*trxReq, err := utils.NewRequest(http.MethodPost, "transactions", bytes.NewBuffer(trxJSON))
	*/
	
	trxReq, err := utils.NewRequest(http.MethodPost, "transactions", bytes.NewBuffer(trxJSON))
     trxReq.Header.Set("Prefer", "return=representation")
	
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal membuat request transaksi"})
	}
	trxRes, err := utils.Client.Do(trxReq)
	if err != nil || trxRes.StatusCode >= 400 {
		body, _ := io.ReadAll(trxRes.Body)
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal menyimpan transaksi",
			"error":   string(body),
		})
	}
	defer trxRes.Body.Close()
   body, err := io.ReadAll(trxRes.Body)
   if err != nil {
    return c.Status(500).JSON(fiber.Map{
        
        "message": "Gagal membaca response dari Supabase",
    })
    }

	var inserted []struct {
		ID string `json:"id"`
	}
	// body, _ = io.ReadAll(trxRes.Body)
   _ = json.Unmarshal(body, &inserted)
   
   if err := json.Unmarshal(body, &inserted); err != nil {
	return c.Status(500).JSON(fiber.Map{
		"message": "Gagal mengurai response dari Supabase",
		"error":   err.Error(),
	})
}

   if len(inserted) == 0 {
	return c.Status(500).JSON(fiber.Map{
		"message": "Transaksi berhasil disimpan, tapi tidak mendapatkan ID transaksi dari Supabase",
	})
    }

    transactionID := inserted[0].ID
	// Tambahkan transaction_id ke setiap item
	for i := range transactionItems {
		transactionItems[i]["transaction_id"] = transactionID
	}

	// Simpan ke transaction_items
	itemsJSON, _ := json.Marshal(transactionItems)
	itemsReq, err := utils.NewRequest(http.MethodPost, "transaction_items", bytes.NewBuffer(itemsJSON))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal membuat request item transaksi"})
	}
	itemsRes, err := utils.Client.Do(itemsReq)
	if err != nil || itemsRes.StatusCode >= 400 {
		body, _ := io.ReadAll(itemsRes.Body)
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal menyimpan item transaksi",
			"error":   string(body),
		})
	}
	defer itemsRes.Body.Close()

	return c.JSON(fiber.Map{
		"message":        "Pembayaran berhasil",
		"transaction_id": transactionID,
	})
}


// GetBranches
func GetBranches(c *fiber.Ctx) error {
	branches, err := supabase.GetBranches()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(branches)
}

func GetBranchByID(c *fiber.Ctx) error {
	id := c.Params("id")

	req, err := utils.NewRequest("GET", "branches?id=eq."+id, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create request",
		})
	}

	res, err := utils.Client.Do(req)
	if err != nil || res.StatusCode != 200 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Branch not found",
		})
	}
	defer res.Body.Close()

	var branches []models.Branch
	if err := json.NewDecoder(res.Body).Decode(&branches); err != nil || len(branches) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Branch not found",
		})
	}

	return c.JSON(branches[0])
}



// CreateBranch
func CreateBranch(c *fiber.Ctx) error {
	var input models.CreateBranchInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	branch := models.Branch{
		ID:      uuid.New().String(),
		Name:    input.Name,
		Address: input.Address,
	}

	if err := supabase.CreateBranch(branch); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
    "message": "Cabang berhasil ditambahkan",
    "branch": branch,
         })
    
}

// UpdateBranch
func UpdateBranch(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is required",
		})
	}

	var input models.CreateBranchInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := supabase.UpdateBranch(id, input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Cabang berhasil diperbarui",
	})
}

// DeleteBranch
func DeleteBranch(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := supabase.DeleteBranch(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Cabang berhasil dihapus",
	})
}

// Login
func Login(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	email := user.Email
	password := user.Password
	
	if email == "" || password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email and password are required",
		})
	}

	req, err := utils.NewRequest("GET", "users?email=eq."+email+"&select=id,email,password,role,branch_id", nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create request",
		})
	}

	res, err := utils.Client.Do(req)
	if err != nil || res.StatusCode != 200 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}
	defer res.Body.Close()

	var users []models.User
	if err := json.NewDecoder(res.Body).Decode(&users); err != nil || len(users) == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	foundUser := users[0]

	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	token, err := utils.GenerateJWT(
		foundUser.ID,
		foundUser.Email,
		foundUser.Role,
		foundUser.BranchID,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return c.JSON(fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":        foundUser.ID,
			"role":      foundUser.Role,
			"branch_id": foundUser.BranchID,
		},
	})
}


func GetProductsByBranch(c *fiber.Ctx) error {
    branchID := c.Params("branch_id")
    query := "products?branch_id=eq." + branchID

    req, err := utils.NewRequest("GET", query, nil)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    res, err := utils.Client.Do(req)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    defer res.Body.Close()

    body, _ := io.ReadAll(res.Body)
    return c.Status(res.StatusCode).Send(body)
}




func CreateProduct(c *fiber.Ctx) error {
    body := c.Body()

    req, err := utils.NewRequest("POST", "products", bytes.NewReader(body))
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    res, err := utils.Client.Do(req)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    defer res.Body.Close()

    resBody, _ := io.ReadAll(res.Body)
    return c.Status(res.StatusCode).Send(resBody)
}




func UpdateProduct(c *fiber.Ctx) error {
    id := c.Params("id")
    path := "products?id=eq." + id

    req, err := utils.NewRequest("PATCH", path, bytes.NewReader(c.Body()))
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    res, err := utils.Client.Do(req)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    defer res.Body.Close()

    resBody, _ := io.ReadAll(res.Body)
    return c.Status(res.StatusCode).Send(resBody)
}




func DeleteProduct(c *fiber.Ctx) error {
    id := c.Params("id")
    path := "products?id=eq." + id

    req, err := utils.NewRequest("DELETE", path, nil)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    res, err := utils.Client.Do(req)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    defer res.Body.Close()

    resBody, _ := io.ReadAll(res.Body)
    return c.Status(res.StatusCode).Send(resBody)
}

// POST /transactions
func CreateTransaction(c *fiber.Ctx) error {
	var trx models.Transaction
	if err := c.BodyParser(&trx); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	data, err := json.Marshal(trx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to encode JSON"})
	}

	req, err := utils.NewRequest("POST", "transactions", bytes.NewReader(data))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := utils.Client.Do(req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	return c.Status(res.StatusCode).Send(body)
}

// GET /branches/:branch_id/transactions
func GetTransactionsByBranch(c *fiber.Ctx) error {
	branchID := c.Params("branch_id")
	req, err := utils.NewRequest("GET", "transactions?branch_id=eq."+branchID, nil)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := utils.Client.Do(req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	return c.Status(res.StatusCode).Send(body)
}

// GET /transactions/:id
func GetTransactionByID(c *fiber.Ctx) error {
	id := c.Params("id")
	req, err := utils.NewRequest("GET", "transactions?id=eq."+id, nil)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := utils.Client.Do(req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	return c.Status(res.StatusCode).Send(body)
}



// POST /transaction-items
func CreateTransactionItem(c *fiber.Ctx) error {
	var item models.TransactionItem
	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	data, err := json.Marshal(item)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to encode JSON"})
	}

	req, err := utils.NewRequest("POST", "transaction_items", bytes.NewReader(data))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := utils.Client.Do(req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	return c.Status(res.StatusCode).Send(body)
}

// GET /transactions/:id/items
func GetTransactionItemsByTransactionID(c *fiber.Ctx) error {
	transactionID := c.Params("id")
	req, err := utils.NewRequest("GET", "transaction_items?transaction_id=eq."+transactionID, nil)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := utils.Client.Do(req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	return c.Status(res.StatusCode).Send(body)
	
}

func CreateStockMovement(c *fiber.Ctx) error {
	var input services.StockMovement
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Ambil stok saat ini dari produk
	product, err := services.GetProductByID(input.ProductID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get product"})
	}

	// Hitung stok baru
	newStock := product.Stock
	if input.Type == "out" {
		if product.Stock < input.Quantity {
			return c.Status(400).JSON(fiber.Map{"error": "Stok tidak cukup"})
		}
		newStock -= input.Quantity
	} else if input.Type == "in" {
		newStock += input.Quantity
	}

	// Update stok produk
	if err := services.UpdateProductStock(input.ProductID, newStock); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update product stock"})
	}

	// Simpan movement ke Supabase
	resp, err := services.AddStockMovement(input)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer resp.Body.Close()

	return c.Status(201).JSON(fiber.Map{"message": "Stock movement recorded"})
}



func GetStockMovements(c *fiber.Ctx) error {
	branchID := c.Query("branch_id")
	productID := c.Query("product_id")
	filter := ""
	if branchID != "" {
		filter += "branch_id=eq." + branchID
	}
	if productID != "" {
		if filter != "" {
			filter += "&"
		}
		filter += "product_id=eq." + productID
	}

	if filter != "" {
		filter += "&"
	}
	filter += "order=created_at.desc"

	req, err := utils.NewRequest("GET", "stock_movements?"+filter, nil)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal membuat request"})
	}

	resp, err := utils.Client.Do(req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Request gagal"})
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return c.Status(resp.StatusCode).Send(body)
	}

	return c.Send(body)
}

func GetStockSummary(c *fiber.Ctx) error {
	branchID := c.Query("branch_id")
	from := c.Query("from")
	to := c.Query("to")

	if branchID == "" || from == "" || to == "" {
		return c.Status(400).JSON(fiber.Map{"error": "branch_id, from, and to are required"})
	}

	path := fmt.Sprintf(
		"stock_movements?branch_id=eq.%s&created_at=gte.%s&created_at=lte.%s",
		branchID, from, to,
	)

	req, err := utils.NewRequest("GET", path, nil)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	resp, err := utils.Client.Do(req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer resp.Body.Close()

	var movements []services.StockMovement
	if err := json.NewDecoder(resp.Body).Decode(&movements); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to parse response"})
	}

	// Hitung ringkasan
	totalIn := 0
	totalOut := 0
	byProduct := make(map[string]map[string]interface{})

	for _, m := range movements {
		if _, ok := byProduct[m.ProductID]; !ok {
			byProduct[m.ProductID] = map[string]interface{}{
				"product_id": m.ProductID,
				"in":         0,
				"out":        0,
			}
		}
		if m.Type == "in" {
			totalIn += m.Quantity
			byProduct[m.ProductID]["in"] = byProduct[m.ProductID]["in"].(int) + m.Quantity
		} else if m.Type == "out" {
			totalOut += m.Quantity
			byProduct[m.ProductID]["out"] = byProduct[m.ProductID]["out"].(int) + m.Quantity
		}
	}

	result := []map[string]interface{}{}
	for _, p := range byProduct {
		result = append(result, p)
	}

	return c.JSON(fiber.Map{
		"total_in":   totalIn,
		"total_out":  totalOut,
		"by_product": result,
	})
}


func GetSalesReport(c *fiber.Ctx) error {
	branchID := c.Query("branch_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if branchID == "" || startDate == "" || endDate == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing query parameters"})
	}

	report, err := services.FetchSalesReport(branchID, startDate, endDate)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(report)
}



func GetTodaySales(c *fiber.Ctx) error {
	branchID := c.Params("branch_id")
	result, err := services.FetchTodaySales(branchID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(result)
}



func GetTopProductsToday(c *fiber.Ctx) error {
	branchID := c.Params("branch_id")
	result, err := services.FetchTopProductsToday(branchID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(result)
}

func GetWeeklySales(c *fiber.Ctx) error {
	branchID := c.Params("branch_id")
	data, err := services.FetchWeeklySales(branchID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(data)
}


func GetTransactionsDaily(c *fiber.Ctx) error {
	branchID := c.Params("branch_id")
	data, err := services.FetchTransactionsDaily(branchID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(data)
}

func GetAverageTransactionValue(c *fiber.Ctx) error {
	branchID := c.Params("branch_id")
	data, err := services.FetchAverageTransactionValue(branchID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"average_transaction_value": data})
}

func GetLowStockAlert(c *fiber.Ctx) error {
	branchID := c.Params("branch_id")
	data, err := services.FetchLowStockAlert(branchID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(data)
}

func GetTopProducts(c *fiber.Ctx) error {
	branchID := c.Params("branch_id")
	data, err := services.FetchTopProducts(branchID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(data)
}


func GetSalesChart(c *fiber.Ctx) error {
	branchID := c.Params("branch_id")
	data, err := services.FetchSalesChart(branchID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(data)
}


func GetLowStock(c *fiber.Ctx) error {
	branchID := c.Params("branch_id")
	data, err := services.FetchLowStock(branchID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(data)
}


