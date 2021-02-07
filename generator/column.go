package generator

import (
	"errors"
	"math/rand"
	"strings"
)

const (
	intType    = "INT"
	floatType  = "FLOAT"
	dateType   = "DATE"
	stringType = "STRING"
)

type Column struct {
	Type           string   `json:"type"`
	Name           string   `json:"name"`
	Distinct       int      `json:"distinct"`
	Mode           string   `json:"mode,omitempty"` // Mode can be column specific
	Offset         int      `json:"offset,omitempty"`
	Start          string   `json:"start,omitempty"`
	End            string   `json:"end,omitempty"`
	Prefix         string   `json:"prefix,omitempty"`
	ValuesList     string   `json:"values,omitempty"`
	ValuesSlice    []string `json:"-"` // Overload to store valuesList sliced
	valueGenerator value    `json:"-"`
}

func (cp *Column) init(defaultMode string) error {
	// Default prefix for string fields
	if cp.Type == "STRING" && cp.Prefix == "" {
		cp.Prefix = "txt_"
	}
	// Mode default to table mode
	if cp.Mode == "" {
		cp.Mode = defaultMode
	}
	// If no table mode, default to alternate
	if cp.Mode == "" {
		cp.Mode = "ALTERNATE"
	}
	// Split valuesList to slices
	if cp.ValuesList != "" {
		cp.ValuesSlice = strings.Split(cp.ValuesList, ";")
	}

	// Initiate value generator
	valueGenerator, err := createValueGenerator(cp.Type)
	if err != nil {
		return err
	}
	cp.valueGenerator = valueGenerator
	cp.valueGenerator.init(*cp)
	return nil
}

// Check that input type is supported
func ChecksSupportedType(t string) bool {
	_, err := createValueGenerator(t)
	return err == nil
}

func createValueGenerator(t string) (value, error) {
	var v value
	switch t {
	case intType:
		v = &intValue{}
	case floatType:
		v = &floatValue{}
	case dateType:
		v = &dateValue{}
	case stringType:
		v = &stringValue{}
	default:
		return nil, errors.New("Unsupported type " + t)
	}
	return v, nil
}

func (c *Column) getValue(line int, total int) string {
	var v int
	if c.Mode == "BLOCK" {
		v = int(float64(line) / (float64(total) / float64(c.Distinct)))
	} else if c.Mode == "RANDOM" {
		v = rand.Intn(c.Distinct)
	} else {
		v = line % c.Distinct
	}

	// Handle forced value list
	if v < len(c.ValuesSlice) {
		return c.ValuesSlice[v]
	}

	return c.valueGenerator.getCurrentValue(v)
}
