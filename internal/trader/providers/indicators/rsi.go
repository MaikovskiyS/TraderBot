package indicators

func (i *indicators) Rsi(values []float64, period int) []float64 {
	rsiValues := make([]float64, len(values))
	gain := 0.0
	loss := 0.0

	for i := 1; i < len(values); i++ {
		change := values[i] - values[i-1]
		if change > 0 {
			gain += change
		} else {
			loss -= change
		}

		if i < period {
			continue
		} else if i == period {
			rsiValues[i] = 100 - (100 / (1 + (gain/float64(period))/(loss/float64(period))))
		} else {
			prevRSI := rsiValues[i-1]
			rs := (prevRSI*13 + gain) / (prevRSI*13 + loss)
			rsiValues[i] = 100 - 100/(1+rs)
		}
	}

	return rsiValues
}
