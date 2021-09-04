package compute

import (
	"fmt"

	"github.com/HenryGeorgist/go-fda/flowfrequencycurves"
	"github.com/HenryGeorgist/go-fda/ratingcurves"
	"github.com/HenryGeorgist/go-fda/stagedamagecurves"
	_gcc "github.com/USACE/go-consequences/compute"
	"github.com/USACE/go-consequences/paireddata"
)

type ExpectedAnnualDamage struct {
	FlowFrequency flowfrequencycurves.UnregulatedFlowFrequencyCurve
	RatingCurve   ratingcurves.RatingCurve
	DamageCurve   stagedamagecurves.StageDamageCurve
}

func (ead ExpectedAnnualDamage) Compute() float64 {
	ff := ead.FlowFrequency.Sample(.5)
	ffpd, ffok := ff.(paireddata.PairedData)
	if !ffok {
		panic("frequency curve is not paired data.")
	}
	rc := ead.RatingCurve.Sample(.5)
	rcpd, rcok := rc.(paireddata.PairedData)
	if !rcok {
		panic("rating curve is not paired data.")
	}
	dc := ead.DamageCurve.Sample(.5)
	dcpd, dcok := dc.(paireddata.PairedData)
	if !dcok {
		panic("stage damage curve is not paired data.")
	}
	fmt.Println(dcpd)
	//compose frequency flow with flow stage
	fspd := ffpd.Compose(rcpd)
	//compose frequency stage with stage damage
	fdpd := fspd.Compose(dcpd)
	//integrate frequency damage
	result := _gcc.ComputeSpecialEAD(ffpd.Xvals, ffpd.Yvals)
	fmt.Println(fdpd)
	return result
}
