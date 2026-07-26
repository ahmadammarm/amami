package qurban

import (
	"time"

	"github.com/google/uuid"
)

type CreatePackageRequest struct {
	Name       string `json:"name" binding:"required"`
	Type       string `json:"type" binding:"required"`
	Price      int64  `json:"price" binding:"required,gt=0"`
	YearHijri  int    `json:"year_hijri" binding:"required"`
	StockTotal int    `json:"stock_total" binding:"required,gt=0"`
}

type PackageResponse struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	Type           string `json:"type"`
	Price          int64  `json:"price"`
	YearHijri      int    `json:"year_hijri"`
	StockTotal     int    `json:"stock_total"`
	StockRemaining int    `json:"stock_remaining"`
}

type CreateBookingRequest struct {
	ShohibulID    uuid.UUID `json:"shohibul_id" binding:"required"`
	PackageID     uint      `json:"package_id" binding:"required"`
	PaymentStatus string    `json:"payment_status" binding:"required"`
	TotalAmount   int64     `json:"total_amount" binding:"required,gt=0"`
}

type BookingResponse struct {
	ID            uuid.UUID `json:"id"`
	ShohibulID    uuid.UUID `json:"shohibul_id"`
	PackageID     uint      `json:"package_id"`
	BookingDate   time.Time `json:"booking_date"`
	PaymentStatus string    `json:"payment_status"`
	TotalAmount   int64     `json:"total_amount"`
}

type CreateAnimalRequest struct {
	TagNumber  string  `json:"tag_number" binding:"required"`
	Type       string  `json:"type" binding:"required"`
	Weight     float64 `json:"weight" binding:"required,gt=0"`
	Status     string  `json:"status" binding:"required"`
	VendorInfo string  `json:"vendor_info"`
}

type AnimalResponse struct {
	ID         uuid.UUID `json:"id"`
	TagNumber  string    `json:"tag_number"`
	Type       string    `json:"type"`
	Weight     float64   `json:"weight"`
	Status     string    `json:"status"`
	VendorInfo string    `json:"vendor_info"`
}
