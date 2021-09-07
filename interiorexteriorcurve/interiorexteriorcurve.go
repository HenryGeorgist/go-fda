package interiorexteriorcurve

import (
	"github.com/USACE/go-consequences/paireddata"
)

type InteriorExteriorCurve struct {
	Curve paireddata.UncertaintyPairedData
}

func (s InteriorExteriorCurve) Sample(randomValue float64) paireddata.ValueSampler {
	return s.Curve.SampleValueSampler(randomValue)
}
