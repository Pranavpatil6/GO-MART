package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pranavpatil6/go_mart/controllers"
	"github.com/pranavpatil6/go_mart/middleware"
)

func SetupRoutes(app *fiber.App) {
	app.Post("/register", controllers.Register)
    app.Post("/login", controllers.Login)
	app.Get("/users",controllers.GetAllUsers)

    // Product
    app.Get("/products", controllers.GetAllProducts)
    app.Get("/products/:id", controllers.GetProduct)
    app.Post("/products", middleware.JWTProtected(),middleware.AdminOnly(), controllers.CreateProduct)
    app.Put("/products/:id", middleware.JWTProtected(),middleware.AdminOnly(), controllers.UpdateProduct)
    app.Delete("/products/:id", middleware.JWTProtected(),middleware.AdminOnly(), controllers.DeleteProduct)

    // Cart
    cart := app.Group("/cart",middleware.JWTProtected())
    cart.Post("/add", controllers.AddToCart)
    cart.Delete("/remove/:id", controllers.RemoveCartItem)
    cart.Get("/", controllers.ViewCart)
    cart.Post("/apply-coupon", controllers.ApplyCoupon)

    // Coupon
    app.Post("/coupons", middleware.JWTProtected(), middleware.AdminOnly(), controllers.CreateCoupon)
    app.Get("/coupons", controllers.GetCoupons)
    app.Get("/coupons/:code", controllers.GetCouponByCode)
    app.Delete("/coupons/:id", middleware.JWTProtected(), middleware.AdminOnly(), controllers.DeleteCoupon)

    //Orders
    orders := app.Group("/orders", middleware.JWTProtected())
    orders.Post("/", controllers.CreateOrder)        
    orders.Get("/", controllers.GetOrders)


}
