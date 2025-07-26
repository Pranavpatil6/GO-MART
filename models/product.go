package models


type Product struct {
    ProductId   uint    `gorm:"primaryKey"`
    Title       string
    Description string
    Price       float64
    Stock       int
}
