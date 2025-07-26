package middleware

import (
    "os"
    "strings"

    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v4"
)


func JWTProtected() fiber.Handler {
    return func(c *fiber.Ctx) error {
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Authorization header"})
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == authHeader {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Malformed token"})
        }

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(os.Getenv("JWT_SECRET")), nil
        })

        if err != nil || !token.Valid {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
        }

        c.Locals("user", claims)
        return c.Next()
    }
}


func AdminOnly() fiber.Handler {
    return func(c *fiber.Ctx) error {
        userClaims, ok := c.Locals("user").(jwt.MapClaims)
        if !ok {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
        }
        role, ok := userClaims["role"].(string)
        if !ok || role != "admin" {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Admin access required"})
        }
        return c.Next()
    }
}
