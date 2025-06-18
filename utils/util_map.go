package utils

import (
	"fmt"
	"reflect"
	"strings"
)

func MergeMaps(dst map[string]any, src map[string]any) map[string]any {
	for key, srcValue := range src {
		if dstValue, ok := dst[key]; ok {
			// Both values are maps, need to merge them recursively
			dstMap, dstOk := dstValue.(map[string]any)
			srcMap, srcOk := srcValue.(map[string]any)
			if dstOk && srcOk {
				dst[key] = MergeMaps(dstMap, srcMap)
			} else {
				// If they are not both maps, overwrite the destination value with the source value
				dst[key] = srcValue
			}
		} else {
			// Key not present in destination map, simply add it
			dst[key] = srcValue
		}
	}
	return dst
}

func DeepMerge(maps ...map[string]any) map[string]any {
	result := make(map[string]any)
	for _, m := range maps {
		for k, v := range m {
			if existing, ok := result[k]; ok {
				if existingMap, ok := existing.(map[string]any); ok {
					if vMap, ok := v.(map[string]any); ok {
						result[k] = DeepMerge(existingMap, vMap)
						continue
					}
				}
			}
			result[k] = v
		}
	}
	return result
}

func MapToStruct(m map[string]any, s any) error {
	v := reflect.ValueOf(s).Elem()
	t := v.Type()

	// Check if s is a pointer to a struct.
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("s must be a pointer to a struct")
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("json")
		tagOptions := strings.Split(tag, ",")

		// Determine the key in the map.
		key := tagOptions[0]
		if key == "" {
			key = field.Name // Fallback to field name if no tag is present.
		}

		// Set the struct field if the key exists in the map.
		if val, ok := m[key]; ok {
			fieldValue := v.Field(i)
			if val == nil {
				continue
			}

			if !fieldValue.CanSet() {
				continue
			}

			val := reflect.ValueOf(val)

			// Handle different field types.
			switch fieldValue.Kind() {
			case reflect.Ptr:
				if val.Kind() == reflect.Map && fieldValue.Type().Elem().Kind() == reflect.Struct {
					ptr := reflect.New(fieldValue.Type().Elem())
					err := MapToStruct(val.Interface().(map[string]any), ptr.Interface())
					if err != nil {
						return err
					}
					fieldValue.Set(ptr)
				} else if val.Type().ConvertibleTo(fieldValue.Type().Elem()) {
					fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
					fieldValue.Elem().Set(val.Convert(fieldValue.Type().Elem()))
				} else {
					return fmt.Errorf("cannot convert value for field %s", field.Name)
				}

			case reflect.Struct:
				if val.Kind() == reflect.Map {
					err := MapToStruct(val.Interface().(map[string]any), fieldValue.Addr().Interface())
					if err != nil {
						return err
					}
				} else {
					return fmt.Errorf("expected map for nested field %s", field.Name)
				}

			case reflect.Slice:
				if val.Kind() != reflect.Slice {
					continue
					//return fmt.Errorf("expected slice for field %s type %s", field.Name, val.Kind())
				}
				slice := reflect.MakeSlice(fieldValue.Type(), val.Len(), val.Len())
				for j := 0; j < val.Len(); j++ {
					elem := slice.Index(j)

					if elem.Kind() == reflect.Ptr {
						elem = reflect.New(fieldValue.Type().Elem().Elem())
						slice.Index(j).Set(elem)
						elem = elem.Elem()
					}

					if elem.Kind() == reflect.Struct && val.Index(j).Kind() == reflect.Map {
						err := MapToStruct(val.Index(j).Interface().(map[string]any), elem.Addr().Interface())
						if err != nil {
							return err
						}
					} else if val.Index(j).Type().ConvertibleTo(elem.Type()) {
						elem.Set(val.Index(j).Convert(elem.Type()))
					} else {
						break
						//return fmt.Errorf("cannot convert slice value for field %s", field.Name)
					}
				}
				fieldValue.Set(slice)

			default:
				if val.Type().ConvertibleTo(fieldValue.Type()) {
					fieldValue.Set(val.Convert(fieldValue.Type()))
				} else {
					continue
					//return fmt.Errorf("cannot convert value for field %s", field.Name)
				}
			}
		}
	}
	return nil
}
func StructToMap(obj any) map[string]any {
	result := make(map[string]any)
	objValue := reflect.ValueOf(obj)

	// 如果是指针，则取其所指向的内容
	if objValue.Kind() == reflect.Ptr {
		objValue = objValue.Elem()
	}

	// 只有结构体才进行处理，否则直接返回空 map
	if objValue.Kind() != reflect.Struct {
		return result
	}

	objType := objValue.Type()

	for i := 0; i < objValue.NumField(); i++ {
		fieldValue := objValue.Field(i)
		fieldType := objType.Field(i)

		// 跳过未导出字段
		if fieldType.PkgPath != "" {
			continue
		}

		// 解析 json tag
		tag := fieldType.Tag.Get("json")
		tagParts := strings.Split(tag, ",")
		tagKey := tagParts[0]
		omitempty := false
		for _, part := range tagParts[1:] {
			if part == "omitempty" {
				omitempty = true
			}
		}

		// 忽略 "-"
		if tagKey == "-" {
			continue
		}

		fieldName := tagKey
		if fieldName == "" {
			fieldName = fieldType.Name
		}

		// 如果有 omitempty 并且字段是零值，则跳过
		if omitempty && isZero(fieldValue) {
			continue
		}

		// 匿名字段嵌入体特殊处理
		if fieldType.Anonymous && fieldValue.Kind() == reflect.Struct {
			embedMap := StructToMap(fieldValue.Interface())
			for k, v := range embedMap {
				result[k] = v
			}
			continue
		}

		switch fieldValue.Kind() {
		case reflect.Ptr:
			if !fieldValue.IsNil() {
				elem := fieldValue.Elem()
				result[fieldName] = handleValue(elem, omitempty)
			} else if !omitempty {
				// 显示为nil，仅在 non-omitempty 时
				result[fieldName] = nil
			}
		case reflect.Struct, reflect.Slice, reflect.Array, reflect.Map, reflect.Interface:
			result[fieldName] = handleValue(fieldValue, omitempty)
		default:
			result[fieldName] = fieldValue.Interface()
		}
	}

	return result
}

