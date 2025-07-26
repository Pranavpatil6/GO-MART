package controllers

import (
	"os"
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pranavpatil6/go_mart/database"
	"github.com/pranavpatil6/go_mart/models"
	"golang.org/x/crypto/bcrypt"
)


func Register(c *fiber.Ctx) error {
    
	var input models.User
    if err := c.BodyParser(&input); err != nil {
        return c.Status(422).JSON(fiber.Map{"error": "Invalid JSON input"})
    }

    if input.Name == "" || input.Email == "" || input.Password == "" {
        return c.Status(422).JSON(fiber.Map{"error": "Missing required field"})
    }

    var exists models.User
    if err := database.DB.Where("email = ?", input.Email).First(&exists).Error; err == nil {
        return c.Status(400).JSON(fiber.Map{"error": "Email already registered"})
    }

    hash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 14)
    user := models.User{
        Name:     input.Name,
        Email:    input.Email,
        Password: string(hash),
        Role:     "user", // default role
    }
    if err := database.DB.Create(&user).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to create user"})
    }
    return c.Status(201).JSON(fiber.Map{"message": "User registered successfully"})
}

func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := os.Getenv("JWT_SECRET")
	t, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(fiber.Map{"token": t})
}
