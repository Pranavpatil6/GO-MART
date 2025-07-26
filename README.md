# GO-MART 🛒

Go-Mart is a modern, scalable e-commerce backend built using Go (Golang), Fiber web framework, GORM ORM, and PostgreSQL database. The project features user authentication with JWT, product management, shopping cart, orders, and coupon system with role-based access control.

---

## 🚀 Features

**User Authentication:** Register, login with hashed passwords and JWT-based authentication.
- **Role-Based Authorization:** Roles like "user" and "admin" with protected routes.
- **Product Management:** CRUD operations for products (admin only).
- **Shopping Cart:** Add, update, remove items from cart; view cart.
- **Order Processing:** Checkout and view order history.
- **Coupons:** Create, delete (admin only), apply coupons to carts.
- **Middleware:** JWT middleware for securing API endpoints.
- **PostgreSQL:** Reliable relational database.
- **Modular Codebase:** Clearly separated layers (models, controllers, routes, middleware, utils).

## 🛠️ Tech Stack

| Component      | Technology               |
|----------------|---------------------------|
| Language       | [Go (Golang)](https://golang.org) |
| Web Framework  | [Fiber](https://gofiber.io/)       |
| ORM            | [GORM](https://gorm.io/)          |
| Database       | PostgreSQL               |
| Auth           | JWT (JSON Web Tokens)    |

---

## 📁 Project Structure

ecommerce/
├── go.mod
├── main.go
├── database/
│ └── db.go
├── models/
│ ├── user.go
│ ├── product.go
│ ├── cart.go
│ ├── order.go
│ └── coupon.go
├── controllers/
│ ├── auth_controller.go
│ ├── product_controller.go
│ ├── cart_controller.go
│ ├── order_controller.go
│ └── coupon_controller.go
├── routes/
│ └── routes.go
├── middleware/
  └── auth.go

## Technologies Used

- **Golang** — Backend language
- **Fiber** — Web framework
- **GORM** — ORM for DB abstraction
- **PostgreSQL** — Database
- **JWT** — Authentication tokens
- **bcrypt** — Password hashing

---

Thank you for using Go-Mart! 🚀