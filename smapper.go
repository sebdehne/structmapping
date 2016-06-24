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

	valueSrc := reflect.ValueOf(src).Elem()
	typeSrc := valueSrc.Type()
	valueDst := reflect.ValueOf(dst).Elem()
	typeDst := valueDst.Type()

	typeConversion := fmt.Sprintf("%s-%s", typeSrc, typeDst)
	fmt.Println("Enter Map() ", typeConversion)
	defer fmt.Println("Leaving Map() ", typeConversion)

	// delegate to custom mapper if needed
	if cMapper, ok := m.customWrappers[typeConversion]; ok {
		fn := cMapper.(func(interface{}, interface{}))
		fn(src, dst)
		return
	}
	// else, continue mapping here

	for i := 0; i < valueSrc.NumField(); i++ {
		srcFieldName := typeSrc.Field(i).Name
		srcFieldType := typeSrc.Field(i).Type
		foundField := false

		// find the field with the same name in dst
		for j := 0; j < valueDst.NumField(); j++ {
			dstFieldName := typeDst.Field(j).Name
			dstFieldType := typeDst.Field(j).Type

			if dstFieldName == srcFieldName {
				foundField = true

				if srcFieldType.Kind() == reflect.Map && dstFieldType.Kind() == reflect.Map {
					// TODO
				} else {
					panic("Do not know how to convert " + srcFieldType.Kind().String() + "->" + dstFieldType.Kind().String())
				}
			}
		}

		if !foundField {
			panic("Could not find field " + srcFieldName + " in dst")
		}
	}
}
