package jamaah

import (
	"math"

	"github.com/ahmadammarm/amami/backend/internal/domain"
)

// Poverty score weights (sum = 100). A higher score means a stronger claim to aid.
const (
	weightIncome     = 50.0
	weightDependents = 20.0
	weightHouse      = 20.0
	weightDisability = 10.0

	// IncomePerCapitaThreshold is the monthly income per household member (IDR)
	// at or above which the income component contributes nothing. It is roughly
	// twice the BPS national poverty line.
	IncomePerCapitaThreshold = 1_200_000.0

	// maxScoredDependents caps the dependents component; more dependents than
	// this do not raise the score further.
	maxScoredDependents = 5
)

var houseStatusPoints = map[domain.HouseStatus]float64{
	domain.HouseNone:   1.0,
	domain.HouseRented: 0.75,
	domain.HouseFamily: 0.5,
	domain.HouseOwned:  0,
}

type MustahikInput struct {
	TotalIncome     int64
	DependentsCount int
	HouseStatus     domain.HouseStatus
	HasDisability   bool
}

// CalculateMustahikScore turns survey answers into a 0-100 poverty score,
// rounded to two decimals.
func CalculateMustahikScore(in MustahikInput) float64 {
	dependents := max(in.DependentsCount, 0)
	income := math.Max(float64(in.TotalIncome), 0)

	perCapita := income / float64(dependents+1)
	incomeScore := weightIncome * clamp01(1-perCapita/IncomePerCapitaThreshold)

	dependentsScore := weightDependents * float64(min(dependents, maxScoredDependents)) / maxScoredDependents

	houseScore := weightHouse * houseStatusPoints[in.HouseStatus]

	disabilityScore := 0.0
	if in.HasDisability {
		disabilityScore = weightDisability
	}

	total := incomeScore + dependentsScore + houseScore + disabilityScore
	return math.Round(total*100) / 100
}

func clamp01(v float64) float64 {
	return math.Min(math.Max(v, 0), 1)
}
