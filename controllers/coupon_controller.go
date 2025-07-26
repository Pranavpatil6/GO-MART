package controllers

import (
	"time"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pranavpatil6/go_mart/database"
	"github.com/pranavpatil6/go_mart/models"
)

func CreateCoupon(c *fiber.Ctx) error {
	var coupon models.Coupon
	if err := c.BodyParser(&coupon); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}


	if coupon.Code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Coupon code is required"})
	}
	if coupon.Discount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Discount must be positive"})
	}
	if coupon.UsageLimit <= 0 {
		coupon.UsageLimit = 100
	}
	if coupon.Expirydate.IsZero() {
		coupon.Expirydate = time.Now().AddDate(0, 1, 0) 
	}

	coupon.TimesUsed = 0 

	if err := database.DB.Create(&coupon).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create coupon (may be duplicate code)"})
	}

	return c.Status(fiber.StatusCreated).JSON(coupon)
}

func GetCoupons(c *fiber.Ctx) error {
	var coupons []models.Coupon
	if err := database.DB.Find(&coupons).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch coupons"})
	}
	return c.JSON(coupons)
}

func GetCouponByCode(c *fiber.Ctx) error {
	code := c.Params("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Coupon code is required"})
	}

	var coupon models.Coupon
	if err := database.DB.Where("code = ?", code).First(&coupon).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Coupon not found"})
	}
	return c.JSON(coupon)
}

func ApplyCoupon(c *fiber.Ctx) error {
	userClaims, ok := c.Locals("user").(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid JWT claims"})
	}
	userIDFloat, ok := userClaims["id"].(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID in token"})
	}
	userID := uint(userIDFloat)

	var input struct {
		Code string `json:"code"`
	}
	if err := c.BodyParser(&input); err != nil || input.Code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Coupon code is required"})
	}

	var cart models.Cart
	if err := database.DB.Preload("Items").Where("user_id = ?", userID).First(&cart).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Cart not found"})
	}

	cartTotal := cart.Total

	var coupon models.Coupon
	if err := database.DB.Where("code = ?", input.Code).First(&coupon).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Coupon not found"})
	}

	if time.Now().After(coupon.Expirydate) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Coupon expired"})
	}

	if coupon.TimesUsed >= coupon.UsageLimit {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Coupon usage limit exceeded"})
	}

	if cartTotal < coupon.MinCartValue {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cart total does not meet minimum value for coupon"})
	}

	var discountAmount float64
	switch coupon.Type {
	case "percent":
		discountAmount = (float64(coupon.Discount) / 100) * cartTotal
	case "fixed":
		discountAmount = float64(coupon.Discount)
	default:
		discountAmount = 0
	}

	// Discount should not exceed the cart total
	if discountAmount > cartTotal {
		discountAmount = cartTotal
	}


	return c.JSON(fiber.Map{
		"code":           coupon.Code,
		"totalAfter":     cartTotal - discountAmount,
		"discountAmount": discountAmount,
		"finalPrice":    cartTotal,
	})
}

func DeleteCoupon(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Coupon ID required"})
	}
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid coupon ID"})
	}

	if err := database.DB.Delete(&models.Coupon{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete coupon"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
