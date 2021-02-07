package generator

import "strconv"

type floatValue struct {
	FloatStart float64
	FloatStep  float64
}

func (v *floatValue) init(c Column) {
	v1, _ := strconv.ParseFloat(c.Start, 32)
	v2, _ := strconv.ParseFloat(c.End, 32)
	if v2 <= v1 {
		v2 = v1 + 1.0
	}
	v.FloatStart = v1
	v.FloatStep = (v2 - v1) / float64(c.Distinct)
}

func (v floatValue) getCurrentValue(position int) string {
	vv := v.FloatStart + v.FloatStep*float64(position)
	return strconv.FormatFloat(vv, 'f', -1, 32)
}
