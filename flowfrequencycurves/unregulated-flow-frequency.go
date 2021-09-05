package flowfrequencycurves

import (
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
	//TODO: use random value to seed a bootstrap on the distribution to create a bootstrap
	for i := 0; i < uffc.Ordinates; i++ {
		probs[i] = float64(float64(i)+0.5) / float64(uffc.Ordinates)
		flows[i] = uffc.Distribution.InvCDF(probs[i])
	}
	return paireddata.PairedData{Xvals: flows, Yvals: probs}
}
