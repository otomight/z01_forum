package templates

import (
	"fmt"
	"forum/internal/utils"
)

func addToStruct(structData interface{},
					values ...interface{}) map[string]interface{} {
	var newData	map[string]interface{}
	var dataMap	map[string]interface{}
	var i		int
	var key		string
	var ok		bool

	dataMap, err := utils.StructToMap(structData)
	if err != nil {
		fmt.Println("Error converting struct to map:", err)
	}
	newData = make(map[string]interface{})
	// add base data
	for key, value := range dataMap {
		newData[key] = value
	}
	// add values to new data
	for i = 0; i < len(values); i += 2 {
		key, ok = values[i].(string)
		if ok {
			newData[key] = values[i + 1]
		}
	}
	return newData
}

// func dict(values ...interface{}) (map[string]interface{}, error) {
// 	var dict	map[string]interface{}
// 	var key		string
// 	var ok		bool
// 	var i		int

// 	if len(values) % 2 != 0 {
// 		return nil, fmt.Errorf("invalid dict call: odd number of arguments")
// 	}
// 	dict = make(map[string]interface{}, len(values) / 2)
// 	for i = 0; i < len(values); i += 2 {
// 		key, ok = values[i].(string)
// 		if !ok {
// 			return nil, fmt.Errorf("dict keys must be strings")
// 		}
// 		dict[key] = values[i + 1]
// 	}
// 	return dict, nil
// }
