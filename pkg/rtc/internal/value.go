package internal

import (
	"strconv"

	"rtc/pkg/rtc"
)

// Value implementation rtc.Value from bytes
type Value struct {
	src []byte
}

// NewValue ...
func NewValue(src []byte) *Value {
	return &Value{src: src}
}

// String ...
func (v *Value) String() string {
	return string(v.src)
}

// MaybeString ...
func (v *Value) MaybeString() (string, error) {
	val := string(v.src)
	if val == "" {
		return "", rtc.ErrNotPresent
	}

	return val, nil
}

// Float64 ...
func (v *Value) Float64() float64 {
	val, _ := v.MaybeFloat64()
	return val
}

// MaybeFloat64 ...
func (v *Value) MaybeFloat64() (float64, error) {
	return strconv.ParseFloat(v.String(), 64)
}

// Bool ...
func (v *Value) Bool() bool {
	val, _ := v.MaybeBool()
	return val
}

// MaybeBool ...
func (v *Value) MaybeBool() (bool, error) {
	return strconv.ParseBool(v.String())
}

// Int ...
func (v *Value) Int() int {
	val, _ := v.MaybeInt()
	return val
}

// MaybeInt ...
func (v *Value) MaybeInt() (int, error) {
	return strconv.Atoi(v.String())
}
