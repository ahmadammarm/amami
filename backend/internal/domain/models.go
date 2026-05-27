package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// --- 1. Core & Auth ---

type Role struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"type:varchar(50);uniqueIndex;not null"`
	Description string `gorm:"type:text"`
}

type Permission struct {
	ID          uint   `gorm:"primaryKey"`
	Code        string `gorm:"type:varchar(50);uniqueIndex;not null"`
	Description string `gorm:"type:text"`
}

type RolePermission struct {
	RoleID       uint       `gorm:"primaryKey"`
	PermissionID uint       `gorm:"primaryKey"`
	Role         Role       `gorm:"constraint:OnDelete:CASCADE;"`
	Permission   Permission `gorm:"constraint:OnDelete:CASCADE;"`
}

type User struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username     string     `gorm:"type:varchar(50);uniqueIndex;not null"`
	Email        string     `gorm:"type:varchar(100);uniqueIndex;not null"`
	PasswordHash string     `gorm:"type:text;not null"`
	RoleID       uint       `gorm:"index"`
	Role         Role       `gorm:"constraint:OnDelete:RESTRICT;"`
	Status       string     `gorm:"type:varchar(20);default:'ACTIVE'"`
	LastLoginAt  *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// --- 2. Jamaah & Profiles ---

type Jamaah struct {
	ID            uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID        *uuid.UUID `gorm:"type:uuid;uniqueIndex"`
	User          *User      `gorm:"constraint:OnDelete:SET NULL;"`
	FullName      string     `gorm:"type:varchar(100);not null;index:idx_jamaah_name,type:gin,expression:full_name gin_trgm_ops"` // GIN Index for fast search
	Phone         string     `gorm:"type:varchar(20)"`
	Address       string     `gorm:"type:text"`
	IsMustahik    bool       `gorm:"default:false;index"`
	MustahikScore float64    `gorm:"type:decimal(5,2);default:0.00"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// --- 3. Finance (The Immutable Ledger) ---

type Fund struct {
	ID             uint   `gorm:"primaryKey"`
	Name           string `gorm:"type:varchar(100);not null"`
	Code           string `gorm:"type:varchar(20);uniqueIndex;not null"`
	CurrentBalance int64  `gorm:"type:bigint;default:0"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Transaction struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	FundID      uint           `gorm:"index:idx_fund_created,priority:1;not null"` // Composite Index Part 1
	Fund        Fund           `gorm:"constraint:OnDelete:RESTRICT;"`
	Type        string         `gorm:"type:varchar(10);not null"` // DEBIT/CREDIT
	Amount      int64          `gorm:"type:bigint;not null;check:amount > 0"`
	Category    string         `gorm:"type:varchar(30);not null;index"`
	ReferenceID *uuid.UUID     `gorm:"type:uuid;index"` // Polymorphic relation ID
	Metadata    datatypes.JSON `gorm:"type:jsonb"`
	CreatedBy   uuid.UUID      `gorm:"type:uuid;not null;index"`
	Creator     User           `gorm:"foreignKey:CreatedBy;constraint:OnDelete:RESTRICT;"`
	CreatedAt   time.Time      `gorm:"index:idx_fund_created,priority:2;index"` // Composite Index Part 2 & Standalone
}

// --- 4. Zakat Module ---

type ZakatDonation struct {
	ID            uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TransactionID uuid.UUID      `gorm:"type:uuid;uniqueIndex;not null"`
	Transaction   Transaction    `gorm:"constraint:OnDelete:RESTRICT;"`
	MuzakkiID     uuid.UUID      `gorm:"type:uuid;index;not null"`
	Muzakki       Jamaah         `gorm:"foreignKey:MuzakkiID;constraint:OnDelete:RESTRICT;"`
	ZakatType     string         `gorm:"type:varchar(30);not null"`
	Metadata      datatypes.JSON `gorm:"type:jsonb"`
}

type MustahikData struct {
	ID              uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	JamaahID        uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`
	Jamaah          Jamaah    `gorm:"constraint:OnDelete:RESTRICT;"`
	TotalIncome     int64     `gorm:"type:bigint"`
	DependentsCount int       `gorm:"type:int"`
	HouseStatus     string    `gorm:"type:varchar(50)"`
	CriteriaNotes   string    `gorm:"type:text"`
	UpdatedAt       time.Time
}

