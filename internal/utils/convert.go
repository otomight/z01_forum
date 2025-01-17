package utils

import (
	"fmt"
	"reflect"
	"strconv"
)

func StructToMap(data interface{}) (map[string]interface{}, error) {
	var	value		reflect.Value
	var	field		reflect.Value
	var	fieldName	string
	var	i			int
	var	ok			bool
	var	result		map[string]interface{}

	if result, ok = data.(map[string]interface{}); ok {
		return result, nil
	}
	value = reflect.ValueOf(data)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return nil, fmt.Errorf("data is not a struct")
	}
	result = make(map[string]interface{})
	for i = 0; i < value.NumField(); i++ {
		field = value.Field(i)
		fieldName = value.Type().Field(i).Name
		result[fieldName] = field.Interface()
	}
	return result, nil
}

func StrSliceToIntSlice(slice []string) ([]int, error) {
	var	result	[]int
	var	i		int
	var	id		int
	var	err		error

	for i = 0; i < len(slice); i++ {
		id, err = strconv.Atoi(slice[i])
		if err != nil {
			return result, err
		}
		result = append(result, id)
	}
	return result, nil
}
