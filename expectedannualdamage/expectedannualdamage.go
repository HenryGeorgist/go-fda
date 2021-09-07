package expectedannualdamage

import (
	"errors"

	"github.com/HenryGeorgist/go-fda/flowfrequencycurves"
	"github.com/HenryGeorgist/go-fda/flowtransformfunction"
	"github.com/HenryGeorgist/go-fda/interiorexteriorcurve"
	"github.com/HenryGeorgist/go-fda/leveefragilitycurve"
	"github.com/HenryGeorgist/go-fda/ratingcurves"
	"github.com/HenryGeorgist/go-fda/stagedamagecurves"
	_gc "github.com/USACE/go-consequences/compute"
	"github.com/USACE/go-consequences/paireddata"
)

type Simulation struct {
	FlowFrequency         *flowfrequencycurves.UnregulatedFlowFrequencyCurve
	FlowTransform         *flowtransformfunction.FlowTransformFunctionCurve
	RatingCurve           *ratingcurves.RatingCurve
	InteriorExteriorCurve *interiorexteriorcurve.InteriorExteriorCurve
	LeveeFragilityCurve   *leveefragilitycurve.LeveeFragilityCurve
	DamageCurve           *stagedamagecurves.StageDamageCurve
}

func (ead Simulation) Compute() (float64, error) {
	if ead.FlowFrequency == nil {
		return 0.0, errors.New("unregulated flow frequeny curve is not defined")
	}
	hasFlowTransform := false
	if ead.FlowTransform != nil {
		hasFlowTransform = true
	}
	if ead.RatingCurve == nil {
		return 0.0, errors.New("rating curve is not defined")
	}
	hasFragilityCurve := false
	if ead.LeveeFragilityCurve != nil {
		hasFragilityCurve = true
	}
	hasInteriorExteriorCurve := false
	if ead.InteriorExteriorCurve != nil {
		hasInteriorExteriorCurve = true
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
	if hasFlowTransform {
		ft := ead.FlowTransform.Sample(.5)
		ftpd, ftok := ft.(paireddata.PairedData)
		if !ftok {
			return 0.0, errors.New("flow transform is not paired data")
		}
		ffpd = ftpd.Compose(ffpd)
	}
	//compose frequency flow with flow stage
	fspd := rcpd.Compose(ffpd)
	//if an interior exterior relationship exists compose it with the frequency stage
	if hasInteriorExteriorCurve {
		ie := ead.InteriorExteriorCurve.Sample(.5)
		iepd, ieok := ie.(paireddata.PairedData)
		if !ieok {
			return 0.0, errors.New("interior exterior is not paired data")
		}
		fspd = iepd.Compose(fspd)
	}
	//if a levee exists, modify stage damage curve
	if hasFragilityCurve {
		lfc := ead.LeveeFragilityCurve.Sample(.5)
		lfcpd, lfcok := lfc.(paireddata.PairedData)
		if !lfcok {
			return 0.0, errors.New("levee fragility curve is not paired data")
		}
		dcpd = leveefragilitycurve.Multiply(dcpd, lfcpd) //this is not quite right. need something other than compose. (multiply?)
	}
	//compose frequency stage with stage damage
	fdpd := dcpd.Compose(fspd)
	//integrate frequency damage
	result := _gc.ComputeEAD(fdpd.Xvals, fdpd.Yvals)
	return result, nil
}
