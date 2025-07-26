package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/pranavpatil6/go_mart/database"
	"github.com/pranavpatil6/go_mart/models"
)


func GetAllProducts(c *fiber.Ctx) error {
    var products []models.Product
    result := database.DB.Find(&products)
    if result.Error != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to retrieve products",
        })
    }

    return c.JSON(products)
}

func GetProduct(c *fiber.Ctx) error {
    idParam := c.Params("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
    }

    var product models.Product
    result := database.DB.First(&product, id)
    if result.Error != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
    }

    return c.JSON(product)
}

func CreateProduct(c *fiber.Ctx) error {
    var product models.Product

    if err := c.BodyParser(&product); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
    }

    result := database.DB.Create(&product)
    if result.Error != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create product"})
    }

    return c.Status(fiber.StatusCreated).JSON(product)
}

func UpdateProduct(c *fiber.Ctx) error {
    idParam := c.Params("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
    }

    var product models.Product
    result := database.DB.First(&product, id)
    if result.Error != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
    }

    var updateData models.Product
    if err := c.BodyParser(&updateData); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
    }

    product.Title = updateData.Title
    product.Description = updateData.Description
    product.Price = updateData.Price
    product.Stock = updateData.Stock

    saveResult := database.DB.Save(&product)
    if saveResult.Error != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update product"})
    }

    return c.JSON(product)
}

func DeleteProduct(c *fiber.Ctx) error {
    idParam := c.Params("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
    }

    var product models.Product
    result := database.DB.First(&product, id)
    if result.Error != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
    }

    deleteResult := database.DB.Delete(&product)
    if deleteResult.Error != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete product"})
    }

    return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"message": "Product Deleted"})
}