// handleValue 递归处理各种类型（struct、slice、array、map、指针）
func handleValue(v reflect.Value, omitempty bool) any {
	if !v.IsValid() {
		return nil
	}

	switch v.Kind() {
	case reflect.Struct:
		// 结构体：再次调用 StructToMap
		return StructToMap(v.Interface())

	case reflect.Slice, reflect.Array:
		// 切片或数组
		if v.Len() == 0 && omitempty {
			return nil
		}
		sliceLen := v.Len()
		sliceResult := make([]any, sliceLen)
		for i := 0; i < sliceLen; i++ {
			elemValue := v.Index(i)
			sliceResult[i] = handleValue(elemValue, omitempty)
		}
		return sliceResult

	case reflect.Map:
		// map 类型
		if v.Len() == 0 && omitempty {
			return nil
		}
		mapResult := make(map[string]any, v.Len())
		for _, key := range v.MapKeys() {
			val := v.MapIndex(key)
			mapResult[key.String()] = handleValue(val, omitempty)
		}
		return mapResult

	case reflect.Ptr:
		if v.IsNil() {
			if !omitempty {
				return nil
			}
			return nil
		}
		elem := v.Elem()
		return handleValue(elem, omitempty)

	case reflect.Interface:
		// interface 类型，取其底层元素并再次处理
		if v.IsNil() {
			if !omitempty {
				return nil
			}
			return nil
		}
		elem := v.Elem()
		return handleValue(elem, omitempty)

	default:
		// 其他基础类型直接返回
		return v.Interface()
	}
}

// 判断 reflect.Value 是否为零值
func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return v.IsNil()
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if !isZero(v.Field(i)) {
				return false
			}
		}
		return true
	default:
		return false
	}
}
