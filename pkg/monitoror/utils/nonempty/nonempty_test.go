package nonempty

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestByte(t *testing.T) {
	var a, b, c, d byte
	a, b, d = 10, 20, 30
	assert.Equal(t, a, Byte(a, b))
	assert.Equal(t, d, Byte(c, d))
}

func TestFloat32(t *testing.T) {
	var a, b, c, d float32
	a, b, d = 10.0, 20.0, 30.0
	assert.Equal(t, a, Float32(a, b))
	assert.Equal(t, d, Float32(c, d))
}

func TestFloat64(t *testing.T) {
	var a, b, c, d float64
	a, b, d = 10.0, 20.0, 30.0
	assert.Equal(t, a, Float64(a, b))
	assert.Equal(t, d, Float64(c, d))
}

func TestInt(t *testing.T) {
	var a, b, c, d int
	a, b, d = 10, 20, 30
	assert.Equal(t, a, Int(a, b))
	assert.Equal(t, d, Int(c, d))
}

func TestInt8(t *testing.T) {
	var a, b, c, d int8
	a, b, d = 10, 20, 30
	assert.Equal(t, a, Int8(a, b))
	assert.Equal(t, d, Int8(c, d))
}

func TestInt16(t *testing.T) {
	var a, b, c, d int16
	a, b, d = 10, 20, 30
	assert.Equal(t, a, Int16(a, b))
	assert.Equal(t, d, Int16(c, d))
}

func TestInt32(t *testing.T) {
	var a, b, c, d int32
	a, b, d = 10, 20, 30
	assert.Equal(t, a, Int32(a, b))
	assert.Equal(t, d, Int32(c, d))
}

func TestInt64(t *testing.T) {
	var a, b, c, d int64
	a, b, d = 10, 20, 30
	assert.Equal(t, a, Int64(a, b))
	assert.Equal(t, d, Int64(c, d))
}

func TestUint(t *testing.T) {
	var a, b, c, d uint
	a, b, d = 10, 20, 30
	assert.Equal(t, a, Uint(a, b))
	assert.Equal(t, d, Uint(c, d))
}

func TestUint8(t *testing.T) {
	var a, b, c, d uint8
	a, b, d = 10, 20, 30
	assert.Equal(t, a, Uint8(a, b))
	assert.Equal(t, d, Uint8(c, d))
}

func TestUint16(t *testing.T) {
	var a, b, c, d uint16
	a, b, d = 10, 20, 30
	assert.Equal(t, a, Uint16(a, b))
	assert.Equal(t, d, Uint16(c, d))
}

func TestUint32(t *testing.T) {
	var a, b, c, d uint32
	a, b, d = 10, 20, 30
	assert.Equal(t, a, Uint32(a, b))
	assert.Equal(t, d, Uint32(c, d))
}

func TestUint64(t *testing.T) {
	var a, b, c, d uint64
	a, b, d = 10, 20, 30
	assert.Equal(t, a, Uint64(a, b))
	assert.Equal(t, d, Uint64(c, d))
}

func TestUintptr(t *testing.T) {
	var a, b, c, d uintptr
	a, b, d = 10, 20, 30
	assert.Equal(t, a, Uintptr(a, b))
	assert.Equal(t, d, Uintptr(c, d))
}

func TestRune(t *testing.T) {
	var a, b, c, d rune
	a, b, d = 10.0, 20.0, 30.0
	assert.Equal(t, a, Rune(a, b))
	assert.Equal(t, d, Rune(c, d))
}

func TestString(t *testing.T) {
	var a, b, c, d string
	a, b, d = "a", "b", "d"
	assert.Equal(t, a, String(a, b))
	assert.Equal(t, d, String(c, d))
}

func TestTime(t *testing.T) {
	var a, b, c, d time.Time
	a, b, d = time.Now(), time.Now().Add(time.Hour), time.Now().Add(time.Hour*24)
	assert.Equal(t, a, Time(a, b))
	assert.Equal(t, d, Time(c, d))
}

type Data struct {
	a string
	b int
}

func TestStruct(t *testing.T) {
	var a, b, c, d Data
	a, b, d = Data{"a", 10}, Data{"b", 10}, Data{"d", 30}
	assert.Equal(t, a, Struct(a, b))
	assert.Equal(t, d, Struct(c, d))
}
