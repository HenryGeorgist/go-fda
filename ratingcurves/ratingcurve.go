package ratingcurves

import (
	"github.com/USACE/go-consequences/paireddata"
)

type RatingCurve struct {
	Curve paireddata.UncertaintyPairedData
}

func (s RatingCurve) Sample(randomValue float64) paireddata.ValueSampler {
	return s.Curve.SampleValueSampler(randomValue)
}
