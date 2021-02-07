package generator

import (
	"testing"
)

func TestFloatValueInit(t *testing.T) {
	var r floatValue
	initValue := Column{Type: "FLOAT", Distinct: 5, Start: "3.00", End: "9.00"}
	r.init(initValue)

	if int64(r.FloatStart) != 3.0 {
		t.Errorf("Incorrect FloatStart %v", r.FloatStart)
	}
	if r.FloatStep != 1.2 {
		t.Errorf("Incorrect FloatStep %v", r.FloatStep)
	}
}

func TestFloatValueGetCurrentValue(t *testing.T) {
	var r floatValue
	initValue := Column{Type: "FLOAT", Distinct: 5, Start: "3.00", End: "9.00"}
	r.init(initValue)
	var e string = "19.8"
	ret := r.getCurrentValue(14)
	if ret != e {
		t.Errorf("Incorrect value %v. Expected %v", ret, e)
	}
}
