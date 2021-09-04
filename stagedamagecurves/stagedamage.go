package stagedamagecurves

import (
	"github.com/USACE/go-consequences/paireddata"
)

type StageDamageCurve struct {
	Curve paireddata.UncertaintyPairedData
}

func (s StageDamageCurve) Sample(randomValue float64) paireddata.ValueSampler {
	return s.Curve.SampleValueSampler(randomValue)
}
