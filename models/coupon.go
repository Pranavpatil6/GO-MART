package models

import (
	"time"
)

type Coupon struct {
	CouponID uint `gorm:"primaryKey"`
	Code string `gorm:"uniqueIndex"`
	Discount int
	MinCartValue float64
	Expirydate time.Time
	Createddate time.Time
	TimesUsed int
	UsageLimit int 
	ProductId int 
	 
}