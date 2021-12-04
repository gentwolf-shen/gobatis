package gobatis

import (
	"reflect"
	"strings"
)

func compare(value interface{}, op, testValue string) bool {
	bl := false

	switch reflect.TypeOf(value).Name() {
	case "int":
		bl = compareInt(value.(int), ToInt(testValue), op)
	case "int8":
		bl = compareInt8(value.(int8), ToInt8(testValue), op)
	case "int16":
		bl = compareInt16(value.(int16), ToInt16(testValue), op)
	case "int32":
		bl = compareInt32(value.(int32), ToInt32(testValue), op)
	case "int64":
		bl = compareInt64(value.(int64), ToInt64(testValue), op)
	case "uint":
		bl = compareUint(value.(uint), ToUint(testValue), op)
	case "uint8":
		bl = compareUint8(value.(uint8), ToUint8(testValue), op)
	case "uint16":
		bl = compareUint16(value.(uint16), ToUint16(testValue), op)
	case "uint32":
		bl = compareUint32(value.(uint32), ToUint32(testValue), op)
	case "uint64":
		bl = compareUint64(value.(uint64), ToUint64(testValue), op)
	case "float32":
		bl = compareFloat32(value.(float32), ToFloat32(testValue), op)
	case "float64":
		bl = compareFloat64(value.(float64), ToFloat64(testValue), op)
	case "string":
		bl = compareString(value.(string), strings.Trim(strings.Trim(testValue, "'"), "\""), op)
	}

	return bl
}

func compareInt(value, target int, op string) bool {
	bl := false
	switch op {
	case "!=":
		bl = value != target

	case "==":
		fallthrough
	case "=":
		bl = value == target

	case ">":
		bl = value > target

	case ">=":
		bl = value >= target

	case "<":
		bl = value < target

	case "<=":
		bl = value <= target

	}

	return bl
}

func compareInt8(value, target int8, op string) bool {
	bl := false
	switch op {
	case "!=":
		bl = value != target

	case "==":
		fallthrough
	case "=":
		bl = value == target

	case ">":
		bl = value > target

	case ">=":
		bl = value >= target

	case "<":
		bl = value < target

	case "<=":
		bl = value <= target

	}

	return bl
}

func compareInt16(value, target int16, op string) bool {
	bl := false
	switch op {
	case "!=":
		bl = value != target

	case "==":
		fallthrough
	case "=":
		bl = value == target

	case ">":
		bl = value > target

	case ">=":
		bl = value >= target

	case "<":
		bl = value < target

	case "<=":
		bl = value <= target

	}

	return bl
}

func compareInt32(value, target int32, op string) bool {
	bl := false
	switch op {
	case "!=":
		bl = value != target

	case "==":
		fallthrough
	case "=":
		bl = value == target

	case ">":
		bl = value > target

	case ">=":
		bl = value >= target

	case "<":
		bl = value < target

	case "<=":
		bl = value <= target

	}

	return bl
}

func compareInt64(value, target int64, op string) bool {
	bl := false
	switch op {
	case "!=":
		bl = value != target

	case "==":
		fallthrough
	case "=":
		bl = value == target

	case ">":
		bl = value > target

	case ">=":
		bl = value >= target

	case "<":
		bl = value < target

	case "<=":
		bl = value <= target

	}

	return bl
}

func compareUint(value, target uint, op string) bool {
	bl := false
	switch op {
	case "!=":
		bl = value != target

	case "==":
		fallthrough
	case "=":
		bl = value == target

	case ">":
		bl = value > target

	case ">=":
		bl = value >= target

	case "<":
		bl = value < target

	case "<=":
		bl = value <= target

	}

	return bl
}

func compareUint8(value, target uint8, op string) bool {
	bl := false
	switch op {
	case "!=":
		bl = value != target

	case "==":
		fallthrough
	case "=":
		bl = value == target

	case ">":
		bl = value > target

	case ">=":
		bl = value >= target

	case "<":
		bl = value < target

	case "<=":
		bl = value <= target

	}

	return bl
}

func compareUint16(value, target uint16, op string) bool {
	bl := false
	switch op {
	case "!=":
		bl = value != target

	case "==":
		fallthrough
	case "=":
		bl = value == target

	case ">":
		bl = value > target

	case ">=":
		bl = value >= target

	case "<":
		bl = value < target

	case "<=":
		bl = value <= target

	}

	return bl
}

func compareUint32(value, target uint32, op string) bool {
	bl := false
	switch op {
	case "!=":
		bl = value != target

	case "==":
		fallthrough
	case "=":
		bl = value == target

	case ">":
		bl = value > target

	case ">=":
		bl = value >= target

	case "<":
		bl = value < target

	case "<=":
		bl = value <= target

	}

	return bl
}

func compareUint64(value, target uint64, op string) bool {
	bl := false
	switch op {
	case "!=":
		bl = value != target

	case "==":
		fallthrough
	case "=":
		bl = value == target

	case ">":
		bl = value > target

	case ">=":
		bl = value >= target

	case "<":
		bl = value < target

	case "<=":
		bl = value <= target

	}

	return bl
}

func compareFloat32(value, target float32, op string) bool {
	bl := false
	switch op {
	case "!=":
		bl = value != target

	case "==":
		fallthrough
	case "=":
		bl = value == target

	case ">":
		bl = value > target

	case ">=":
		bl = value >= target

	case "<":
		bl = value < target

	case "<=":
		bl = value <= target

	}

	return bl
}

func compareFloat64(value, target float64, op string) bool {
	bl := false
	switch op {
	case "!=":
		bl = value != target

	case "==":
		fallthrough
	case "=":
		bl = value == target

	case ">":
		bl = value > target

	case ">=":
		bl = value >= target

	case "<":
		bl = value < target

	case "<=":
		bl = value <= target

	}

	return bl
}

func compareString(value, target string, op string) bool {
	bl := false
	switch op {
	case "!=":
		bl = value != target
	case "==":
		fallthrough
	case "=":
		bl = value == target

	case ">":
		bl = value > target

	case ">=":
		bl = value >= target

	case "<":
		bl = value < target

	case "<=":
		bl = value <= target

	}

	return bl
}
