package jamaah

import (
	"testing"

	"github.com/ahmadammarm/amami/backend/internal/domain"
)

func TestCalculateMustahikScore(t *testing.T) {
	tests := []struct {
		name string
		in   MustahikInput
		want float64
	}{
		{
			name: "no income, homeless, disabled, no dependents",
			in:   MustahikInput{TotalIncome: 0, DependentsCount: 0, HouseStatus: domain.HouseNone, HasDisability: true},
			want: 80,
		},
		{
			name: "maximum score",
			in:   MustahikInput{TotalIncome: 0, DependentsCount: 5, HouseStatus: domain.HouseNone, HasDisability: true},
			want: 100,
		},
		{
			name: "income per capita at threshold contributes nothing",
			in:   MustahikInput{TotalIncome: 2_400_000, DependentsCount: 1, HouseStatus: domain.HouseOwned},
			want: 4,
		},
		{
			name: "half the threshold, renting",
			in:   MustahikInput{TotalIncome: 600_000, DependentsCount: 0, HouseStatus: domain.HouseRented},
			want: 40,
		},
		{
			name: "dependents are capped",
			in:   MustahikInput{TotalIncome: 100_000_000, DependentsCount: 12, HouseStatus: domain.HouseFamily},
			want: 30,
		},
		{
			name: "well-off owner scores zero",
			in:   MustahikInput{TotalIncome: 50_000_000, DependentsCount: 0, HouseStatus: domain.HouseOwned},
			want: 0,
		},
		{
			name: "negative inputs are treated as zero",
			in:   MustahikInput{TotalIncome: -1, DependentsCount: -3, HouseStatus: domain.HouseOwned},
			want: 50,
		},
		{
			name: "unknown house status adds nothing",
			in:   MustahikInput{TotalIncome: 0, DependentsCount: 0, HouseStatus: "CASTLE"},
			want: 50,
		},
		{
			name: "rounded to two decimals",
			in:   MustahikInput{TotalIncome: 1_000_000, DependentsCount: 2, HouseStatus: domain.HouseOwned},
			want: 44.11,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateMustahikScore(tt.in)
			if got != tt.want {
				t.Fatalf("CalculateMustahikScore(%+v) = %v, want %v", tt.in, got, tt.want)
			}
			if got < 0 || got > 100 {
				t.Fatalf("score %v is outside 0-100", got)
			}
		})
	}
}
