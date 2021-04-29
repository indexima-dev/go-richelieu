package generator

import (
	"testing"
)

func TestStringValueInit(t *testing.T) {
	var r intValue
	initValue := Column{Type: "STRING", Distinct: 33, Offset: 1000, Prefix: "prefix_"}
	r.init(initValue)

	if r.Prefix != "prefix_" {
		t.Errorf("Incorrect prefix %s.", r.Prefix)
	}
	if r.Offset != 1000 {
		t.Errorf("Incorrect offset %d.", r.Offset)
	}
}

func TestStringValueGetCurrentValue(t *testing.T) {
	var r intValue
	initValue := Column{Type: "STRING", Distinct: 25, Offset: 1000, Prefix: "prefix_"}
	r.init(initValue)
	ret := r.getCurrentValue(24)
	if ret != "prefix_1024" {
		t.Errorf("Incorrect value %s.", ret)
	}
}
