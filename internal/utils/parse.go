package utils

import (
	"net/http"
	"reflect"
)

func ParseForm(r *http.Request, form interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	formValue := reflect.ValueOf(form).Elem()
	formType := formValue.Type()
	for i := 0; i < formType.NumField(); i++ {
		field := formType.Field(i)
		formFieldName := field.Tag.Get("form")
		if formFieldName == "" {
			formFieldName = field.Name
		}
		fieldValue := r.FormValue(formFieldName)
		formValue.FieldByName(field.Name).SetString(fieldValue)
	}
	return nil
}
