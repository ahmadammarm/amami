package qurban

import (
	"testing"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	qurbanDto "github.com/ahmadammarm/amami/backend/internal/dto/qurban"
	qurbanRepo "github.com/ahmadammarm/amami/backend/internal/repository/qurban"
	"github.com/ahmadammarm/amami/backend/internal/testutil/pgtest"
	"github.com/ahmadammarm/amami/backend/pkg/apperror"
)

// Quota edits interleaved with bookings must never lose a booking's decrement:
// at the end, remaining stock is exactly the total minus the bookings made.
func TestQuotaEditsDuringBookings(t *testing.T) {
	const bookings, edits = 30, 10

	db := pgtest.New(t)
	svc := NewQurbanService(qurbanRepo.NewQurbanRepository(db))
	shohibul := pgtest.CreateJamaah(t, db, false)

	base := qurbanDto.PackageRequest{Name: "Kambing", Type: "KAMBING", Price: 3_000_000, YearHijri: 1447, StockTotal: 40}
	pkg, err := svc.CreatePackage(base)
	if err != nil {
		t.Fatalf("create package: %v", err)
	}

	errs := pgtest.RunConcurrently(bookings+edits, func(i int) error {
		if i < bookings {
			_, err := svc.CreateBooking(qurbanDto.CreateBookingRequest{
				ShohibulID: shohibul.ID, PackageID: pkg.ID, PaymentStatus: domain.PaymentPending,
			})
			return err
		}
		edit := base
		edit.StockTotal = 35 + i%2*10 // alternate between 35 and 45
		_, err := svc.UpdatePackage(pkg.ID, edit)
		return err
	})
	for _, err := range errs {
		if err != nil && !apperror.IsKind(err, apperror.KindConflict) {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	var stored domain.QurbanPackage
	if err := db.First(&stored, pkg.ID).Error; err != nil {
		t.Fatalf("reload package: %v", err)
	}
	var booked int64
	if err := db.Model(&domain.QurbanBooking{}).Where("package_id = ?", pkg.ID).Count(&booked).Error; err != nil {
		t.Fatalf("count bookings: %v", err)
	}

	if want := stored.StockTotal - int(booked); stored.StockRemaining != want {
		t.Fatalf("stock remaining = %d, want total %d - booked %d = %d",
			stored.StockRemaining, stored.StockTotal, booked, want)
	}
	if successes := pgtest.CountNil(errs[:bookings]); int64(successes) != booked {
		t.Fatalf("%d bookings reported success but %d were stored", successes, booked)
	}
}
