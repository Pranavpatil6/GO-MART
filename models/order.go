package models

type Order struct {
	Id     uint
	UserId uint
	Total  float64
	Status string
	Items  []OrderItem
}

type OrderItem struct {
	Id        uint `json:"primaryKey"`
	OrderId   uint
	ProductId uint
	Quantity  int
	Price     float64
}
