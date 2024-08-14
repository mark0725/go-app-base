package options

import (
	"errors"
	"reflect"
)

type Option func(*map[string]interface{})

func With(key string, value interface{}) Option {
	return func(m *map[string]interface{}) {
		(*m)[key] = value
	}
}

func WithMap(m map[string]interface{}) Option {
	return func(m2 *map[string]interface{}) {
		for k, v := range m {
			(*m2)[k] = v
		}
	}
}

func WithMapMerge(m map[string]interface{}) Option {
	return func(m2 *map[string]interface{}) {
		for k, v := range m {
			if _, ok := (*m2)[k]; !ok {
				(*m2)[k] = v
			}
		}
	}
}

func createOptions(m map[string]interface{}, opts ...Option) {
	for _, opt := range opts {
		opt(&m)
	}
}

func NewMapOptions(opts ...Option) map[string]interface{} {
	m := make(map[string]interface{})
	createOptions(m, opts...)
	return m
}

func NewOptions(v interface{}, opts ...Option) error {

	if reflect.TypeOf(v) == reflect.TypeOf(map[string]interface{}{}) {
		createOptions(v.(map[string]interface{}), opts...)
		return nil
	}

	m := NewMapOptions(opts...)
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		return errors.New("v must be a pointer")
	}

	//判断v是否为struct
	if reflect.TypeOf(v).Elem().Kind() != reflect.Struct {
		return errors.New("v must be a struct")
	}

	val := reflect.ValueOf(v).Elem()
	for key, value := range m {
		field := val.FieldByName(key)
		if field.IsValid() && field.CanSet() {
			fieldValue := reflect.ValueOf(value)
			if fieldValue.Type().AssignableTo(field.Type()) {
				field.Set(fieldValue)
			}
		}
	}

	return nil
}
