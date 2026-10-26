package domain

// String-backed enums persisted as varchar columns. Request DTOs validate the
// same values with `oneof` binding tags; keep both in sync when adding a value.

type TransactionType string

const (
	TransactionCredit TransactionType = "CREDIT"
	TransactionDebit  TransactionType = "DEBIT"
)

// Opposite returns the type that cancels this one out in a reversal entry.
func (t TransactionType) Opposite() TransactionType {
	if t == TransactionCredit {
		return TransactionDebit
	}
	return TransactionCredit
}

const (
	CategoryReversal = "REVERSAL"
	CategoryZakatIn  = "Penerimaan Zakat"
	CategoryZakatOut = "Penyaluran Zakat"
)

type UserStatus string

const (
	UserActive                UserStatus = "ACTIVE"
	UserSuspended             UserStatus = "SUSPENDED"
	UserPendingPasswordChange UserStatus = "PENDING_PASSWORD_CHANGE"
)

type AssetStatus string

const (
	AssetGood        AssetStatus = "GOOD"
	AssetLoaned      AssetStatus = "LOANED"
	AssetNeedsRepair AssetStatus = "NEEDS_REPAIR"
	AssetRepair      AssetStatus = "REPAIR"
	AssetBroken      AssetStatus = "BROKEN"
)

type AgendaStatus string

const (
	AgendaScheduled AgendaStatus = "SCHEDULED"
	AgendaOngoing   AgendaStatus = "ONGOING"
	AgendaCompleted AgendaStatus = "COMPLETED"
	AgendaCancelled AgendaStatus = "CANCELLED"
)

type ZakatUnit string

const (
	ZakatUnitIDR   ZakatUnit = "IDR"
	ZakatUnitKG    ZakatUnit = "KG"
	ZakatUnitLiter ZakatUnit = "LITER"
)

type PaymentStatus string

const (
	PaymentPending PaymentStatus = "PENDING"
	PaymentPaid    PaymentStatus = "PAID"
)

type AnimalStatus string

const (
	AnimalAvailable   AnimalStatus = "AVAILABLE"
	AnimalAssigned    AnimalStatus = "ASSIGNED"
	AnimalSlaughtered AnimalStatus = "SLAUGHTERED"
)

type HouseStatus string

const (
	HouseOwned  HouseStatus = "OWNED"
	HouseFamily HouseStatus = "FAMILY"
	HouseRented HouseStatus = "RENTED"
	HouseNone   HouseStatus = "NONE"
)
