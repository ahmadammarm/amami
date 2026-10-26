package qurban

import (
	"time"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/google/uuid"
)

type PackageRequest struct {
	Name       string `json:"name" binding:"required,max=100"`
	Type       string `json:"type" binding:"required,oneof=KAMBING DOMBA SAPI SAPI_KOLEKTIF"`
	Price      int64  `json:"price" binding:"required,gt=0"`
	YearHijri  int    `json:"year_hijri" binding:"required,gte=1400,lte=1600"`
	StockTotal int    `json:"stock_total" binding:"required,gt=0"`
}

type CreateBookingRequest struct {
	ShohibulID    uuid.UUID            `json:"shohibul_id" binding:"required"`
	PackageID     uint                 `json:"package_id" binding:"required"`
	PaymentStatus domain.PaymentStatus `json:"payment_status" binding:"required,oneof=PENDING PAID"`
}

type AnimalRequest struct {
	TagNumber  string              `json:"tag_number" binding:"required,max=50"`
	Type       string              `json:"type" binding:"required,oneof=KAMBING DOMBA SAPI"`
	Weight     float64             `json:"weight" binding:"required,gt=0,lt=10000"`
	Status     domain.AnimalStatus `json:"status" binding:"required,oneof=AVAILABLE ASSIGNED SLAUGHTERED"`
	VendorInfo string              `json:"vendor_info" binding:"max=1000"`
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

type BookingResponse struct {
	ID            uuid.UUID            `json:"id"`
	ShohibulID    uuid.UUID            `json:"shohibul_id"`
	ShohibulName  string               `json:"shohibul_name"`
	PackageID     uint                 `json:"package_id"`
	PackageName   string               `json:"package_name"`
	BookingDate   time.Time            `json:"booking_date"`
	PaymentStatus domain.PaymentStatus `json:"payment_status"`
	TotalAmount   int64                `json:"total_amount"`
}

type AnimalResponse struct {
	ID         uuid.UUID           `json:"id"`
	TagNumber  string              `json:"tag_number"`
	Type       string              `json:"type"`
	Weight     float64             `json:"weight"`
	Status     domain.AnimalStatus `json:"status"`
	VendorInfo string              `json:"vendor_info"`
}

func ToPackageResponse(p *domain.QurbanPackage) PackageResponse {
	return PackageResponse{
		ID:             p.ID,
		Name:           p.Name,
		Type:           p.Type,
		Price:          p.Price,
		YearHijri:      p.YearHijri,
		StockTotal:     p.StockTotal,
		StockRemaining: p.StockRemaining,
	}
}

// ToBookingResponse expects Shohibul and Package to be loaded.
func ToBookingResponse(b *domain.QurbanBooking) BookingResponse {
	return BookingResponse{
		ID:            b.ID,
		ShohibulID:    b.ShohibulID,
		ShohibulName:  b.Shohibul.FullName,
		PackageID:     b.PackageID,
		PackageName:   b.Package.Name,
		BookingDate:   b.BookingDate,
		PaymentStatus: b.PaymentStatus,
		TotalAmount:   b.TotalAmount,
	}
}

func ToAnimalResponse(a *domain.QurbanAnimal) AnimalResponse {
	return AnimalResponse{
		ID:         a.ID,
		TagNumber:  a.TagNumber,
		Type:       a.Type,
		Weight:     a.Weight,
		Status:     a.Status,
		VendorInfo: a.VendorInfo,
	}
}
