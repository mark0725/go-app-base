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

func MapToStruct2(m map[string]any, s any) error {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("s must be a non-nil pointer to struct")
	}
	v = v.Elem()
	t := v.Type()
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("s must be a pointer to a struct")
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("json")
		tagOpts := strings.Split(tag, ",")
		key := tagOpts[0]
		if key == "" {
			key = field.Name
		}

		val, ok := m[key]
		if !ok || val == nil {
			continue
		}

		dst := v.Field(i)
		if !dst.CanSet() {
			continue
		}

		src := reflect.ValueOf(val)

		switch dst.Kind() {

		case reflect.Ptr:
			if src.Kind() == reflect.Map && dst.Type().Elem().Kind() == reflect.Struct {
				ptr := reflect.New(dst.Type().Elem())
				if err := MapToStruct(src.Interface().(map[string]any), ptr.Interface()); err != nil {
					return err
				}
				dst.Set(ptr)
			} else if src.Type().ConvertibleTo(dst.Type().Elem()) {
				ptr := reflect.New(dst.Type().Elem())
				ptr.Elem().Set(src.Convert(dst.Type().Elem()))
				dst.Set(ptr)
			} else {
				return fmt.Errorf("cannot convert value for field %s", field.Name)
			}

		case reflect.Struct:
			if src.Kind() != reflect.Map {
				return fmt.Errorf("expected map for nested field %s", field.Name)
			}
			if err := MapToStruct(src.Interface().(map[string]any), dst.Addr().Interface()); err != nil {
				return err
			}

		case reflect.Slice, reflect.Array:
			if src.Kind() != reflect.Slice && src.Kind() != reflect.Array {
				return fmt.Errorf("expected slice/array for field %s", field.Name)
			}
			elemTyp := dst.Type().Elem()
			n := src.Len()

			switch dst.Kind() {
			case reflect.Slice:
				out := reflect.MakeSlice(dst.Type(), n, n)
				if err := fillSequential(out, src, elemTyp, field.Name); err != nil {
					return err
				}
				dst.Set(out)

			case reflect.Array:
				// 数组长度固定，超出部分截断，缺失部分保持零值
				out := reflect.New(dst.Type()).Elem()
				if err := fillSequential(out, src, elemTyp, field.Name); err != nil {
					return err
				}
				dst.Set(out)
			}

		default: // 基础类型
			if src.Type().ConvertibleTo(dst.Type()) {
				dst.Set(src.Convert(dst.Type()))
			}
		}
	}
	return nil
}
func MapToStruct(m map[string]any, s any) error {
	rv := reflect.ValueOf(s)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("s must be a pointer to a struct")
	}

	v := rv.Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		fieldVal := v.Field(i)
		if !fieldVal.CanSet() {
			continue
		}

		tag := sf.Tag.Get("json")
		name := strings.Split(tag, ",")[0]
		if name == "" {
			name = sf.Name
		}
		raw, ok := m[name]
		if !ok || raw == nil {
			continue
		}

		srcVal := reflect.ValueOf(raw)

		//------------------------------------------------------------------
		// 指针
		//------------------------------------------------------------------
		switch fieldVal.Kind() {

		case reflect.Ptr:
			elemTyp := fieldVal.Type().Elem()
			if srcVal.Kind() == reflect.Map && elemTyp.Kind() == reflect.Struct {
				ptr := reflect.New(elemTyp)
				if err := MapToStruct(srcVal.Interface().(map[string]any), ptr.Interface()); err != nil {
					return err
				}
				fieldVal.Set(ptr)
				continue
			}
			if srcVal.Type().ConvertibleTo(elemTyp) {
				ptr := reflect.New(elemTyp)
				ptr.Elem().Set(srcVal.Convert(elemTyp))
				fieldVal.Set(ptr)
				continue
			}
			return fmt.Errorf("cannot convert %v to %v (field %s)",
				srcVal.Type(), elemTyp, sf.Name)

		//------------------------------------------------------------------
		// 结构体
		//------------------------------------------------------------------
		case reflect.Struct:
			if srcVal.Kind() != reflect.Map {
				return fmt.Errorf("expected map for nested field %s", sf.Name)
			}
			if err := MapToStruct(srcVal.Interface().(map[string]any), fieldVal.Addr().Interface()); err != nil {
				return err
			}

		//------------------------------------------------------------------
		// 切片
		//------------------------------------------------------------------
		case reflect.Slice:
			if srcVal.Kind() != reflect.Slice {
				return fmt.Errorf("expected slice for field %s", sf.Name)
			}

			l := srcVal.Len()
			dstSlice := reflect.MakeSlice(fieldVal.Type(), l, l)
			elemTyp := fieldVal.Type().Elem()

			for j := 0; j < l; j++ {
				srcElem := srcVal.Index(j)

				// 关键修复：把 interface{} 外壳全部剥掉
				for srcElem.Kind() == reflect.Interface {
					srcElem = srcElem.Elem()
				}

				dstElem := dstSlice.Index(j)

				if elemTyp.Kind() == reflect.Ptr {
					dstElem.Set(reflect.New(elemTyp.Elem()))
					dstElem = dstElem.Elem()
				}

				if dstElem.Kind() == reflect.Struct && srcElem.Kind() == reflect.Map {
					if err := MapToStruct(srcElem.Interface().(map[string]any), dstElem.Addr().Interface()); err != nil {
						return err
					}
					continue
				}

				if srcElem.Type().ConvertibleTo(dstElem.Type()) {
					dstElem.Set(srcElem.Convert(dstElem.Type()))
					continue
				}
				return fmt.Errorf("cannot convert slice elem for field %s", sf.Name)
			}
			fieldVal.Set(dstSlice)

		//------------------------------------------------------------------
		// 固定数组
		//------------------------------------------------------------------
		case reflect.Array:
			if srcVal.Kind() != reflect.Slice && srcVal.Kind() != reflect.Array {
				return fmt.Errorf("expected array/slice for field %s", sf.Name)
			}
			l := srcVal.Len()
			if l > fieldVal.Len() {
				l = fieldVal.Len()
			}
			for j := 0; j < l; j++ {
				srcElem := srcVal.Index(j)
				// 同样剥掉 interface 包装
				for srcElem.Kind() == reflect.Interface {
					srcElem = srcElem.Elem()
				}
				dstElem := fieldVal.Index(j)

				if dstElem.Kind() == reflect.Struct && srcElem.Kind() == reflect.Map {
					if err := MapToStruct(srcElem.Interface().(map[string]any), dstElem.Addr().Interface()); err != nil {
						return err
					}
					continue
				}
				if srcElem.Type().ConvertibleTo(dstElem.Type()) {
					dstElem.Set(srcElem.Convert(dstElem.Type()))
					continue
				}
				return fmt.Errorf("cannot convert array elem for field %s", sf.Name)
			}
		case reflect.Map:
			if srcVal.Kind() != reflect.Map {
				return fmt.Errorf("expected map for field %s", sf.Name)
			}

			dstMap := reflect.MakeMap(fieldVal.Type())
			keyTyp := fieldVal.Type().Key()
			valTyp := fieldVal.Type().Elem()

			for _, k := range srcVal.MapKeys() {
				srcElem := srcVal.MapIndex(k)
				for srcElem.Kind() == reflect.Interface {
					srcElem = srcElem.Elem()
				}

				if !k.Type().ConvertibleTo(keyTyp) {
					return fmt.Errorf("cannot convert map key %v to %v (field %s)",
						k.Type(), keyTyp, sf.Name)
				}
				dstKey := k.Convert(keyTyp)

				var dstVal reflect.Value

				switch {
				// value 是结构体 --------------------------------------------------
				case valTyp.Kind() == reflect.Struct && srcElem.Kind() == reflect.Map:
					dstVal = reflect.New(valTyp).Elem()
					if err := MapToStruct(srcElem.Interface().(map[string]any),
						dstVal.Addr().Interface()); err != nil {
						return err
					}

				// value 是 *结构体 ------------------------------------------------
				case valTyp.Kind() == reflect.Ptr &&
					valTyp.Elem().Kind() == reflect.Struct &&
					srcElem.Kind() == reflect.Map:

					dstVal = reflect.New(valTyp.Elem())
					if err := MapToStruct(srcElem.Interface().(map[string]any),
						dstVal.Interface()); err != nil {
						return err
					}

				// value 还是 Map，递归处理 ---------------------------------------
				case valTyp.Kind() == reflect.Map && srcElem.Kind() == reflect.Map:
					// 这里是关键修复
					nestedMap := reflect.MakeMap(valTyp)
					// 复用当前函数实现递归复制
					if err := copyMap(srcElem, nestedMap, sf.Name); err != nil {
						return err
					}
					dstVal = nestedMap

				// 其它可直接转换 --------------------------------------------------
				default:
					if !srcElem.Type().ConvertibleTo(valTyp) {
						return fmt.Errorf("cannot convert map value %v to %v (field %s)",
							srcElem.Type(), valTyp, sf.Name)
					}
					dstVal = srcElem.Convert(valTyp)
				}

				dstMap.SetMapIndex(dstKey, dstVal)
			}

			fieldVal.Set(dstMap)
		default:
			if srcVal.Type().ConvertibleTo(fieldVal.Type()) {
				fieldVal.Set(srcVal.Convert(fieldVal.Type()))
			} else {
				return fmt.Errorf("cannot convert %v to %v (field %s)",
					srcVal.Type(), fieldVal.Type(), sf.Name)
			}
		}
	}
	return nil
}

