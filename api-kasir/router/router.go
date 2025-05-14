package routes

import (
	"github.com/gofiber/fiber/v2"
	"api-kasir/handlers"
	"api-kasir/middleware"
	jwtware "github.com/gofiber/jwt/v3"
	jwt "github.com/golang-jwt/jwt/v4"
	"api-kasir/utils"
	"encoding/json"
)
func SetupProtectedRoutes(app *fiber.App) {
    // Middleware untuk semua endpoint yang butuh login
    app.Use(jwtware.New(jwtware.Config{
        SigningKey: []byte("582889w67888"), // Harus sama dengan yang dipakai saat buat token
        ErrorHandler: func(c *fiber.Ctx, err error) error {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Unauthorized",
            })
        },
    }))

    // Di bawah ini adalah contoh endpoint yang butuh token
    app.Get("/profile", func(c *fiber.Ctx) error {
    user := c.Locals("user").(*jwt.Token)
    claims := user.Claims.(jwt.MapClaims)

    email := claims["email"].(string)
    role := claims["role"].(string)
    branchID := claims["branch_id"].(string)

    // Ambil data cabang dari Supabase
    req, err := utils.NewRequest("GET", "branches?id=eq."+branchID, nil)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to create request",
        })
    }

    resp, err := utils.Client.Do(req)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to fetch branch data",
        })
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Supabase responded with status " + resp.Status,
        })
    }

    // Baca hasil response
    var branches []struct {
        ID      string `json:"id"`
        Name    string `json:"name"`
        Address string `json:"address"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&branches); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to decode branch data",
        })
    }

    if len(branches) == 0 {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Branch not found",
        })
    }

    branch := branches[0]

    // Kirim data lengkap
    return c.JSON(fiber.Map{
        "email":          email,
        "role":           role,
        "branch_id":      branchID,
        "branch_name":    branch.Name,
        "branch_address": branch.Address,
    })
})

}


func SetRoutes(app *fiber.App) {
    
	app.Post("/branches", handlers.CreateBranch)
	app.Put("/branches/:id", handlers.UpdateBranch)
	app.Delete("/branches/:id", handlers.DeleteBranch)
	app.Get("/branches", middleware.JWTMiddleware, handlers.GetBranches)
	app.Get("/branches/:id", middleware.JWTMiddleware, handlers.GetBranchByID)
	
	
	app.Post("/login", handlers.Login)
	SetupProtectedRoutes(app)
    
    
    app.Get("/branches/:branch_id/products", handlers.GetProductsByBranch)
    
    app.Post("/products", handlers.CreateProduct)
    app.Patch("/products/:id", handlers.UpdateProduct)
    app.Delete("/products/:id", handlers.DeleteProduct)
    app.Post("/transactions", handlers.CreateTransaction)
    app.Get("/branches/:branch_id/transactions", handlers.GetTransactionsByBranch)
    app.Get("/transactions/:id", handlers.GetTransactionByID)
    app.Post("/transaction-items", handlers.CreateTransactionItem)
    app.Get("/transactions/:id/items", handlers.GetTransactionItemsByTransactionID)
    app.Get("/stock-movements", handlers.GetStockMovements)
    app.Post("/stock-movements", handlers.CreateStockMovement)
    app.Get("/stock-summary", handlers.GetStockSummary)
    
    app.Get("/reports/sales", handlers.GetSalesReport)
    
    app.Get("/branches/:branch_id/dashboard/today-sales", handlers.GetTodaySales)
    
    app.Get("/branches/:branch_id/dashboard/top-products", handlers.GetTopProducts)
    app.Get("/branches/:branch_id/dashboard/sales-chart", handlers.GetSalesChart)
    app.Get("/branches/:branch_id/dashboard/low-stock", handlers.GetLowStock)
    //
	
	app.Get("/branches/:branch_id/dashboard/top-products", handlers.GetTopProductsToday)
	app.Get("/branches/:branch_id/dashboard/weekly-sales", handlers.GetWeeklySales)
	
	
	
}


