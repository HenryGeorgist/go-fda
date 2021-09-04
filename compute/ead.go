package compute

import (
	"github.com/HenryGeorgist/go-fda/flowfrequencycurves"
	"github.com/HenryGeorgist/go-fda/ratingcurves"
	"github.com/HenryGeorgist/go-fda/stagedamagecurves"
)

type ExpectedAnnualDamage struct {
	FlowFrequency flowfrequencycurves.UnregulatedFlowFrequencyCurve
	RatingCurve   ratingcurves.RatingCurve
	DamageCurve   stagedamagecurves.StageDamageCurve
}

func (ead ExpectedAnnualDamage) Compute() float64 {
	/*
		ff := ead.FlowFrequency.Sample(.5)
		rc := ead.RatingCurve.Sample(.5)
		dc := ead.DamageCurve.Sample(.5)
	*/
	return .5
}
