package domain

type PermissionCode string

const (
	PermUserManage      PermissionCode = "user:manage"
	PermLedgerRead      PermissionCode = "ledger:read"
	PermLedgerWrite     PermissionCode = "ledger:write"
	PermZakatRead       PermissionCode = "zakat:read"
	PermZakatManage     PermissionCode = "zakat:manage"
	PermJamaahManage    PermissionCode = "jamaah:manage"
	PermMustahikSurvey  PermissionCode = "mustahik:survey"
	PermQurbanRead      PermissionCode = "qurban:read"
	PermQurbanManage    PermissionCode = "qurban:manage"
	PermInventoryRead   PermissionCode = "inventory:read"
	PermInventoryManage PermissionCode = "inventory:manage"
	PermAgendaRead      PermissionCode = "agenda:read"
	PermAgendaManage    PermissionCode = "agenda:manage"
	PermSettingsManage  PermissionCode = "settings:manage"
)

const (
	RoleSuperAdmin = "SUPER_ADMIN"
	RoleBendahara  = "BENDAHARA"
	RoleTakmir     = "TAKMIR"
	RoleSekretaris = "SEKRETARIS"
	RoleJamaah     = "JAMAAH"
)

// PermissionDescriptions lists every permission the system knows about.
var PermissionDescriptions = map[PermissionCode]string{
	PermUserManage:      "Manage users and roles",
	PermLedgerRead:      "Read the ledger and financial reports",
	PermLedgerWrite:     "Create funds, transactions and reversal entries",
	PermZakatRead:       "Read zakat collections and distributions",
	PermZakatManage:     "Record zakat collection and distribution",
	PermJamaahManage:    "Manage the community registry",
	PermMustahikSurvey:  "Submit mustahik surveys that set the poverty score",
	PermQurbanRead:      "Read qurban packages, bookings and animals",
	PermQurbanManage:    "Manage qurban packages, bookings and animals",
	PermInventoryRead:   "Read the asset registry",
	PermInventoryManage: "Manage assets and loans",
	PermAgendaRead:      "Read agendas and their documentation",
	PermAgendaManage:    "Manage agendas and upload documentation",
	PermSettingsManage:  "Manage mosque profile and system settings",
}

// RolePermissions is the authoritative grant table, mirroring docs/RBAC_VIEW_MATRIX.md.
// The seeder syncs role_permissions to exactly this map.
var RolePermissions = map[string][]PermissionCode{
	RoleSuperAdmin: {
		PermUserManage, PermLedgerRead, PermZakatRead, PermZakatManage, PermJamaahManage,
		PermQurbanRead, PermInventoryRead, PermAgendaRead, PermAgendaManage, PermSettingsManage,
	},
	RoleBendahara: {
		PermLedgerRead, PermLedgerWrite, PermZakatRead, PermZakatManage, PermJamaahManage,
		PermQurbanRead, PermAgendaRead,
	},
	RoleTakmir: {
		PermLedgerRead, PermZakatRead, PermZakatManage, PermJamaahManage, PermMustahikSurvey,
		PermQurbanRead, PermQurbanManage, PermInventoryRead, PermInventoryManage,
		PermAgendaRead, PermAgendaManage,
	},
	RoleSekretaris: {
		PermLedgerRead, PermZakatRead, PermJamaahManage, PermMustahikSurvey,
		PermQurbanRead, PermQurbanManage, PermInventoryRead, PermInventoryManage,
		PermAgendaRead, PermAgendaManage, PermSettingsManage,
	},
	RoleJamaah: {
		PermInventoryRead, PermAgendaRead,
	},
}
