package structmapping

import (
	"reflect"
	"fmt"
)

func New() *StructMapper {
	return &StructMapper{customWrappers:make(map[string]interface{})}
}

type StructMapper struct {
	customWrappers map[string]interface{}
}

func (s *StructMapper) Add(w interface{}) {
	funcType := reflect.ValueOf(w).Type()
	s.customWrappers[fmt.Sprintf("%s-%s", funcType.In(0), funcType.In(1))] = w
}

func (m *StructMapper) Map(src, dst interface{}) {

	typeConversion := fmt.Sprintf("%s-%s", reflect.ValueOf(src).Type(), reflect.ValueOf(dst).Type())
	fmt.Println("Enter Map() ", typeConversion)
	defer fmt.Println("Leaving Map() ", typeConversion)

	// delegate to custom mapper if needed
	if cMapper, ok := m.customWrappers[typeConversion]; ok {
		fnValue := reflect.ValueOf(cMapper)
		fnValue.Call([]reflect.Value{reflect.ValueOf(src), reflect.ValueOf(dst)})
		return
	}

	// else, continue mapping of each field here

	// unwrap pointer
	valueSrc := reflect.ValueOf(src).Elem()
	valueDst := reflect.ValueOf(dst).Elem()

	m.mapValue(valueSrc, valueDst)
}

func (m *StructMapper) mapValue(src, dst reflect.Value) {
	fmt.Println("Enter mapValue for ", src, src.Type(), src.Kind(), dst, dst.Type(), dst.Kind())
	defer fmt.Println("Leave mapValue for ", src, src.Type(), src.Kind(), dst, dst.Type(), dst.Kind())
	if src.Kind() == reflect.Map && dst.Kind() == reflect.Map {
		dst.Set(reflect.MakeMap(dst.Type()))
		for _, k := range src.MapKeys() {
			dstValue := reflect.New(dst.Type().Elem())
			m.mapValue(src.MapIndex(k), dstValue)
			dst.SetMapIndex(k, dstValue)
		}
	} else if src.Kind() == reflect.Ptr && dst.Kind() == reflect.Ptr {
		m.mapValue(src.Elem(), dst.Elem())
	} else if src.Kind() == reflect.Struct && dst.Kind() == reflect.Struct {
		for i := 0; i < src.NumField(); i++ {
			srcFieldName := src.Type().Field(i).Name
			foundField := false

			// find the field with the same name in dst
			for j := 0; j < dst.NumField(); j++ {
				dstFieldName := dst.Type().Field(j).Name

				if dstFieldName == srcFieldName {
					foundField = true
					m.mapValue(src.Field(i), dst.Field(j))
				}
			}

			if !foundField {
				panic("Could not find field " + srcFieldName + " in dst")
			}
		}
	} else if isPrimitive(src.Kind()) && isPrimitive(dst.Kind()) {
		// all primitive types remain, simply copy the value over
		reflect.New(dst.Type()).Elem().Set(src)
	} else {
		panic("Do not know how to convert " + src.Kind().String() + "->" + dst.Kind().String())
	}
}

func isPrimitive(v reflect.Kind) bool {
	return v == reflect.Bool ||
	v == reflect.Int ||
	v == reflect.Int8 ||
	v == reflect.Int16 ||
	v == reflect.Int32 ||
	v == reflect.Int64 ||
	v == reflect.Uint ||
	v == reflect.Uint8 ||
	v == reflect.Uint16 ||
	v == reflect.Uint32 ||
	v == reflect.Uint64 ||
	v == reflect.Float32 ||
	v == reflect.Float64 ||
	v == reflect.Complex64 ||
	v == reflect.Complex128 ||
	v == reflect.String
}