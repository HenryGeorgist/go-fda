package expectedannualdamage

import (
	"fmt"
	"testing"

	"github.com/HenryGeorgist/go-fda/flowfrequencycurves"
	"github.com/HenryGeorgist/go-fda/flowtransformfunction"
	"github.com/HenryGeorgist/go-fda/leveefragilitycurve"
	"github.com/HenryGeorgist/go-fda/ratingcurves"
	"github.com/HenryGeorgist/go-fda/stagedamagecurves"
	"github.com/HydrologicEngineeringCenter/go-statistics/statistics"
	"github.com/USACE/go-consequences/paireddata"
)

func Test_EAD(t *testing.T) {
	n := statistics.LogPearsonIIIDistribution{Mean: 3.368, StandardDeviation: .246, Skew: .668}
	//create a flow frequency curve based on lp3
	ff := flowfrequencycurves.UnregulatedFlowFrequencyCurve{Distribution: n, Ordinates: 1000}
	//create some basic data that will fill future structs
	flows := []float64{14000.0, 13000.0, 12000.0, 11000.0, 7000.0, 5000.0, 4000.0, 3000.0, 1500.0, 1000.0, 900.0, 850.0}
	for i, j := 0, len(flows)-1; i < j; i, j = i+1, j-1 {
		flows[i], flows[j] = flows[j], flows[i]
	}
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

	sim := Simulation{FlowFrequency: &ff, RatingCurve: &rc, DamageCurve: &dc}
	ead, err := sim.Compute()
	if err != nil {
		fmt.Println(err)
	}
	expected := 47.6930981139668
	if ead != expected {
		t.Errorf("expected %f, got %f", expected, ead)
	}
}
func Test_EAD_FlowTransform(t *testing.T) {
	n := statistics.LogPearsonIIIDistribution{Mean: 3.368, StandardDeviation: .246, Skew: .668}
	//create a flow frequency curve based on lp3
	ff := flowfrequencycurves.UnregulatedFlowFrequencyCurve{Distribution: n, Ordinates: 1000}
	//create some basic data that will fill future structs
	inflows := []float64{15000.0, 14000.0, 13000.0, 12000.0, 10000.0, 7000.0, 5000.0, 4000.0, 4000.0, 1500.0, 1000.0, 850.0}
	flows := []float64{14000.0, 13000.0, 12000.0, 11000.0, 7000.0, 5000.0, 4000.0, 3000.0, 1500.0, 1000.0, 900.0, 850.0}
	for i, j := 0, len(flows)-1; i < j; i, j = i+1, j-1 {
		flows[i], flows[j] = flows[j], flows[i]
		inflows[i], inflows[j] = inflows[j], inflows[i]
	}
	flowDists := make([]statistics.ContinuousDistribution, len(flows))
	for i, f := range flows {
		dist, _ := statistics.InitDeterministic(f)
		flowDists[i] = dist
	}
	//create a flow transform
	ftc := flowtransformfunction.FlowTransformFunctionCurve{Curve: paireddata.UncertaintyPairedData{Xvals: inflows, Yvals: flowDists}}

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

	sim := Simulation{FlowFrequency: &ff, RatingCurve: &rc, DamageCurve: &dc, FlowTransform: &ftc}
	ead, err := sim.Compute()
	if err != nil {
		fmt.Println(err)
	}
	expected := 36.30599217405613
	if ead != expected {
		t.Errorf("expected %f, got %f", expected, ead)
	}
}
func Test_EAD_Levee(t *testing.T) {
	n := statistics.LogPearsonIIIDistribution{Mean: 3.368, StandardDeviation: .246, Skew: .668}
	//create a flow frequency curve based on lp3
	ff := flowfrequencycurves.UnregulatedFlowFrequencyCurve{Distribution: n, Ordinates: 1000}
	//create some basic data that will fill future structs
	flows := []float64{14000.0, 13000.0, 12000.0, 11000.0, 7000.0, 5000.0, 4000.0, 3000.0, 1500.0, 1000.0, 900.0, 850.0}
	for i, j := 0, len(flows)-1; i < j; i, j = i+1, j-1 {
		flows[i], flows[j] = flows[j], flows[i]
	}
	stages := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0}
	stageDists := make([]statistics.ContinuousDistribution, len(stages))
	for i, s := range stages {
		dist, _ := statistics.InitDeterministic(s)
		stageDists[i] = dist
	}
	//create a rating curve based on fake data
	rc := ratingcurves.RatingCurve{Curve: paireddata.UncertaintyPairedData{Xvals: flows, Yvals: stageDists}}
	//create a levee
	leveestages := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0}
	leveefailure := []float64{.01, .01, .01, .01, .5, .75}
	failureDists := make([]statistics.ContinuousDistribution, len(leveefailure))
	for i, f := range leveefailure {
		dist, _ := statistics.InitDeterministic(f)
		failureDists[i] = dist
	}
	lc := leveefragilitycurve.LeveeFragilityCurve{Curve: paireddata.UncertaintyPairedData{Xvals: leveestages, Yvals: failureDists}}
	//create some basic data for damages
	damages := []float64{10.0, 20.0, 30.0, 40.0, 50.0, 60.0, 70.0, 80.0, 90.0, 100.0, 110.0, 120.0}
	damageDists := make([]statistics.ContinuousDistribution, len(damages))
	for i, d := range damages {
		dist, _ := statistics.InitDeterministic(d)
		damageDists[i] = dist
	}
	dc := stagedamagecurves.StageDamageCurve{Curve: paireddata.UncertaintyPairedData{Xvals: stages, Yvals: damageDists}}

	sim := Simulation{FlowFrequency: &ff, RatingCurve: &rc, DamageCurve: &dc, LeveeFragilityCurve: &lc}
	ead, err := sim.Compute()
	if err != nil {
		fmt.Println(err)
	}
	expected := 21.930429356104298
	if ead != expected {
		t.Errorf("expected %f, got %f", expected, ead)
	}
}
