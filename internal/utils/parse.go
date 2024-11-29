package utils

import (
	"net/http"
	"reflect"
)

func fillFormValueWithSliceStr(
	r *http.Request, field reflect.Value,
	formFieldName string, structField reflect.StructField,
) {
	var	fieldValues	[]string
	var	slice		reflect.Value
	var	i			int

	fieldValues = r.Form[formFieldName]
	if len(fieldValues) > 0 {
		// Create a slice to store the values
		slice = reflect.MakeSlice(
			structField.Type, len(fieldValues), len(fieldValues),
		)
		for i = 0; i < len(fieldValues); i++ {
			slice.Index(i).SetString(fieldValues[i])
		}
		field.Set(slice)
	}
}

func ParseStringForm(r *http.Request, form interface{}) error {
	var	formValue		reflect.Value
	var	formType		reflect.Type
	var	structField		reflect.StructField
	var	formFieldName	string
	var	fieldValue		string
	var	i				int
	var	err				error

	if err = r.ParseForm(); err != nil {
		return err
	}
	formValue = reflect.ValueOf(form).Elem()
	formType = formValue.Type()
	for i = 0; i < formType.NumField(); i++ {
		structField = formType.Field(i)
		formFieldName = structField.Tag.Get("form")
		if formFieldName == "" {
			formFieldName = structField.Name
		}
		// case field value is type []string
		if structField.Type.Kind() == reflect.Slice &&
				structField.Type.Elem().Kind() == reflect.String {
			fillFormValueWithSliceStr(
				r, formValue.Field(i), formFieldName, structField,
			)
		} else {
			fieldValue = r.FormValue(formFieldName)
			formValue.FieldByName(structField.Name).SetString(fieldValue)
		}
	}
	return nil
}
