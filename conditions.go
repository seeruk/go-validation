package validation

import "reflect"

// IsEmpty checks if a value is "empty". Basically, if it's a value's zero value, nil, or has no
// length it's considered empty.
func IsEmpty(value reflect.Value) bool {
	if !value.IsValid() {
		return true
	}

	if value.IsZero() {
		return true
	}

	switch value.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		if value.Len() == 0 {
			return true
		}
	}

	return false
}

// IsNillable returns true if the given value is of a nillable type.
func IsNillable(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.Array, reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		return true
	}

	return false
}
