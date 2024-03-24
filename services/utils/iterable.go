package utils

import (
	"reflect"
)

func Contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

func MapToInterfaceSlice(m interface{}) []interface{} {
	var s []interface{}
	switch reflect.TypeOf(m).Kind() {
	case reflect.Map:
		v := reflect.ValueOf(m)
		keys := v.MapKeys()
		for _, key := range keys {
			s = append(s, v.MapIndex(key).Interface())
		}
	}
	return s
}