func copyMap(src, dst reflect.Value, fld string) error {
	keyTyp := dst.Type().Key()
	valTyp := dst.Type().Elem()

	for _, k := range src.MapKeys() {
		v := src.MapIndex(k)
		for v.Kind() == reflect.Interface {
			v = v.Elem()
		}

		if !k.Type().ConvertibleTo(keyTyp) {
			return fmt.Errorf("cannot convert map key %v to %v (field %s)",
				k.Type(), keyTyp, fld)
		}
		dstKey := k.Convert(keyTyp)

		var dstVal reflect.Value
		switch {
		case valTyp.Kind() == reflect.Map && v.Kind() == reflect.Map:
			nested := reflect.MakeMap(valTyp)
			if err := copyMap(v, nested, fld); err != nil {
				return err
			}
			dstVal = nested

		case valTyp.Kind() == reflect.Struct && v.Kind() == reflect.Map:
			dstVal = reflect.New(valTyp).Elem()
			if err := MapToStruct(v.Interface().(map[string]any),
				dstVal.Addr().Interface()); err != nil {
				return err
			}

		default:
			if !v.Type().ConvertibleTo(valTyp) {
				return fmt.Errorf("cannot convert map value %v to %v (field %s)",
					v.Type(), valTyp, fld)
			}
			dstVal = v.Convert(valTyp)
		}
		dst.SetMapIndex(dstKey, dstVal)
	}
	return nil
}

// fillSequential 将 src 中的元素依次写入 out，支持元素为结构体或基础类型。
// out 可以是 slice 也可以是 array。
func fillSequential(out, src reflect.Value, elemTyp reflect.Type, fieldName string) error {
	limit := min(out.Len(), src.Len())
	for j := 0; j < limit; j++ {
		dstElem := out.Index(j)
		srcElem := src.Index(j)

		// 如果目标是 *T，需要先创建指针
		if dstElem.Kind() == reflect.Ptr {
			ptr := reflect.New(elemTyp.Elem())
			dstElem.Set(ptr)
			dstElem = ptr.Elem()
		}

		switch dstElem.Kind() {
		case reflect.Struct:
			if srcElem.Kind() != reflect.Map {
				return fmt.Errorf("expected map in %s[%d]", fieldName, j)
			}
			if err := MapToStruct(srcElem.Interface().(map[string]any), dstElem.Addr().Interface()); err != nil {
				return err
			}
		default:
			if !srcElem.Type().ConvertibleTo(dstElem.Type()) {
				return fmt.Errorf("cannot convert %s[%d]", fieldName, j)
			}
			dstElem.Set(srcElem.Convert(dstElem.Type()))
		}
	}
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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
