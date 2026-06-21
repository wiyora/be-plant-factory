package validator

import (
	"reflect"
	"strings"
)

func (v *Validate) Register() {
	v.val.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}

		return name
	})
}
