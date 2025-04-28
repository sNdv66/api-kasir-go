package routes

import (
	"github.com/gofiber/fiber/v2"
	"api-kasir/handlers"
	"api-kasir/middleware"
	jwtware "github.com/gofiber/jwt/v3"
	jwt "github.com/golang-jwt/jwt/v4"
	
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

        return c.JSON(fiber.Map{
            "email":     email,
            "role":      role,
            "branch_id": branchID,
        })
    })
}


func SetRoutes(app *fiber.App) {
	//app.Get("/branches", handlers.GetBranches)
	app.Post("/branches", handlers.CreateBranch)
	app.Put("/branches/:id", handlers.UpdateBranch)
	app.Delete("/branches/:id", handlers.DeleteBranch)
	app.Post("/login", handlers.Login)
	SetupProtectedRoutes(app)
    app.Get("/branches", middleware.JWTMiddleware, handlers.GetBranches)
	
	
}


