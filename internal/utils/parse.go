package utils

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
)

func isPointerToStruct(i interface{}) bool {
	var	value	reflect.Value

	value = reflect.ValueOf(i)
	if value.Kind() == reflect.Ptr {
		if value.Elem().Kind() == reflect.Struct {
			return true
		}
	}
	return false
}

func fillFieldValueWithSlice(
	r *http.Request, field reflect.Value,
	fieldTag string, fieldType reflect.Type,
) {
	var	fieldValues	[]string
	var	slice		reflect.Value
	var	i			int

	fieldValues = r.Form[fieldTag]
	if len(fieldValues) > 0 {
		// Create a slice to store the values
		slice = reflect.MakeSlice(
			fieldType, len(fieldValues), len(fieldValues),
		)
		for i = 0; i < len(fieldValues); i++ {
			slice.Index(i).SetString(fieldValues[i])
		}
		field.Set(slice)
	}
}

func fillFieldFile(
	r *http.Request, field reflect.Value, fieldTag string,
) {
	var	formFile	FormFile
	var	err			error

	if !field.CanSet() {
		log.Printf("Cannot set field %s\n", fieldTag)
		return
	}
	formFile = FormFile{}
	formFile.File, formFile.FileHeader, err = r.FormFile(fieldTag)
	if err != nil {
		return
	}
	field.Set(reflect.ValueOf(&formFile))
}

func parseField(
	r *http.Request, form reflect.Value, index int,
	structField reflect.StructField,
) {
	var	fieldValue	string
	var	fieldTag	string
	var	fieldType	reflect.Type

	fieldTag = structField.Tag.Get("form") // empty if no tag form
	fieldType = structField.Type
	switch {
	case fieldTag == "":
		return
	case fieldType.Kind() == reflect.Slice &&
	fieldType.Elem().Kind() == reflect.String: // slice of string
		fillFieldValueWithSlice(
			r, form.Field(index), fieldTag, fieldType,
		)
	case fieldType.Kind() == reflect.String: // string
		fieldValue = r.FormValue(fieldTag)
		form.Field(index).SetString(fieldValue)
	case fieldType.Kind() == reflect.Pointer &&
	fieldType.Elem().Kind() == reflect.Struct &&
	fieldType.Elem() == reflect.TypeOf(FormFile{}): // FormFile
		fillFieldFile(r, form.Field(index), fieldTag)
	default:
		log.Printf(
			"Type of %s not supported, skip.\n", fieldType.Name(),
		)
	}
}

func requestParseForm(r *http.Request) error {
	var	contentType	string
	var	err			error

	contentType = r.Header.Get("Content-Type")
	fmt.Println(contentType)
	if contentType == "" {
		log.Println("No content provided with the request.")
		return nil
	}
	if strings.Contains(contentType, "multipart/form-data") {
		if err = r.ParseMultipartForm(20 << 20); err != nil {
			return err
		}
	} else if err = r.ParseForm(); err != nil {
		return err
	}
	return nil
}

// Fill values of struct pointed by `PtToformStruct`
// with the form received on request `r`.
// Only values with the tag `form` and a key that match with
// any value of the form from the request will be written.
func ParseForm(r *http.Request, PtToformStruct interface{}) error {

	var	form			reflect.Value
	var	i				int
	var	structField		reflect.StructField
	var	err				error

	if err = requestParseForm(r); err != nil {
		log.Printf("Error at request form parsing: %v\n", err)
		return err
	}
	if (!isPointerToStruct(PtToformStruct)) {
		return fmt.Errorf("not a pointer to struct")
	}
	form = reflect.ValueOf(PtToformStruct).Elem()
	for i = 0; i < form.NumField(); i++ {
		structField = form.Type().Field(i) // metadata of the field
		parseField(r, form, i, structField)
	}
	return nil
}
