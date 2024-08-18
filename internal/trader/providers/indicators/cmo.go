package indicators

import "math"

// Функция для расчета CMO
func (i *indicators) Cmo(open []float64, close []float64) []float64 {
	m1 := make([]float64, len(open))
	m2 := make([]float64, len(open))
	cmoValues := make([]float64, len(open))

	for i := 1; i < len(open); i++ {
		src1 := (open[i] - open[i-1]) / open[i-1]
		src2 := (close[i] - close[i-1]) / close[i-1]
		m1[i] = math.Max(src1, src2)
		m2[i] = math.Min(src1, src2)

		sm1 := sum(m1)
		sm2 := sum(m2)
		cmoValues[i] = 100 * (sm1 - sm2) / (sm1 + sm2)
	}

	return cmoValues
}

func sum(slice []float64) float64 {
	total := 0.0
	for _, value := range slice {
		total += value
	}
	return total
}
