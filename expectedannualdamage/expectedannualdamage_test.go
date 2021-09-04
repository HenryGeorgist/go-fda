package expectedannualdamage

import (
	"fmt"
	"testing"

	flowfrequencycurves "github.com/HenryGeorgist/go-fda/flow-frequency-curves"
	ratingcurves "github.com/HenryGeorgist/go-fda/rating-curves"
	stagedamagecurves "github.com/HenryGeorgist/go-fda/stage-damage-curves"
	"github.com/HydrologicEngineeringCenter/go-statistics/statistics"
	"github.com/USACE/go-consequences/paireddata"
)

func Test_EAD(t *testing.T) {
	n := statistics.LogPearsonIIIDistribution{Mean: 3.368, StandardDeviation: .246, Skew: .668}
	//create a flow frequency curve based on lp3
	ff := flowfrequencycurves.UnregulatedFlowFrequencyCurve{Distribution: n, Ordinates: 100}
	//create some basic data that will fill future structs
	flows := []float64{14000.0, 13000.0, 12000.0, 11000.0, 7000.0, 5000.0, 4000.0, 3000.0, 1500.0, 1000.0, 900.0, 850.0}
	stages := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0}
	stageDists := make([]statistics.ContinuousDistribution, len(stages))
	for i, s := range stages {
		dist, _ := statistics.InitDeterministic(s)
		stageDists[i] = dist
	}
	//create a rating curve based on fake data
	rc := ratingcurves.RatingCurve{Curve: paireddata.UncertaintyPairedData{Xvals: flows, Yvals: stageDists}}

	//create some basic data for damages
	damages := []float64{10.0, 20.0, 30.0, 40.0, 50.0, 60.0, 70.0, 80.0, 90.0, 100.0, 110.0, 120.0}
	damageDists := make([]statistics.ContinuousDistribution, len(damages))
	for i, d := range damages {
		dist, _ := statistics.InitDeterministic(d)
		damageDists[i] = dist
	}
	dc := stagedamagecurves.StageDamageCurve{Curve: paireddata.UncertaintyPairedData{Xvals: stages, Yvals: damageDists}}

	sim := Simulation{FlowFrequency: ff, RatingCurve: rc, DamageCurve: dc}
	ead := sim.Compute()
	fmt.Println(ead)
}
