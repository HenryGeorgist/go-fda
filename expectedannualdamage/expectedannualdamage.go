package expectedannualdamage

import (
	"errors"

	"github.com/HenryGeorgist/go-fda/flowfrequencycurves"
	"github.com/HenryGeorgist/go-fda/ratingcurves"
	"github.com/HenryGeorgist/go-fda/stagedamagecurves"
	_gc "github.com/USACE/go-consequences/compute"
	"github.com/USACE/go-consequences/paireddata"
)

type Simulation struct {
	FlowFrequency *flowfrequencycurves.UnregulatedFlowFrequencyCurve
	RatingCurve   *ratingcurves.RatingCurve
	DamageCurve   *stagedamagecurves.StageDamageCurve
}

func (ead Simulation) Compute() (float64, error) {
	if ead.FlowFrequency == nil {
		return 0.0, errors.New("unregulated flow frequeny curve is not defined")
	}
	if ead.RatingCurve == nil {
		return 0.0, errors.New("rating curve is not defined")
	}
	if ead.DamageCurve == nil {
		return 0.0, errors.New("damage curve is not defined")
	}
	ff := ead.FlowFrequency.DeterministicSample()
	ffpd, ffok := ff.(paireddata.PairedData)
	if !ffok {
		return 0.0, errors.New("flow frequeny curve is not paired data")
	}
	rc := ead.RatingCurve.Sample(.5)
	rcpd, rcok := rc.(paireddata.PairedData)
	if !rcok {
		return 0.0, errors.New("rating curve is not paired data")
	}
	dc := ead.DamageCurve.Sample(.5)
	dcpd, dcok := dc.(paireddata.PairedData)
	if !dcok {
		return 0.0, errors.New("stage damage curve is not paired data")
	}
	//compose frequency flow with flow stage
	fspd := rcpd.Compose(ffpd)
	//compose frequency stage with stage damage
	fdpd := dcpd.Compose(fspd)
	//integrate frequency damage
	result := _gc.ComputeEAD(fdpd.Xvals, fdpd.Yvals)
	return result, nil
}
