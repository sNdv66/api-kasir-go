package middleware

import (
    "github.com/gofiber/fiber/v2"
    "api-kasir/utils"
)

func JWTMiddleware(c *fiber.Ctx) error {
    authHeader := c.Get("Authorization")
    if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid token"})
    }

    token := authHeader[7:]
    claims, err := utils.ValidateJWT(token)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
    }

    // Simpan data dari token ke context
    c.Locals("user_id", claims.UserID)
    c.Locals("email", claims.Email)
    c.Locals("role", claims.Role)
    c.Locals("branch_id", claims.BranchID)

    return c.Next()
}