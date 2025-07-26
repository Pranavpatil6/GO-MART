package controllers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pranavpatil6/go_mart/database"
	"github.com/pranavpatil6/go_mart/models"
)

// CreateCoupon: Adds a new coupon
func CreateCoupon(c *fiber.Ctx) error {
    var input models.Coupon
    if err := c.BodyParser(&input); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
    }
    if input.Code == "" || input.Discount <= 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Code and positive discount required"})
    }
    input.Createddate = time.Now()
    if err := database.DB.Create(&input).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create coupon"})
    }
    return c.Status(fiber.StatusCreated).JSON(input)
}

// GetCoupons: List all coupons
func GetCoupons(c *fiber.Ctx) error {
    var coupons []models.Coupon
    if err := database.DB.Find(&coupons).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not fetch coupons"})
    }
    return c.JSON(coupons)
}

// GetCouponByCode: Retrieve a single coupon by code
func GetCouponByCode(c *fiber.Ctx) error {
    code := c.Params("code")
    var coupon models.Coupon
    if err := database.DB.Where("code = ?", code).First(&coupon).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Coupon not found"})
    }
    return c.JSON(coupon)
}

// ApplyCoupon: Validates application of coupon (does not increment TimesUsed, just checks validity)
func ApplyCoupon(c *fiber.Ctx) error {
    type ApplyInput struct {
        Code       string  `json:"code"`
        CartTotal  float64 `json:"cart_total"`
    }
    var req ApplyInput
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
    }

    var coupon models.Coupon
    if err := database.DB.Where("code = ?", req.Code).First(&coupon).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Coupon not found"})
    }

    if time.Now().After(coupon.Expirydate) {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Coupon expired"})
    }
    if coupon.TimesUsed >= coupon.UsageLimit {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Usage limit exceeded"})
    }
    if req.CartTotal < coupon.MinCartValue {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cart value does not meet minimum"})
    }


    var discountValue float64
    discountValue = float64(coupon.Discount) // use as percent for now
    discounted := req.CartTotal * ((100 - discountValue) / 100.0)

    return c.JSON(fiber.Map{
        "code":         coupon.Code,
        "original":     req.CartTotal,
        "discounted":   discounted,
        "discountType": "percent", // update if you distinguish fixed/percent
    })
}


func DeleteCoupon(c *fiber.Ctx) error {
    idStr := c.Params("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid coupon ID"})
    }
    if err := database.DB.Delete(&models.Coupon{}, id).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete coupon"})
    }
    return c.SendStatus(fiber.StatusNoContent)
}
