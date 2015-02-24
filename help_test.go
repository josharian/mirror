// This file contains testing helpers.

package mirror

import (
	"math"
	"math/cmplx"
	"reflect"
)

func FlexEqual(x, y interface{}) bool {
	var isfloat, iscomplex bool
	var xf, yf float64
	var xc, yc complex128

	switch x := x.(type) {
	case float32:
		if _, ok := y.(float32); !ok {
			return false
		}
		xf, yf = float64(x), float64(y.(float32))
		isfloat = true
	case float64:
		if _, ok := y.(float64); !ok {
			return false
		}
		xf, yf = float64(x), y.(float64)
		isfloat = true
	case complex64:
		if _, ok := y.(complex64); !ok {
			return false
		}
		xc, yc = complex128(x), complex128(y.(complex64))
		iscomplex = true
	case complex128:
		if _, ok := y.(complex128); !ok {
			return false
		}
		xc, yc = complex128(x), y.(complex128)
		iscomplex = true
	}

	if isfloat {
		switch {
		case math.IsNaN(xf) && math.IsNaN(yf):
			return true
		case math.IsInf(xf, 1) && math.IsInf(yf, 1):
			return true
		case math.IsInf(xf, -1) && math.IsInf(yf, -1):
			return true
		}
	}

	if iscomplex {
		switch {
		case cmplx.IsNaN(xc) && cmplx.IsNaN(yc):
			return true
		case cmplx.IsInf(xc) && cmplx.IsInf(yc):
			return true
		}
	}

	return reflect.DeepEqual(x, y)
}
