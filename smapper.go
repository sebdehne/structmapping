package structmapping

import (
	"reflect"
	"fmt"
)

type MappingMode int

//noinspection GoUnusedConst
const (
	debug = false
	SrcFieldBased MappingMode = iota
	DstFieldBased
)

func New(mappingMode MappingMode, ignoreUnknownField bool) *StructMapper {
	return &StructMapper{
		customWrappers:make(map[string]interface{}),
		mappingMode:mappingMode,
		ignoreUnknownField:ignoreUnknownField}
}

type StructMapper struct {
	customWrappers map[string]interface{}
	mappingMode MappingMode
	ignoreUnknownField bool
}

func (s *StructMapper) Add(w interface{}) {
	funcType := reflect.ValueOf(w).Type()

	if funcType.Kind() != reflect.Func || funcType.NumIn() != 2 {
		panic("Only functions with two arguments can be used")
	} else if funcType.In(0).Kind() != reflect.Struct {
		panic("First argument must be a struct")
	} else if funcType.In(1).Kind() != reflect.Ptr || funcType.In(1).Elem().Kind() != reflect.Struct {
		panic("Second argument must be pointer to a struct")
	}

	s.customWrappers[fmt.Sprintf("%s-%s", funcType.In(0), funcType.In(1).Elem())] = w
}

func (m *StructMapper) Map(src, dst interface{}) {

	if debug {
		typeConversion := fmt.Sprintf("%s-%s", reflect.ValueOf(src).Type(), reflect.ValueOf(dst).Type())
		fmt.Println("Enter Map() ", typeConversion)
		defer fmt.Println("Leaving Map() ", typeConversion)
	}

	m.mapValue(reflect.ValueOf(src).Elem(), reflect.ValueOf(dst).Elem())
}

func (m *StructMapper) mapValue(src, dst reflect.Value) {
	if debug {
		fmt.Println("Enter mapValue for ", src, src.Type(), src.Kind(), dst, dst.Type(), dst.Kind())
		defer fmt.Println("Leave mapValue for ", src, src.Type(), src.Kind(), dst, dst.Type(), dst.Kind())
	}

	if src.Kind() == reflect.Map && dst.Kind() == reflect.Map {
		dst.Set(reflect.MakeMap(dst.Type()))
		for _, k := range src.MapKeys() {
			dstValue := reflect.New(dst.Type().Elem()).Elem()
			m.mapValue(src.MapIndex(k), dstValue)
			dst.SetMapIndex(k, dstValue)
		}
	} else if src.Kind() == reflect.Slice && dst.Kind() == reflect.Slice {
		dst.Set(reflect.MakeSlice(dst.Type(), src.Len(), src.Cap()))
		for i := 0; i < src.Len(); i++ {
			m.mapValue(src.Index(i), dst.Index(i))
		}
	} else if src.Kind() == reflect.Ptr && dst.Kind() == reflect.Ptr {
		if src.Elem().IsValid() {
			dst.Set(reflect.New(dst.Type()))
			m.mapValue(src.Elem(), dst.Elem())
		}
	} else if src.Kind() == reflect.Struct && dst.Kind() == reflect.Struct {

		// delegate to custom mapper if needed
		typeConversion := fmt.Sprintf("%s-%s", src.Type(), dst.Type())
		if cMapper, ok := m.customWrappers[typeConversion]; ok {
			reflect.ValueOf(cMapper).Call([]reflect.Value{src, dst.Addr()})
			return
		}

		baseFields := src
		targetFields := dst
		if m.mappingMode == DstFieldBased {
			baseFields = dst
			targetFields = src
		}

		for i := 0; i < baseFields.NumField(); i++ {
			srcFieldName := baseFields.Type().Field(i).Name
			foundField := false

			// find the field with the same name in target
			for j := 0; j < targetFields.NumField(); j++ {
				dstFieldName := targetFields.Type().Field(j).Name

				if dstFieldName == srcFieldName {
					foundField = true
					m.mapValue(src.Field(i), dst.Field(j))
				}
			}

			if !foundField && !m.ignoreUnknownField {
				panic("Could not find field " + srcFieldName + " in dst")
			}
		}
	} else if isPrimitive(src.Kind()) && isPrimitive(dst.Kind()) {
		// all primitive types remain, simply copy the value over
		dst.Set(src)
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