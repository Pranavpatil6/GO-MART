package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pranavpatil6/go_mart/database"
	"github.com/pranavpatil6/go_mart/models"

	"strconv"
)

// CreateOrder creates a new order for a user (checkout)
func CreateOrder(c *fiber.Ctx) error {
    type OrderInput struct {
        UserID uint `json:"user_id"` // In production, get from JWT!
    }
    var input OrderInput
    if err := c.BodyParser(&input); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
    }

    // Get user's cart and items
    var cart models.Cart
    if err := database.DB.Preload("Items.Product").Where("user_id = ?", input.UserID).First(&cart).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Cart not found"})
    }
    if len(cart.Items) == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cart is empty"})
    }

    // Create order items from cart items
    var orderItems []models.OrderItem
    var total float64
    for _, ci := range cart.Items {
        oi := models.OrderItem{
            ProductId: ci.ProductID,
            Quantity:  ci.Quantity,
            Price:     ci.Price,
        }
        orderItems = append(orderItems, oi)
        total += float64(ci.Quantity) * ci.Price
    }

    // Create new order
    order := models.Order{
        UserId: input.UserID,
        Total:  total,
        Status: "pending",
        Items:  orderItems,
    }
    if err := database.DB.Create(&order).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create order"})
    }

    // Optionally, clear user's cart after order
    if err := database.DB.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err == nil {
        cart.Total = 0
        database.DB.Save(&cart)
    }

    return c.Status(fiber.StatusCreated).JSON(order)
}

// GetOrders retrieves a list of a user's orders
func GetOrders(c *fiber.Ctx) error {
    userIdStr := c.Query("user_id")
    userId, err := strconv.Atoi(userIdStr)
    if err != nil || userId < 1 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
    }

    var orders []models.Order
    if err := database.DB.Preload("Items").Where("user_id = ?", userId).Find(&orders).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get orders"})
    }

    return c.JSON(orders)
}

// GetOrderById returns full details of a single order
func GetOrderById(c *fiber.Ctx) error {
    idStr := c.Params("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid order ID"})
    }

    var order models.Order
    if err := database.DB.Preload("Items").First(&order, id).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Order not found"})
    }

    return c.JSON(order)
}

