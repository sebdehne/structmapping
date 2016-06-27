package structmapping

import (
	"reflect"
	"fmt"
)

// Takes a pointer to a struct and searches the fields for nil-maps and slices and initializes them
func NewWithMaps(ptr interface{}) interface{} {
	if reflect.ValueOf(ptr).Kind() != reflect.Ptr {
		panic("Argument must be a pointer")
	}
	makeMapsInternal(reflect.ValueOf(ptr).Elem())

	return ptr
}

func makeMapsInternal(v reflect.Value) {

	if debug {
		fmt.Println("Enter MakeMaps for ", v, v.Type(), v.Kind())
		defer fmt.Println("Leave mapValue for ", v, v.Type(), v.Kind())
	}

	if v.Kind() == reflect.Map  && v.IsNil() {
		v.Set(reflect.MakeMap(v.Type()))
	} else if v.Kind() == reflect.Slice && v.IsNil() {
		v.Set(reflect.MakeSlice(v.Type(), v.Len(), v.Cap()))
	} else if v.Kind() == reflect.Ptr && !v.IsNil() {
		makeMapsInternal(v.Elem())
	} else if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			makeMapsInternal(v.Field(i))
		}
	}

}
