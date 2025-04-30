package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"api-kasir/router"
	"api-kasir/utils"
    "encoding/json"
	"github.com/fatih/color"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

var verboseMode bool

func main() {
	if len(os.Args) == 1 {
		showWelcomeMessage()
		showUsage()
		os.Exit(0)
	}

	if err := run(); err != nil {
		log.Fatalf("🚨 Application error: %v", err)
	}
}

func run() error {
	command := os.Args[1]

	// Cek apakah --verbose diaktifkan
	for _, arg := range os.Args {
		if arg == "--verbose" {
			verboseMode = true
			break
		}
	}

	switch command {
	case "start":
		return startAPIServer()
	case "help":
		showEndpoints()
		return nil
	case "detail":
		if len(os.Args) < 3 {
			return fmt.Errorf("masukkan path endpoint setelah perintah 'detail'")
		}
		return showEndpointDetail(os.Args[2])
	default:
		showUsage()
		return fmt.Errorf("command tidak dikenali: %s", command)
	}
}

func startAPIServer() error {
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}

	requiredEnvVars := []string{"SUPABASE_URL", "SUPABASE_API_KEY"}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			return fmt.Errorf("%s is required", envVar)
		}
	}

	supabaseUrl := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_API_KEY")
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	utils.Init()

	if verboseMode {
		fmt.Printf("%s %s\n", red("Supabase URL:"), supabaseUrl)
		fmt.Printf("%s %s*****\n", red("Supabase Key:"), supabaseKey[:3])
		fmt.Printf("%s %s\n", red("Server starting on port:"), port)
	}

	// Setup Fiber app
	app := fiber.New(fiber.Config{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
	}))

	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "Asia/Jakarta",
	}))

	// Set routes
	routes.SetRoutes(app)  // Ubah ini sesuai dengan nama package Anda

	// Handle graceful shutdown
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-shutdownChan
		log.Println(red("Shutting down server gracefully..."))
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = app.ShutdownWithContext(ctx)
	}()

	// Start server
	fmt.Println()
	fmt.Println(cyan("🚀 Server running on http://127.0.0.1:") + port)
	fmt.Println(yellow("🔗 Connected to Supabase project"))
	fmt.Println()

	return app.Listen(":" + port)
}

func showWelcomeMessage() {
	fmt.Println("=== Selamat datang di API Kasir CLI ===")
}

func showUsage() {
	fmt.Println("Penggunaan:")
	fmt.Println("  api-kasir start            - Menjalankan API server")
	fmt.Println("  api-kasir help             - Menampilkan daftar endpoint")
	fmt.Println("  api-kasir detail [path]    - Menampilkan detail endpoint")
	fmt.Println("  --verbose                  - Menampilkan informasi detail saat start")
}

func showEndpoints() {
	fmt.Println("Daftar Endpoint yang tersedia:")
	fmt.Println("  GET    /branches")
	fmt.Println("  POST   /branches")
	fmt.Println("  PUT    /branches/:id")
	fmt.Println("  DELETE /branches/:id")
	fmt.Println("  POST   /login")
	fmt.Println("  POST   /products")
}

type EndpointDetail struct {
	Method      string
	Path        string
	Description string
	Input       any
	Headers     string
	Output      any
	Errors      []map[string]interface{} // Changed from ErrorResponses to Errors for consistency
}

