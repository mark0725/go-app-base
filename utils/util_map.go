package utils

import (
	"fmt"
	"reflect"
	"strings"
)

func MergeMaps(dst, src map[string]interface{}) map[string]interface{} {
	for key, srcValue := range src {
		if dstValue, ok := dst[key]; ok {
			// Both values are maps, need to merge them recursively
			dstMap, dstOk := dstValue.(map[string]interface{})
			srcMap, srcOk := srcValue.(map[string]interface{})
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

func MapToStruct(m map[string]interface{}, s interface{}) error {
	v := reflect.ValueOf(s).Elem()
	t := v.Type()

	// Check if s is a pointer to a struct
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("s must be a pointer to a struct")
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("json")
		tagOptions := strings.Split(tag, ",")

		// Determine the key in the map
		key := tagOptions[0]
		if key == "" {
			continue
		}

		// Set the struct field if the key exists in the map
		if val, ok := m[key]; ok {
			fieldValue := v.Field(i)
			if val == nil {
				continue
			}

			if !fieldValue.CanSet() {
				continue
			}
			val := reflect.ValueOf(val)

			// Handle different field types
			switch fieldValue.Kind() {
			case reflect.Ptr:
				if val.Kind() == reflect.Map && fieldValue.Type().Elem().Kind() == reflect.Struct {
					ptr := reflect.New(fieldValue.Type().Elem())
					err := MapToStruct(val.Interface().(map[string]interface{}), ptr.Interface())
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
					err := MapToStruct(val.Interface().(map[string]interface{}), fieldValue.Addr().Interface())
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
				slice := reflect.MakeSlice(fieldValue.Type(), val.Len(), val.Cap())
				for j := 0; j < val.Len(); j++ {
					elem := reflect.New(fieldValue.Type().Elem()).Elem()
					elem.Set(val.Index(j).Convert(fieldValue.Type().Elem()))
					slice.Index(j).Set(elem)
				}
				fieldValue.Set(slice)

			default:
				if val.Type().ConvertibleTo(fieldValue.Type()) {
					fieldValue.Set(val.Convert(fieldValue.Type()))
				} else {
					return fmt.Errorf("cannot convert value for field %s", field.Name)
				}
			}
		}
	}
	return nil
}
