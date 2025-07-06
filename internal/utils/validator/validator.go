package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func New() *validator.Validate {
	validate := validator.New()

	// Using the names which have been specified for JSON representations of structs, rather than normal Go field names.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return validate
}
