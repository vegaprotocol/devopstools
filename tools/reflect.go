package tools

import "reflect"

func StructSize[T any](variable T) int {
	return int(reflect.TypeOf(variable).Size())
}
