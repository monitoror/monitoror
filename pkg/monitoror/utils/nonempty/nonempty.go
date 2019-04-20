package nonempty

import (
	"reflect"
	"time"
)

//Byte if value == byte zero value, return def otherwise return value
func Byte(value, def byte) byte {
	if isZero(value) {
		return def
	}
	return value
}

//Float32 if value == float32 zero value, return def otherwise return value
func Float32(value, def float32) float32 {
	if isZero(value) {
		return def
	}
	return value
}

//Float64 if value == float64 zero value, return def otherwise return value
func Float64(value, def float64) float64 {
	if isZero(value) {
		return def
	}
	return value
}

//Int if value == int zero value, return def otherwise return value
func Int(value, def int) int {
	if isZero(value) {
		return def
	}
	return value
}

//Int8 if value == int8 zero value, return def otherwise return value
func Int8(value, def int8) int8 {
	if isZero(value) {
		return def
	}
	return value
}

//Int16 if value == int16 zero value, return def otherwise return value
func Int16(value, def int16) int16 {
	if isZero(value) {
		return def
	}
	return value
}

//Int32 if value == int32 zero value, return def otherwise return value
func Int32(value, def int32) int32 {
	if isZero(value) {
		return def
	}
	return value
}

//Int64 if value == int64 zero value, return def otherwise return value
func Int64(value, def int64) int64 {
	if isZero(value) {
		return def
	}
	return value
}

//Uint if value == uint zero value, return def otherwise return value
func Uint(value, def uint) uint {
	if isZero(value) {
		return def
	}
	return value
}

//Uint8 if value == uint8 zero value, return def otherwise return value
func Uint8(value, def uint8) uint8 {
	if isZero(value) {
		return def
	}
	return value
}

//Uint16 if value == uint16 zero value, return def otherwise return value
func Uint16(value, def uint16) uint16 {
	if isZero(value) {
		return def
	}
	return value
}

//Uint32 if value == uint32 zero value, return def otherwise return value
func Uint32(value, def uint32) uint32 {
	if isZero(value) {
		return def
	}
	return value
}

//Uint64 if value == uint64 zero value, return def otherwise return value
func Uint64(value, def uint64) uint64 {
	if isZero(value) {
		return def
	}
	return value
}

//Uintptr if value == uintptr zero value, return def otherwise return value
func Uintptr(value, def uintptr) uintptr {
	if isZero(value) {
		return def
	}
	return value
}

//Rune if value == rune zero value, return def otherwise return value
func Rune(value, def rune) rune {
	if isZero(value) {
		return def
	}
	return value
}

//String if value == string zero value, return def otherwise return value
func String(value, def string) string {
	if isZero(value) {
		return def
	}
	return value
}

//Time if value == time.Time zero value, return def otherwise return value
func Time(value, def time.Time) time.Time {
	if isZero(value) {
		return def
	}
	return value
}

//Struct if value == struct zero value, return def otherwise return value
func Struct(value, def interface{}) interface{} {
	if isZero(value) {
		return def
	}
	return value
}

func isZero(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}
