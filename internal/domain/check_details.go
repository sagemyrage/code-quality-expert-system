package domain

type CheckDetails struct {
	Check           Check
	Metrics         CheckMetrics
	Recommendations []CheckRecommendation
}
