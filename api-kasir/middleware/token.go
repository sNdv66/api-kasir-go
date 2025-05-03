package middleware

import (
    "github.com/gofiber/fiber/v2"
    "api-kasir/utils"
)

// JWT Middleware untuk memverifikasi token
func JWTMiddleware(c *fiber.Ctx) error {
    authHeader := c.Get("Authorization")
    if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
    return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid token"})
    }
    token := authHeader[7:]

    if token == "" {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid token"})
    }

    claims, err := utils.ValidateJWT(token)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
    }

    // Tambahkan data ke context untuk digunakan di handler
    c.Locals("user", claims)
    return c.Next()
}