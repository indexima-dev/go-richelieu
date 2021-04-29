package generator

import (
	"testing"
)

func TestIntValueInit(t *testing.T) {
	var r intValue
	initValue := Column{Type: "INT", Distinct: 33, Offset: 1000}
	r.init(initValue)
	if r.Prefix != "" {
		t.Errorf("Incorrect prefix %s.", r.Prefix)
	}
	if r.Offset != 1000 {
		t.Errorf("Incorrect offset %d.", r.Offset)
	}
}

func TestIntValueGetCurrentValue(t *testing.T) {
	var r intValue
	initValue := Column{Type: "INT", Distinct: 25, Offset: 1000}
	r.init(initValue)
	ret := r.getCurrentValue(24)
	if ret != "1024" {
		t.Errorf("Incorrect value %s.", ret)
	}
}
