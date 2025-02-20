package ThunderORM

import "reflect"

// Scanning returns a slice of pointers to the fields of the struct pointed to by u.
// It is used to build the targets for row scanning.
func Scanning(u interface{}) []interface{} {
	val := reflect.ValueOf(u).Elem()
	targets := make([]interface{}, val.NumField())
	for i := 0; i < val.NumField(); i++ {
		targets[i] = val.Field(i).Addr().Interface()
	}
	return targets
}
