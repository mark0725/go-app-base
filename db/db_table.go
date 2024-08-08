package db

import (
	"errors"
	"reflect"
)

type DBField struct {
	Id      string
	VarName string
	Type    string
	Len     string
	Name    string
	Comment string
}

// 解析结构体标签的函数
func GetTableFields(s interface{}) ([]DBField, error) {
	// 获取传入变量的类型
	var fields []DBField
	t := reflect.TypeOf(s)

	// 确保传入的是一个 struct 类型
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, errors.New("expected struct")
	}

	// 迭代结构体的所有字段
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		var dbField DBField

		if fieldTag := field.Tag.Get("field-id"); fieldTag == "" {
			continue
		} else {
			dbField.Id = fieldTag
		}

		if fieldTag := field.Tag.Get("field-name"); fieldTag != "" {
			dbField.Name = fieldTag
		}

		if fieldTag := field.Tag.Get("field-type"); fieldTag != "" {
			dbField.Type = fieldTag
		}

		if fieldTag := field.Tag.Get("field-len"); fieldTag != "" {
			dbField.Len = fieldTag
		}

		if fieldTag := field.Tag.Get("field-comment"); fieldTag != "" {
			dbField.Comment = fieldTag
		}

		dbField.VarName = field.Name

		fields = append(fields, dbField)

	}

	return fields, nil
}

func GetTableFieldIds(s interface{}) ([]string, error) {
	// 获取传入变量的类型
	var fields []string
	t := reflect.TypeOf(s)

	// 确保传入的是一个 struct 类型
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, errors.New("expected struct")
	}

	// 迭代结构体的所有字段
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if fieldTag := field.Tag.Get("field-id"); fieldTag != "" {
			fields = append(fields, fieldTag)
		}
	}

	return fields, nil
}
