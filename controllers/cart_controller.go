package controllers

import (
    "strconv"

    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v5"
    "github.com/pranavpatil6/go_mart/database"
    "github.com/pranavpatil6/go_mart/models"
)

func AddToCart(c *fiber.Ctx) error {
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
        ProductID uint `json:"product_id"`
        Quantity  int  `json:"quantity"`
    }
    if err := c.BodyParser(&input); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON input"})
    }
    if input.Quantity < 1 {
        input.Quantity = 1 // Default minimum quantity
    }

    // Verify product exists
    var product models.Product
    if err := database.DB.First(&product, input.ProductID).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
    }

    // Get or create user's cart
    var cart models.Cart
    err := database.DB.Preload("Items").Where("user_id = ?", userID).First(&cart).Error
    if err != nil {
        // Cart not found: create new
        cart = models.Cart{UserID: userID}
        if err := database.DB.Create(&cart).Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create cart"})
        }
    }

    // Check if product already in cart
    var cartItem models.CartItem
    found := false
    for _, item := range cart.Items {
        if item.ProductID == input.ProductID {
            cartItem = item
            found = true
            break
        }
    }

    if found {
        // Update quantity and price
        cartItem.Quantity += input.Quantity
        cartItem.Price = product.Price
        if err := database.DB.Save(&cartItem).Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update cart item"})
        }
    } else {
        // Add new cart item
        cartItem = models.CartItem{
            CartID:    cart.ID,
            ProductID: input.ProductID,
            Quantity:  input.Quantity,
            Price:     product.Price,
        }
        if err := database.DB.Create(&cartItem).Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add item to cart"})
        }
    }

    // Recalculate and update cart total
    if err := updateCartTotal(&cart); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update cart total"})
    }

    // Reload cart with fresh data
    if err := database.DB.Preload("Items.Product").First(&cart, cart.ID).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch updated cart"})
    }

    return c.JSON(cart)
}


// RemoveCartItem removes a cart item by ID
func RemoveCartItem(c *fiber.Ctx) error {
    userClaims, ok := c.Locals("user").(jwt.MapClaims)
    if !ok {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid JWT claims"})
    }

    userIDFloat, ok := userClaims["id"].(float64)
    if !ok {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID in token"})
    }
    userID := uint(userIDFloat)

    // Cart item ID from URL
    idStr := c.Params("id")
    itemID, err := strconv.Atoi(idStr)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid cart item ID"})
    }

    var cartItem models.CartItem
    if err := database.DB.First(&cartItem, itemID).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Cart item not found"})
    }

    // Confirm cart belongs to user
    var cart models.Cart
    if err := database.DB.First(&cart, cartItem.CartID).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cart not found"})
    }
    if cart.UserID != userID {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You cannot delete items from another user's cart"})
    }

    // Delete the cart item
    if err := database.DB.Delete(&cartItem).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to remove cart item"})
    }

    // Update cart total
    if err := updateCartTotal(&cart); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update cart total"})
    }

    return c.SendStatus(fiber.StatusNoContent) // 204 No Content
}

// ViewCart returns the user's full cart with items and product details
func ViewCart(c *fiber.Ctx) error {
    userClaims, ok := c.Locals("user").(jwt.MapClaims)
    if !ok {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid JWT claims"})
    }

    userIDFloat, ok := userClaims["id"].(float64)
    if !ok {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID in token"})
    }
    userID := uint(userIDFloat)

    var cart models.Cart
    if err := database.DB.Preload("Items.Product").Where("user_id = ?", userID).First(&cart).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Cart not found"})
    }

    return c.JSON(cart)
}

// updateCartTotal recalculates the cart total price based on cart items
func updateCartTotal(cart *models.Cart) error {
    var items []models.CartItem
    if err := database.DB.Where("cart_id = ?", cart.ID).Find(&items).Error; err != nil {
        return err
    }

    var total float64 = 0
    for _, item := range items {
        total += float64(item.Quantity) * item.Price
    }

    cart.Total = total
    return database.DB.Save(cart).Error
}
