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
	"github.com/fatih/color"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

var verboseMode bool

var (
    hiRedPrint     = color.New(color.FgHiRed).PrintlnFunc()
    hiGreenPrint   = color.New(color.FgHiGreen).PrintlnFunc()
    hiYellowPrint    = color.New(color.FgHiBlue).PrintlnFunc()
    
    hiMagentaPrint = color.New(color.FgHiMagenta).PrintlnFunc()
    hiCyanPrint    = color.New(color.FgHiCyan).PrintlnFunc()
    hiWhitePrint   = color.New(color.FgHiWhite).PrintlnFunc()
    hiBlackPrint   = color.New(color.FgHiBlack).PrintlnFunc()
)


func main() {
	if len(os.Args) == 1 {
		showWelcomeMessage()
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
		showUsage()
		return nil
	case "-m":
	     showEndpoints()
	     return nil
	case "-v":
	     showVersion()
	     return nil
	default:
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
    hiCyanPrint("Use commands :")
    
    hiYellowPrint("go run main.go help")
    
}

func showUsage() {
	hiCyanPrint("Perintah yang tersedia:")
	hiYellowPrint("  go run main.go [perintah]")

	hiCyanPrint("\nDaftar perintah:")
	hiYellowPrint("  start               -> Menjalankan server")
	hiYellowPrint("  start --verbose     -> Menjalankan server dengan mode debug")
	hiYellowPrint("  -m                  -> Menampilkan daftar endpoint")
	hiYellowPrint("  -v                  -> version")
}

func showVersion(){
    hiCyanPrint("version 0.1.1")
}

func showEndpoints() {
	hiCyanPrint("Daftar Endpoint yang tersedia saat ini:\n")

	// Helper untuk format
	printEP := func(method, path, desc string) {
		fmt.Printf("[%s]  %s\n      -> %s\n\n", method, path, desc)
	}

	// Autentikasi
	printEP("POST", "/login", "Login dan menerima token JWT")

	// Cabang
	printEP("GET", "/branches", "Mengambil daftar cabang")
	printEP("POST", "/branches", "Menambahkan cabang baru")
	printEP("PUT", "/branches/:id", "Memperbarui data cabang berdasarkan ID")
	printEP("DELETE", "/branches/:id", "Menghapus cabang berdasarkan ID")

	// Produk
	printEP("GET", "/branches/:branch_id/products", "Mengambil daftar produk berdasarkan cabang")
	printEP("POST", "/products", "Menambahkan produk baru")
	printEP("PATCH", "/products/:id", "Memperbarui data produk berdasarkan ID")
	printEP("DELETE", "/products/:id", "Menghapus produk berdasarkan ID")

	// Transaksi
	printEP("POST", "/transactions", "Membuat transaksi baru")
	printEP("GET", "/branches/:branch_id/transactions", "Mengambil daftar transaksi berdasarkan cabang")
	printEP("GET", "/transactions/:id", "Mengambil detail transaksi berdasarkan ID")

	// Item Transaksi
	printEP("POST", "/transaction-items", "Menambahkan item ke transaksi")
	printEP("GET", "/transactions/:id/items", "Mengambil daftar item berdasarkan transaksi ID")

	// Pergerakan Stok
	printEP("GET", "/stock-movements", "Mengambil histori pergerakan stok")
	printEP("POST", "/stock-movements", "Menambahkan pergerakan stok")
	printEP("GET", "/stock-summary", "Mengambil ringkasan stok")

	// Laporan
	printEP("GET", "/reports/sales", "Mendapatkan laporan penjualan")

	// Dashboard
	printEP("GET", "/branches/:branch_id/dashboard/today-sales", "Penjualan hari ini")
	printEP("GET", "/branches/:branch_id/dashboard/top-products", "Produk terlaris hari ini")
	printEP("GET", "/branches/:branch_id/dashboard/weekly-sales", "Penjualan mingguan")
	printEP("GET", "/branches/:branch_id/dashboard/sales-chart", "Grafik penjualan")
	printEP("GET", "/branches/:branch_id/dashboard/low-stock", "Produk dengan stok rendah")
}