package utils

import "reflect"

func StructToMap(item interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	val := reflect.ValueOf(item)
	typ := reflect.TypeOf(item)

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Skip zero value fields; i.e. fields that are not set
		if field.IsZero() {
			continue
		}

		tag := fieldType.Tag.Get("json")
		if tag == "" {
			tag = fieldType.Name
		} else {
			// Remove ",omitempty" if present
			tag = tag[:len(tag)-len(",omitempty")]
		}

		result[tag] = field.Interface()
	}

	return result
}
