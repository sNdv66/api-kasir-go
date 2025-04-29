package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"api-kasir/utils"
	routes "api-kasir/router"
	
)

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

var verboseMode bool

func run() error {
	command := os.Args[1]
    
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

	// Validate required environment variables
	requiredEnvVars := []string{"SUPABASE_URL", "SUPABASE_API_KEY"}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			return fmt.Errorf("%s is required", envVar)
		}
	}

	supabaseUrl := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_API_KEY")

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Initialize dependencies
	utils.Init()

	// Verbose logging
	
	if verboseMode {
	fmt.Printf("%s %s\n", red("Supabase URL:"), supabaseUrl)
	fmt.Printf("%s %s*****\n", red("Supabase Key:"), supabaseKey[:3])
	fmt.Printf("%s %s\n", red("Server starting on port:"), port)
	}

	// Create Fiber app with configuration
	app := fiber.New(fiber.Config{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	})

	// Middlewares
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
	}))

	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "Asia/Jakarta",
	}))

	// Setup routes
	routes.SetRoutes(app)

	// Graceful Shutdown
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-shutdownChan
		log.Println(red("Shutting down server gracefully..."))
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = app.ShutdownWithContext(ctx)
	}()

	// Print startup message
	fmt.Println()
	fmt.Println(cyan("🚀 Server running on http://127.0.0.1:") + port)
	fmt.Println(yellow("🔗 Connected to Supabase project"))
	fmt.Println()

	return app.Listen(":" + port)
}



func showWelcomeMessage() {
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	
	fmt.Println()
	fmt.Println(cyan("   ╔══════════════════════════════════════════╗"))
	fmt.Println(cyan("   ║                                          ║"))
	fmt.Println(cyan("   ║          ") + yellow("A P I   K A S I R") + cyan("          ║"))
	fmt.Println(cyan("   ║                                          ║"))
	fmt.Println(cyan("   ╚══════════════════════════════════════════╝"))
	fmt.Println()
}

func showUsage() {
    green := color.New(color.FgGreen).SprintFunc()
    magenta := color.New(color.FgMagenta).SprintFunc()
    cyan := color.New(color.FgCyan).SprintFunc()
    yellow := color.New(color.FgYellow).SprintFunc()

    fmt.Println(yellow("\n  ╭───────────────────────────────────────────╮"))
    fmt.Println(yellow("  │  ") + green("🚀  CARA PENGGUNAAN API") + yellow("           │"))
    fmt.Println(yellow("  ╰───────────────────────────────────────────╯"))
    
    fmt.Println(magenta("\n  ┌─────────── PERINTAH UTAMA ──────────────┐"))
    fmt.Printf("  │ %-45s │\n", cyan("▶  Menjalankan server:"))
    fmt.Printf("  │ %-45s │\n", "    go run main.go start")
    fmt.Println(magenta("  ├───────────────────────────────────────────┤"))
    fmt.Printf("  │ %-45s │\n", cyan("▶  Menampilkan endpoint:"))
    fmt.Printf("  │ %-45s │\n", "    go run main.go help")
    fmt.Println(magenta("  └───────────────────────────────────────────┘"))
    
    fmt.Println(magenta("\n  ┌─────────── VARIABEL LINGKUNGAN ────────┐"))
    fmt.Printf("  │ %-45s │\n", cyan("🔧  SUPABASE_URL"))
    fmt.Printf("  │ %-45s │\n", cyan("🔑  SUPABASE_API_KEY"))
    fmt.Printf("  │ %-45s │\n", cyan("🎚️  PORT (opsional, default: 3000)"))
    fmt.Println(magenta("  └───────────────────────────────────────────┘"))
    
    fmt.Println(green("\n  💡 Tips: Buat file .env untuk menyimpan konfigurasi"))
    fmt.Println(yellow("  ╰───────────────────────────────────────────╯\n"))
}

func showEndpoints() {
    green := color.New(color.FgGreen).SprintFunc()
    cyan := color.New(color.FgCyan).SprintFunc()
    yellow := color.New(color.FgYellow).SprintFunc()
    
    fmt.Println(green("\n╔══════════════════════════════════════╗"))
    fmt.Println(green("║       ENDPOINT API YANG TERSEDIA        ║"))
    fmt.Println(green("╚══════════════════════════════════════╝"))
    fmt.Println()
    fmt.Println(yellow("  ┌──────────────────────────────────────┐"))
    fmt.Printf("  │ %-10s %-25s │\n", cyan("GET"), "/branches : Mendapatkan data cabang")
    fmt.Println(yellow("  ├──────────────────────────────────────┤"))
    fmt.Printf("  │ %-10s %-25s │\n", cyan("POST"), "/branches : Menambahkan cabang baru")
    fmt.Println(yellow("  ├──────────────────────────────────────┤"))
    fmt.Printf("  │ %-10s %-25s │\n", cyan("PUT"), "/branches/:id : Edit cabang ")
    fmt.Println(yellow("  ├──────────────────────────────────────┤"))
    
    fmt.Printf("  │ %-10s %-25s │\n", cyan("DELETE"), "/branches/:id : Menghapus cabang")
    fmt.Println(yellow("  └──────────────────────────────────────┘"))
    fmt.Printf("  │ %-10s %-25s │\n", cyan("POST"), "/login : dengan token jwt")
    fmt.Println(yellow("  └──────────────────────────────────────┘"))
    fmt.Println()
    fmt.Println(green("  Use 'curl' or tools like Postman to access these endpoints"))
    fmt.Println()
    fmt.Println(green("  go run main.go start --verbose untuk debug "))
    fmt.Println()
    
}