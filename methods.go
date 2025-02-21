package ThunderORM

import (
	"fmt"
	"reflect"

	"github.com/fatih/structs"
)

// IsStruct returns true if the given interface is a struct.
func IsStruct(i interface{}) bool {
	return structs.IsStruct(i)
}

// Name returns the name of the struct.
func Name(i interface{}) string {
	return structs.Name(i)
}

// Fields returns the names of the fields in the struct.
func Fields(i interface{}) []string {
	return structs.Names(i)
}

// Values returns the values of the fields in the struct.
func Values(i interface{}) []interface{} {
	return structs.Values(i)
}

// Get returns the value of the field in a struct.
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
