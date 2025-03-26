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

	// Handle pointer to struct
	if objValue.Kind() == reflect.Ptr {
		objValue = objValue.Elem()
	}

	if objValue.Kind() != reflect.Struct {
		fmt.Println("Input is not a struct")
		return result
	}

	objType := objValue.Type()

	for i := 0; i < objValue.NumField(); i++ {
		fieldValue := objValue.Field(i)
		fieldType := objType.Field(i)
		tag := fieldType.Tag.Get("json")
		tagOptions := strings.Split(tag, ",")
		key := tagOptions[0]
		//omitempty := tagOptions[1]
		if key == "" {
			continue
		}

		fieldName := key

		// Process the field value depending on its kind.
		switch fieldValue.Kind() {
		case reflect.Ptr:
			if !fieldValue.IsNil() {
				result[fieldName] = fieldValue.Elem().Interface()
			}
		case reflect.Struct:
			result[fieldName] = StructToMap(fieldValue.Interface())
		case reflect.Slice:
			sliceLen := fieldValue.Len()
			sliceResult := make([]any, sliceLen)
			for j := 0; j < sliceLen; j++ {
				sliceElem := fieldValue.Index(j).Interface()
				sliceResult[j] = sliceElem
			}
			result[fieldName] = sliceResult
		default:
			result[fieldName] = fieldValue.Interface()
		}
	}

	return result
}
