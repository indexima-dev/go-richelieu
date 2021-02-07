package generator

type value interface {
	getCurrentValue(position int) string
	init(c Column)
}
