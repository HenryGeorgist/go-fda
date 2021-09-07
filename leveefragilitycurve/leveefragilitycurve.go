package leveefragilitycurve

import (
	"github.com/USACE/go-consequences/paireddata"
)

type LeveeFragilityCurve struct {
	Curve paireddata.UncertaintyPairedData
}

func (s LeveeFragilityCurve) Sample(randomValue float64) paireddata.ValueSampler {
	return s.Curve.SampleValueSampler(randomValue)
}
func Multiply(damagecurve paireddata.PairedData, leveecurve paireddata.PairedData) paireddata.PairedData {
	belowFragilityCurveValue := 0.0
	aboveFragilityCurveValue := 1.0
	newXvals := make([]float64, 0)
	newYvals := make([]float64, 0)
	if damagecurve.Xvals[0] < leveecurve.Xvals[0] {
		//cacluate no damage until the bottom of the fragility curve
		bottom := leveecurve.Xvals[0]
		for _, dcx := range damagecurve.Xvals {
			if dcx < bottom {
				//set to zero
				newXvals = append(newXvals, dcx)
				newYvals = append(newYvals, belowFragilityCurveValue)
			} else {
				//create a point on the curve just below the bottom of the levee at damage zero.
				newXvals = append(newXvals, bottom-.000000000001)
				newYvals = append(newYvals, belowFragilityCurveValue)
				//create a point at the bottom of the fragility curve
				newXvals = append(newXvals, bottom)
				damage := damagecurve.SampleValue(bottom) * leveecurve.Yvals[0]
				newYvals = append(newYvals, damage)
				break
			}
		}
	}
	for idx, lcx := range leveecurve.Xvals {
		//modify
		damage := damagecurve.SampleValue(lcx) * leveecurve.Yvals[idx]
		newXvals = append(newXvals, lcx)
		newYvals = append(newYvals, damage)
	}
	if leveecurve.Xvals[len(leveecurve.Xvals)-1] < damagecurve.Xvals[len(damagecurve.Xvals)-1] {
		//add in the damage curve ordinates without modification.
		top := leveecurve.Xvals[len(leveecurve.Xvals)-1]
		newXvals = append(newXvals, top)
		damage := damagecurve.SampleValue(top) * leveecurve.Yvals[len(leveecurve.Yvals)-1]
		newYvals = append(newYvals, damage)
		//create a point at the bottom of the fragility curve
		newXvals = append(newXvals, top+.00000001)
		damageabove := damagecurve.SampleValue(top+.00000001) * aboveFragilityCurveValue
		newYvals = append(newYvals, damageabove)
		for idx, dcx := range damagecurve.Xvals {
			if dcx > top {
				//set to max val
				newXvals = append(newXvals, dcx)
				damage := damagecurve.Yvals[idx] * aboveFragilityCurveValue
				newYvals = append(newYvals, damage)

			}
		}
	}
	return paireddata.PairedData{Xvals: newXvals, Yvals: newYvals}
}
