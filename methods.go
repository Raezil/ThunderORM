package ThunderORM

import (
	"fmt"
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
	return structs.Map(i)[key]
}

// Set updates the field in a struct pointer with the given value.
// It returns an error if the field is missing or not settable.
func Set(i interface{}, field string, value interface{}) error {
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("Set: input must be a non-nil pointer")
	}
	elem := v.Elem()
	f := elem.FieldByName(field)
	if !f.IsValid() {
		return fmt.Errorf("Set: no such field %s", field)
	}
	if !f.CanSet() {
		return fmt.Errorf("Set: cannot set field %s", field)
	}
	f.Set(reflect.ValueOf(value))
	return nil
}
