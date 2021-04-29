package generator

import (
	"testing"
	"time"
)

func TestDatetimeValueInit(t *testing.T) {
	var r datetimeValue
	initValue := Column{Type: "DATETIME", Distinct: 5, Start: "2020-01-01 15:00:00", End: "2020-01-02 15:05:00"}
	r.init(initValue)

	if r.DateStart != time.Date(2020, 01, 01, 15, 00, 00, 000000000, time.UTC).Unix() {
		t.Errorf("Incorrect DateStart %v", r.DateStart)
	}
	if r.DateStep != 17340 {
		t.Errorf("Incorrect DateStep %v", r.DateStep)
	}
}

func TestDatetimeValueGetCurrentValue(t *testing.T) {
	var r datetimeValue
	initValue := Column{Type: "DATETIME", Distinct: 5, Start: "2020-01-01 15:00:00", End: "2020-01-02 15:05:00"}
	r.init(initValue)
	var e string = "2020-01-04 11:26:00 +0100 CET"
	ret := r.getCurrentValue(14)
	if ret != e {
		t.Errorf("Incorrect value %v. Expected %v", ret, e)
	}
}
