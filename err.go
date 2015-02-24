package mirror

import (
	"fmt"
	"reflect"
)

func mustSameType(x, y reflect.Value) {
	if x.Type() != y.Type() {
		panic(fmt.Sprintf("%v and %v have different types: %v != %v", x, y, x.Type(), y.Type()))
	}
}

func mustUnsigned(x reflect.Value) {
	switch x.Type().Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return
	}

	panic(fmt.Sprintf("%x is not unsigned", x.Type()))
}