var endpointDetails = []EndpointDetail{
	{
		Method:      "POST",
		Path:        "/login",
		Description: "Login dan mendapatkan JWT token",
		Input: map[string]string{
			"email":    "admin@example.com",
			"password": "123456",
		},
		Headers: "Content-Type: application/json",
		Output: map[string]any{
			"token": "jwt_token",
			"user": map[string]string{
				"id":   "uuid",
				"role": "admin",
			},
		},
		Errors: []map[string]interface{}{
			{
				"code":   400,
				"error":  "Invalid request",
				"detail": "Format email atau password salah",
			},
			{
				"code":   401,
				"error":  "Unauthorized",
				"detail": "Email atau password tidak valid",
			},
		},
	},
	{
		Method:      "GET",
		Path:        "/branches",
		Description: "Mendapatkan daftar semua cabang",
		Headers:     "Authorization: Bearer <token>\nContent-Type: application/json",
		Output: []map[string]interface{}{
			{
				"id":        "branch_123",
				"name":      "Cabang Jakarta",
				"address":   "Jl. Sudirman No. 1",
				"createdAt": "2023-01-01T00:00:00Z",
			},
		},
		Errors: []map[string]interface{}{
			{
				"code":   401,
				"error":  "Unauthorized",
				"detail": "Token tidak valid atau kadaluarsa",
			},
			{
				"code":   500,
				"error":  "Internal Server Error",
				"detail": "Gagal mengambil data cabang",
			},
		},
	},
	{
		Method:      "POST",
		Path:        "/branches",
		Description: "Membuat cabang baru",
		Headers:     "Authorization: Bearer <token>\nContent-Type: application/json",
		Input: map[string]interface{}{
			"name":    "Nama Cabang",
			"address": "Alamat Lengkap",
		},
		Output: map[string]interface{}{
			"message": "Cabang berhasil ditambahkan",
			"branch": map[string]interface{}{
				"id":        "generated-uuid",
				"name":      "Nama Cabang",
				"address":   "Alamat Lengkap",
				"createdAt": "2023-01-01T00:00:00Z",
			},
		},
		Errors: []map[string]interface{}{
			{
				"code":   400,
				"error":  "Bad Request",
				"detail": "Nama atau alamat tidak boleh kosong",
			},
		},
	},
	{
		Method:      "PUT",
		Path:        "/branches/:id",
		Description: "Memperbarui data cabang",
		Headers:     "Authorization: Bearer <token>\nContent-Type: application/json",
		Input: map[string]interface{}{
			"name":    "Nama Cabang Updated",
			"address": "Alamat Updated",
		},
		Output: map[string]interface{}{
			"message": "Cabang berhasil diperbarui",
			"branch": map[string]interface{}{
				"id":        "existing-uuid",
				"name":      "Nama Cabang Updated",
				"address":   "Alamat Updated",
				"updatedAt": "2023-01-02T00:00:00Z",
			},
		},
	},
	{
		Method:      "DELETE",
		Path:        "/branches/:id",
		Description: "Menghapus cabang",
		Headers:     "Authorization: Bearer <token>\nContent-Type: application/json",
		Output: map[string]interface{}{
			"message": "Cabang berhasil dihapus",
			"deletedAt": "2023-01-03T00:00:00Z",
		},
	},
    {
        Method:      "GET",
        Path:        "/branches/:branch_id/products",
        Description: "Mendapatkan daftar produk berdasarkan cabang",
        Headers:     "Authorization: Bearer <token>\nContent-Type: application/json",
        Output: []map[string]interface{}{
            {
                "id":        "prod_123",
                "name":      "Produk A",
                "price":     100000,
                "stock":     50,
                "branch_id": "branch_123",
                "created_at": "2023-01-01T00:00:00Z",
            },
        },
        Errors: []map[string]interface{}{
            {
                "code":    500,
                "error":   "Internal Server Error",
                "detail": "Gagal mengambil data produk",
            },
        },
    },
    // POST /products
    {
        Method:      "POST",
        Path:        "/products",
        Description: "Menambahkan produk baru",
        Headers:     "Authorization: Bearer <token>\nContent-Type: application/json",
        Input: map[string]interface{}{
            "name":      "Produk Baru",
            "price":     150000,
            "stock":     100,
            "branch_id": "branch_123",
        },
        Output: map[string]interface{}{
            "id":        "prod_124",
            "name":      "Produk Baru",
            "price":     150000,
            "stock":     100,
            "branch_id": "branch_123",
            "created_at": "2023-01-02T00:00:00Z",
        },
        Errors: []map[string]interface{}{
            {
                "code":    500,
                "error":   "Internal Server Error",
                "detail": "Gagal menambahkan produk",
            },
        },
    },
    // PATCH /products/:id
    {
        Method:      "PATCH",
        Path:        "/products/:id",
        Description: "Memperbarui data produk",
        Headers:     "Authorization: Bearer <token>\nContent-Type: application/json",
        Input: map[string]interface{}{
            "price": 175000,
            "stock": 75,
        },
        Output: map[string]interface{}{
            "id":        "prod_123",
            "name":      "Produk A",
            "price":     175000,
            "stock":     75,
            "branch_id": "branch_123",
            "updated_at": "2023-01-03T00:00:00Z",
        },
        Errors: []map[string]interface{}{
            {
                "code":    500,
                "error":   "Internal Server Error",
                "detail": "Gagal memperbarui produk",
            },
        },
    },
    // DELETE /products/:id
    {
        Method:      "DELETE",
        Path:        "/products/:id",
        Description: "Menghapus produk",
        Headers:     "Authorization: Bearer <token>\nContent-Type: application/json",
        Output: map[string]interface{}{
            "message": "Produk berhasil dihapus",
            "id":      "prod_123",
        },
        Errors: []map[string]interface{}{
            {
                "code":    500,
                "error":   "Internal Server Error",
                "detail": "Gagal menghapus produk",
            },
        },
    },

}

func showEndpointDetail(path string) error {
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	endpoints := []EndpointDetail{}
	for _, ep := range endpointDetails {
		if ep.Path == path {
			endpoints = append(endpoints, ep)
		}
	}

	if len(endpoints) == 0 {
		return fmt.Errorf(red("Endpoint tidak ditemukan: %s"), path)
	}

	fmt.Println()
	fmt.Println(green("🔍 Detail Endpoint"))
	fmt.Println("─────────────────────────────────────")
	
	for _, ep := range endpoints {
		fmt.Println(yellow("Method:     "), cyan(ep.Method))
		fmt.Println(yellow("Path:       "), ep.Path)
		fmt.Println(yellow("Deskripsi:  "), ep.Description)
		fmt.Println(yellow("Headers:    "), ep.Headers)
		
		if ep.Input != nil {
			fmt.Println(yellow("\nContoh Input:"))
			printPrettyJSON(ep.Input)
		}
		
		if ep.Output != nil {
			fmt.Println(yellow("\nContoh Output:"))
			printPrettyJSON(ep.Output)
		}
		
		if len(ep.Errors) > 0 {
			fmt.Println(yellow("\nError Responses:"))
			for _, err := range ep.Errors {
				fmt.Printf("  %s %s: %s\n", 
					red("Code:"), err["code"], err["error"])
				fmt.Printf("  %s %s\n\n", 
					yellow("Detail:"), err["detail"])
			}
		}
		
		fmt.Println("─────────────────────────────────────")
	}
	return nil
}

// Helper function to print pretty JSON
func printPrettyJSON(data interface{}) {
	jsonData, err := json.MarshalIndent(data, "  ", "  ")
	if err != nil {
		fmt.Println("  Error formatting JSON:", err)
		return
	}
	fmt.Println(string(jsonData))
}