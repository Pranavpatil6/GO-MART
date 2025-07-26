package models

import (
    "gorm.io/gorm"
)

type Cart struct {
    gorm.Model
    UserID uint       `gorm:"not null;index"`
    User   User       `gorm:"foreignKey:UserID"`
    Items  []CartItem `gorm:"foreignKey:CartID"` 
    Total  float64
}

type CartItem struct {
    gorm.Model
    CartID    uint    `gorm:"not null;index"`        
    ProductID uint    `gorm:"not null;index"`          
    Product   Product `gorm:"foreignKey:ProductID"`    
    Quantity  int     `gorm:"not null"`
    Price     float64 `gorm:"not null"`
}
