package flowtransformfunction

import (
	"github.com/USACE/go-consequences/paireddata"
)

type FlowTransformFunctionCurve struct {
	Curve paireddata.UncertaintyPairedData
}

func (s FlowTransformFunctionCurve) Sample(randomValue float64) paireddata.ValueSampler {
	return s.Curve.SampleValueSampler(randomValue)
}
