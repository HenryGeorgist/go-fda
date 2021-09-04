package stagefrequencycurves

import (
	"github.com/USACE/go-consequences/paireddata"
)

type StageFrequencyCurve struct {
	Curve paireddata.UncertaintyPairedData
}

func (s StageFrequencyCurve) Sample(randomValue float64) paireddata.ValueSampler {
	return s.Curve.SampleValueSampler(randomValue)
}
