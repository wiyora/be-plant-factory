package helper

import "reflect"

func IsEmptyStruct(s any) bool {
	return reflect.ValueOf(s).IsZero()
}
