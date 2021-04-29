package generator

import "strconv"

type intValue struct {
	Prefix string
	Offset int
}

func (v *intValue) init(c Column) {
	v.Prefix = c.Prefix
	v.Offset = c.Offset
}

func (v intValue) getCurrentValue(position int) string {
	return v.Prefix + strconv.Itoa(position+v.Offset)
}