type ZakatDistribution struct {
	ID            uuid.UUID   `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TransactionID uuid.UUID   `gorm:"type:uuid;uniqueIndex;not null"`
	Transaction   Transaction `gorm:"constraint:OnDelete:RESTRICT;"`
	MustahikID    uuid.UUID   `gorm:"type:uuid;index;not null"`
	Mustahik      Jamaah      `gorm:"foreignKey:MustahikID;constraint:OnDelete:RESTRICT;"`
	Amount        int64       `gorm:"type:bigint;not null"`
	DistributedAt time.Time   `gorm:"default:now()"`
}

// --- 5. Kurban Module ---

type QurbanPackage struct {
	ID             uint   `gorm:"primaryKey"`
	Name           string `gorm:"type:varchar(100);not null"`
	Type           string `gorm:"type:varchar(20)"`
	Price          int64  `gorm:"type:bigint;not null"`
	YearHijri      int    `gorm:"index;not null"`
	StockTotal     int
	StockRemaining int `gorm:"check:stock_remaining >= 0"`
}

type QurbanBooking struct {
	ID            uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ShohibulID    uuid.UUID      `gorm:"type:uuid;index;not null"`
	Shohibul      Jamaah         `gorm:"foreignKey:ShohibulID;constraint:OnDelete:RESTRICT;"`
	PackageID     uint           `gorm:"index;not null"`
	Package       QurbanPackage  `gorm:"constraint:OnDelete:RESTRICT;"`
	BookingDate   time.Time      `gorm:"default:now()"`
	PaymentStatus string         `gorm:"type:varchar(20)"`
	TotalAmount   int64          `gorm:"type:bigint"`
	Metadata      datatypes.JSON `gorm:"type:jsonb"`
}

type QurbanAnimal struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TagNumber  string    `gorm:"type:varchar(50);uniqueIndex;not null"`
	Type       string    `gorm:"type:varchar(20)"`
	Weight     float64   `gorm:"type:decimal(6,2)"`
	Status     string    `gorm:"type:varchar(20);index"`
	VendorInfo string    `gorm:"type:text"`
}

// --- 6. Inventory Module ---

type Asset struct {
	ID            uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name          string     `gorm:"type:varchar(255);index;not null"`
	SKU           string     `gorm:"type:varchar(50);uniqueIndex"`
	PurchaseDate  *time.Time
	PurchasePrice int64      `gorm:"type:bigint"`
	CurrentStatus string     `gorm:"type:varchar(20)"`
	Location      string     `gorm:"type:varchar(100)"`
}

type AssetLoan struct {
	ID             uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	AssetID        uuid.UUID  `gorm:"type:uuid;index;not null"`
	Asset          Asset      `gorm:"constraint:OnDelete:RESTRICT;"`
	JamaahID       uuid.UUID  `gorm:"type:uuid;index;not null"`
	Jamaah         Jamaah     `gorm:"constraint:OnDelete:RESTRICT;"`
	LoanDate       time.Time  `gorm:"not null"`
	DueDate        *time.Time
	ReturnDate     *time.Time `gorm:"index"` // Indexed to quickly find active loans
	ConditionNotes string     `gorm:"type:text"`
}

// --- 7. System & Utility Module ---

type SystemSetting struct {
	ID       uint   `gorm:"primaryKey"`
	Key      string `gorm:"type:varchar(50);uniqueIndex;not null"`
	Value    string `gorm:"type:text;not null"`
	IsSecret bool   `gorm:"default:false"`
}

type Memo struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Title     string    `gorm:"type:varchar(200);not null"`
	Content   string    `gorm:"type:text;not null"`
	AuthorID  uuid.UUID `gorm:"type:uuid;not null"`
	Author    User      `gorm:"foreignKey:AuthorID;constraint:OnDelete:RESTRICT;"`
	IsPinned  bool      `gorm:"default:false"`
	CreatedAt time.Time `gorm:"index"`
}

// --- 8. Logistics & Event Module ---

type Agenda struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Title       string     `gorm:"type:varchar(200);not null"`
	Description string     `gorm:"type:text"`
	StartTime   time.Time  `gorm:"not null;index"`
	EndTime     time.Time  `gorm:"not null"`
	Location    string     `gorm:"type:varchar(100)"`
	Status      string     `gorm:"type:varchar(20);index"`
	CreatedByID *uuid.UUID `gorm:"type:uuid"`
	Creator     *User      `gorm:"foreignKey:CreatedByID;constraint:OnDelete:SET NULL;"`
}

type AgendaDocumentation struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	AgendaID   uuid.UUID `gorm:"type:uuid;index;not null"`
	Agenda     Agenda    `gorm:"constraint:OnDelete:CASCADE;"`
	FileURL    string    `gorm:"type:text;not null"`
	FileType   string    `gorm:"type:varchar(50)"`
	UploadedBy *uuid.UUID `gorm:"type:uuid"`
	Uploader   *User     `gorm:"foreignKey:UploadedBy;constraint:OnDelete:SET NULL;"`
	UploadedAt time.Time `gorm:"default:now()"`
}
