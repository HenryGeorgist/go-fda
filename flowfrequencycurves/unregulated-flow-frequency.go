package flowfrequencycurves

import (
	"math/rand"

	"github.com/HydrologicEngineeringCenter/go-statistics/statistics"
	"github.com/USACE/go-consequences/paireddata"
)

type UnregulatedFlowFrequencyCurve struct {
	Distribution            statistics.ContinuousDistribution
	EquivalantYearsOfRecord int
	Ordinates               int
}

func (uffc UnregulatedFlowFrequencyCurve) Sample(randomValue float64) paireddata.ValueSampler {
	flows := make([]float64, uffc.Ordinates)
	probs := make([]float64, uffc.Ordinates)
	bootstraplpiii := bootstrap(uffc, randomValue)
	for i := 0; i < uffc.Ordinates; i++ {
		probs[i] = float64(float64(i)+0.5) / float64(uffc.Ordinates)
		flows[i] = bootstraplpiii.InvCDF(probs[i])
	}
	return paireddata.PairedData{Xvals: flows, Yvals: probs}
}
func (uffc UnregulatedFlowFrequencyCurve) DeterministicSample() paireddata.ValueSampler {
	flows := make([]float64, uffc.Ordinates)
	probs := make([]float64, uffc.Ordinates)
	for i := 0; i < uffc.Ordinates; i++ {
		probs[i] = float64(float64(i)+0.5) / float64(uffc.Ordinates)
		flows[i] = uffc.Distribution.InvCDF(probs[i])
	}
	return paireddata.PairedData{Xvals: flows, Yvals: probs}
}
func bootstrap(uffc UnregulatedFlowFrequencyCurve, randomValue float64) statistics.LogPearsonIIIDistribution {
	rngsrc := rand.NewSource(int64(randomValue))
	rng := rand.New(rngsrc)
	sample := make([]float64, uffc.EquivalantYearsOfRecord)
	for i := 0; i < uffc.EquivalantYearsOfRecord; i++ {
		sample[i] = uffc.Distribution.InvCDF(rng.Float64())
	}
	lpiii := statistics.LogPearsonIIIDistribution{}
	lpiii.Fit(sample)
	return lpiii
}
