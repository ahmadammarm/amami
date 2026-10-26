package domain

import (
	"slices"
	"testing"
)

// expectedGrants is docs/RBAC_VIEW_MATRIX.md section 4 written out by hand, so a
// change to RolePermissions that drifts from the documented matrix fails here.
var expectedGrants = map[PermissionCode][]string{
	PermUserManage:      {RoleSuperAdmin},
	PermLedgerRead:      {RoleSuperAdmin, RoleBendahara, RoleTakmir, RoleSekretaris},
	PermLedgerWrite:     {RoleBendahara},
	PermZakatRead:       {RoleSuperAdmin, RoleBendahara, RoleTakmir, RoleSekretaris},
	PermZakatManage:     {RoleSuperAdmin, RoleBendahara, RoleTakmir},
	PermJamaahManage:    {RoleSuperAdmin, RoleBendahara, RoleTakmir, RoleSekretaris},
	PermMustahikSurvey:  {RoleTakmir, RoleSekretaris},
	PermQurbanRead:      {RoleSuperAdmin, RoleBendahara, RoleTakmir, RoleSekretaris},
	PermQurbanManage:    {RoleTakmir, RoleSekretaris},
	PermInventoryRead:   {RoleSuperAdmin, RoleTakmir, RoleSekretaris, RoleJamaah},
	PermInventoryManage: {RoleTakmir, RoleSekretaris},
	PermAgendaRead:      {RoleSuperAdmin, RoleBendahara, RoleTakmir, RoleSekretaris, RoleJamaah},
	PermAgendaManage:    {RoleSuperAdmin, RoleTakmir, RoleSekretaris},
	PermSettingsManage:  {RoleSuperAdmin, RoleSekretaris},
}

func TestRolePermissionsMatchMatrix(t *testing.T) {
	actual := map[PermissionCode][]string{}
	for role, codes := range RolePermissions {
		seen := map[PermissionCode]bool{}
		for _, code := range codes {
			if seen[code] {
				t.Errorf("%s is granted %s twice", role, code)
			}
			seen[code] = true
			actual[code] = append(actual[code], role)
		}
	}

	for code, roles := range expectedGrants {
		got := slices.Sorted(slices.Values(actual[code]))
		want := slices.Sorted(slices.Values(roles))
		if !slices.Equal(got, want) {
			t.Errorf("%s: granted to %v, matrix says %v", code, got, want)
		}
	}
	for code := range actual {
		if _, ok := expectedGrants[code]; !ok {
			t.Errorf("%s is granted but missing from the matrix", code)
		}
	}
}

func TestEveryPermissionHasDescription(t *testing.T) {
	for code := range expectedGrants {
		if PermissionDescriptions[code] == "" {
			t.Errorf("%s has no description", code)
		}
	}
	if len(PermissionDescriptions) != len(expectedGrants) {
		t.Errorf("descriptions list %d codes, matrix lists %d", len(PermissionDescriptions), len(expectedGrants))
	}
}
