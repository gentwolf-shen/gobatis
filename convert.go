package gobatis

import (
	"strconv"
)

func ToInt(str string) int {
	i, _ := strconv.Atoi(str)
	return i
}

func ToInt8(str string) int8 {
	i, _ := strconv.Atoi(str)
	return int8(i)
}

func ToInt16(str string) int16 {
	i, _ := strconv.Atoi(str)
	return int16(i)
}

func ToInt32(str string) int32 {
	i, _ := strconv.Atoi(str)
	return int32(i)
}

func ToInt64(str string) int64 {
	i, _ := strconv.ParseInt(str, 10, 64)
	return i
}

func ToUint(str string) uint {
	return uint(ToUint64(str))
}

func ToUint8(str string) uint8 {
	return uint8(ToUint64(str))
}

func ToUint16(str string) uint16 {
	return uint16(ToUint64(str))
}

func ToUint32(str string) uint32 {
	return uint32(ToUint64(str))
}

func ToUint64(str string) uint64 {
	i, _ := strconv.ParseUint(str, 10, 64)
	return i
}

func ToFloat32(str string) float32 {
	i, _ := strconv.ParseFloat(str, 32)
	return float32(i)
}

func ToFloat64(str string) float64 {
	i, _ := strconv.ParseFloat(str, 64)
	return i
}
