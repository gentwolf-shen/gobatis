package gobatis

import (
	"fmt"
	"reflect"
)

func toStr(n interface{}) string {
	if reflect.ValueOf(n).Kind() == reflect.String {
		return n.(string)
	}
	return fmt.Sprintf("%v", n)
}
