package generator

import "time"

type dateValue struct {
	DateStart int64
	DateStep  int64
}

func (v *dateValue) init(c Column) {
	if c.Start == "" {
		c.Start = "2020-01-01 00:00:00"
	}
	if c.End == "" {
		c.End = "2020-12-31 00:00:00"
	}
	v1, _ := time.Parse("2006-01-02 15:04:05", c.Start)
	v2, _ := time.Parse("2006-01-02 15:04:05", c.End)
	v.DateStart = v1.Unix()
	v.DateStep = (v2.Unix() - v1.Unix()) / int64(c.Distinct)
}

func (v dateValue) getCurrentValue(position int) string {
	vv := v.DateStart + v.DateStep*int64(position)
	tm := time.Unix(vv, 0)
	return tm.Local().String()
}
