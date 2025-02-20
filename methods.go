package ThunderORM

import (
	"reflect"

	"github.com/fatih/structs"
)

func IsStruct(i interface{}) bool {
	return structs.IsStruct(i)
}

func Name(i interface{}) string {
	return structs.Name(i)
}

func Fields(i interface{}) []string {
	return structs.Names(i)
}

func Values(i interface{}) []interface{} {
	return structs.Values(i)
}

func Get(i interface{}, key string) interface{} {
	if !IsStruct(i) {
		return nil
	}
	m := structs.Map(i)
	return m[key]
}

func Set(i interface{}, field string, value interface{}) interface{} {
	reflect.ValueOf(i).Elem().FieldByName(field).Set(reflect.ValueOf(value))
	return i
}
