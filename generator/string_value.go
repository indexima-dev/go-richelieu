package generator

import "strconv"

type stringValue struct {
	Prefix string
	Offset int
}

func (v *stringValue) init(c Column) {
	v.Prefix = c.Prefix
	v.Offset = c.Offset
}

func (v stringValue) getCurrentValue(position int) string {
	return v.Prefix + strconv.Itoa(position+v.Offset)
}
