package zero

import "reflect"

func IsStructZero(s any) bool {
	if reflect.ValueOf(s).IsZero() {
		return true
	}

	return false
}
