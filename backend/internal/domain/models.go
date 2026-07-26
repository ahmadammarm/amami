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
	JamaahProfile *Jamaah   `gorm:"foreignKey:UserID"`
	Status       string     `gorm:"size:50;default:'ACTIVE'"`
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
	TransactionID *uuid.UUID     `gorm:"type:uuid;index"`
	Transaction   *Transaction   `gorm:"constraint:OnDelete:RESTRICT;"`
	MuzakkiID     uuid.UUID      `gorm:"type:uuid;index;not null"`
	Muzakki       Jamaah         `gorm:"foreignKey:MuzakkiID;constraint:OnDelete:RESTRICT;"`
	ZakatType     string         `gorm:"type:varchar(30);not null"`
	AmountOrQty   float64        `gorm:"column:amount;type:decimal(12,2);not null;default:0"`
	Unit          string         `gorm:"type:varchar(20);default:'IDR'"` // IDR, KG, LITER
	Description   string         `gorm:"type:text"`
	Metadata      datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt     time.Time      `gorm:"default:now()"`
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
	TransactionID *uuid.UUID  `gorm:"type:uuid;index"`
	Transaction   *Transaction `gorm:"constraint:OnDelete:RESTRICT;"`
	MustahikID    uuid.UUID   `gorm:"type:uuid;index;not null"`
	Mustahik      Jamaah      `gorm:"foreignKey:MustahikID;constraint:OnDelete:RESTRICT;"`
	AmountOrQty   float64     `gorm:"column:amount;type:decimal(12,2);not null;default:0"`
	Unit          string      `gorm:"type:varchar(20);default:'IDR'"`
	Description   string      `gorm:"type:text"`
	DistributedAt time.Time   `gorm:"default:now()"`
}

// --- 5. Kurban Module ---

type QurbanPackage struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	Name           string `gorm:"type:varchar(100);not null" json:"name"`
	Type           string `gorm:"type:varchar(20)" json:"type"`
	Price          int64  `gorm:"type:bigint;not null" json:"price"`
	YearHijri      int    `gorm:"index;not null" json:"year_hijri"`
	StockTotal     int    `json:"stock_total"`
	StockRemaining int    `gorm:"check:stock_remaining >= 0" json:"stock_remaining"`
}

type QurbanBooking struct {
	ID            uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ShohibulID    uuid.UUID      `gorm:"type:uuid;index;not null" json:"shohibul_id"`
	Shohibul      Jamaah         `gorm:"foreignKey:ShohibulID;constraint:OnDelete:RESTRICT;" json:"shohibul"`
	PackageID     uint           `gorm:"index;not null" json:"package_id"`
	Package       QurbanPackage  `gorm:"constraint:OnDelete:RESTRICT;" json:"package"`
	BookingDate   time.Time      `gorm:"default:now()" json:"booking_date"`
	PaymentStatus string         `gorm:"type:varchar(30)" json:"payment_status"`
	TotalAmount   int64          `gorm:"type:bigint" json:"total_amount"`
	Metadata      datatypes.JSON `gorm:"type:jsonb" json:"metadata"`
}

type QurbanAnimal struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TagNumber  string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"tag_number"`
	Type       string    `gorm:"type:varchar(20)" json:"type"`
	Weight     float64   `gorm:"type:decimal(6,2)" json:"weight"`
	Status     string    `gorm:"type:varchar(30);index" json:"status"`
	VendorInfo string    `gorm:"type:text" json:"vendor_info"`
}

// --- 6. Inventory Module ---

type Asset struct {
	ID            uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name          string     `gorm:"type:varchar(255);index;not null" json:"name"`
	SKU           string     `gorm:"type:varchar(50);uniqueIndex" json:"sku"`
	PurchaseDate  *time.Time `json:"purchase_date"`
	PurchasePrice int64      `gorm:"type:bigint" json:"purchase_price"`
	CurrentStatus string     `gorm:"type:varchar(30)" json:"current_status"`
	Location      string     `gorm:"type:varchar(100)" json:"location"`
}

type AssetLoan struct {
	ID             uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	AssetID        uuid.UUID  `gorm:"type:uuid;index;not null" json:"asset_id"`
	Asset          Asset      `gorm:"constraint:OnDelete:RESTRICT;" json:"asset"`
	JamaahID       uuid.UUID  `gorm:"type:uuid;index;not null" json:"jamaah_id"`
	Jamaah         Jamaah     `gorm:"constraint:OnDelete:RESTRICT;" json:"jamaah"`
	LoanDate       time.Time  `gorm:"not null" json:"loan_date"`
	DueDate        *time.Time `json:"due_date"`
	ReturnDate     *time.Time `gorm:"index" json:"return_date"`
	ConditionNotes string     `gorm:"type:text" json:"condition_notes"`
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
	ID          uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Title       string     `gorm:"type:varchar(200);not null" json:"title"`
	Description string     `gorm:"type:text" json:"description"`
	Day         string     `gorm:"type:varchar(20);not null" json:"day"`
	Time        string     `gorm:"type:varchar(20);not null" json:"time"`
	Location    string     `gorm:"type:varchar(100)" json:"location"`
	Status      string     `gorm:"type:varchar(30);index" json:"status"`
	CreatedByID *uuid.UUID `gorm:"type:uuid" json:"created_by_id"`
	Creator     *User      `gorm:"foreignKey:CreatedByID;constraint:OnDelete:SET NULL;" json:"creator"`
}

type AgendaDocumentation struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	AgendaID   uuid.UUID `gorm:"type:uuid;index;not null" json:"agenda_id"`
	Agenda     Agenda    `gorm:"constraint:OnDelete:CASCADE;" json:"agenda"`
	FileURL    string    `gorm:"type:text;not null" json:"file_url"`
	FileType   string    `gorm:"type:varchar(50)" json:"file_type"`
	UploadedBy *uuid.UUID `gorm:"type:uuid" json:"uploaded_by"`
	Uploader   *User     `gorm:"foreignKey:UploadedBy;constraint:OnDelete:SET NULL;" json:"uploader"`
	UploadedAt time.Time `gorm:"default:now()" json:"uploaded_at"`
}

type AuditLog struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;index;not null"`
	User      User      `gorm:"constraint:OnDelete:RESTRICT;"`
	Action    string    `gorm:"type:varchar(50);not null"`
	Entity    string    `gorm:"type:varchar(50);not null"`
	EntityID  string    `gorm:"type:varchar(100)"`
	Metadata  string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"index;default:now()"`
}
