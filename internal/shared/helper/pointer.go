package helper

import (
	"errors"
	"reflect"
)

func MustBePointer[T any](value any) error {
	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return errors.New("destination must be a non-nil pointer")
	}

	return nil
}
