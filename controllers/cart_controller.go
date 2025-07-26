package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/pranavpatil6/go_mart/database"
	"github.com/pranavpatil6/go_mart/models"
)

// AddToCart adds a product to the user's cart or increments quantity if already present
func AddToCart(c *fiber.Ctx) error {
    type AddToCartInput struct {
        UserID    uint `json:"user_id"`    // In production, extract from JWT
        ProductID uint `json:"product_id"`
        Quantity  int  `json:"quantity"`
    }

    var input AddToCartInput
    if err := c.BodyParser(&input); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
    }
    if input.Quantity < 1 {
        input.Quantity = 1
    }

    // Verify product exists
    var product models.Product
    if err := database.DB.First(&product, input.ProductID).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
    }

    // Find or create cart for user
    var cart models.Cart
    err := database.DB.Preload("Items").Where("user_id = ?", input.UserID).First(&cart).Error
    if err != nil {
        // Cart not found, create new
        cart = models.Cart{UserID: input.UserID, Total: 0}
        if err := database.DB.Create(&cart).Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create cart"})
        }
    }

    // Check if item already exists in cart
    var cartItem models.CartItem
    itemFound := false
    for _, item := range cart.Items {
        if item.ProductID == input.ProductID {
            cartItem = item
            itemFound = true
            break
        }
    }

    if itemFound {
        // Update quantity
        cartItem.Quantity += input.Quantity
        cartItem.Price = product.Price // update price just in case
        if err := database.DB.Save(&cartItem).Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update cart item"})
        }
    } else {
        // Add new item
        cartItem = models.CartItem{
            CartID:    cart.ID,
            ProductID: input.ProductID,
            Quantity:  input.Quantity,
            Price:     product.Price,
        }
        if err := database.DB.Create(&cartItem).Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not add item to cart"})
        }
    }

    // Update cart total
    if err := updateCartTotal(&cart); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update cart total"})
    }

    // Reload the cart with items to return fresh data
    if err := database.DB.Preload("Items.Product").First(&cart, cart.ID).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch cart data"})
    }

    return c.JSON(cart)
}

// UpdateCart updates the quantity of a cart item
func UpdateCart(c *fiber.Ctx) error {
    type UpdateCartInput struct {
        UserID    uint `json:"user_id"`    // From JWT in real app
        ProductID uint `json:"product_id"`
        Quantity  int  `json:"quantity"`
    }

    var input UpdateCartInput
    if err := c.BodyParser(&input); err != nil || input.Quantity < 1 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input or quantity"})
    }

    // Find user's cart
    var cart models.Cart
    if err := database.DB.Where("user_id = ?", input.UserID).First(&cart).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Cart not found"})
    }

    // Find the cart item
    var cartItem models.CartItem
    if err := database.DB.Where("cart_id = ? AND product_id = ?", cart.ID, input.ProductID).First(&cartItem).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Cart item not found"})
    }

    cartItem.Quantity = input.Quantity
    if err := database.DB.Save(&cartItem).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update cart item"})
    }

    // Update cart total
    if err := updateCartTotal(&cart); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update cart total"})
    }

    // Return updated cart data
    if err := database.DB.Preload("Items.Product").First(&cart, cart.ID).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch updated cart"})
    }

    return c.JSON(cart)
}

// RemoveFromCart removes a cart item by its ID
func RemoveFromCart(c *fiber.Ctx) error {
    idParam := c.Params("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid cart item ID"})
    }

    // Find the cart item
    var cartItem models.CartItem
    if err := database.DB.First(&cartItem, id).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Cart item not found"})
    }

    // Find the cart to update total later
    var cart models.Cart
    if err := database.DB.First(&cart, cartItem.CartID).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cart not found"})
    }

    if err := database.DB.Delete(&cartItem).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to remove cart item"})
    }

    // Update cart total
    if err := updateCartTotal(&cart); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update cart total"})
    }

    return c.SendStatus(fiber.StatusNoContent) // 204
}

// ViewCart returns all the cart items for a specific user
func ViewCart(c *fiber.Ctx) error {
    userIDParam := c.Query("user_id")
    userID, err := strconv.Atoi(userIDParam)
    if err != nil || userID < 1 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
    }

    var cart models.Cart
    if err := database.DB.Preload("Items.Product").Where("user_id = ?", userID).First(&cart).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Cart not found"})
    }

    return c.JSON(cart)
}

// Helper function to update the total price of the cart
func updateCartTotal(cart *models.Cart) error {
    var total float64
    var cartItems []models.CartItem

    if err := database.DB.Where("cart_id = ?", cart.ID).Find(&cartItems).Error; err != nil {
        return err
    }

    for _, item := range cartItems {
        total += float64(item.Quantity) * item.Price
    }

    cart.Total = total
    return database.DB.Save(cart).Error
}
