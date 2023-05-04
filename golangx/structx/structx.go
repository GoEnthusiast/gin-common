package structx

import (
	"errors"
	"fmt"
	"reflect"
)

// GetFields
/*
	返回结构体的全部字段
*/
func GetFields(i interface{}) ([]string, error) {
	v := reflect.ValueOf(i)
	t := v.Type()

	if t.Kind() == reflect.Ptr {
		v = v.Elem()
		t = v.Type()
	}

	if t.Kind() != reflect.Struct {
		return nil, errors.New("input parameter is not a struct")
	}

	fields := make([]string, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		fields = append(fields, t.Field(i).Name)
	}
	return fields, nil
}

// GetTagValues
/*
	指定结构体的tag, 获取对应tag值
*/
func GetTagValues(i interface{}, tagStr string) ([]string, error) {
	v := reflect.ValueOf(i)
	t := v.Type()

	if t.Kind() == reflect.Ptr {
		v = v.Elem()
		t = v.Type()
	}

	if t.Kind() != reflect.Struct {
		return nil, errors.New("input parameter is not a struct")
	}

	fields := make([]string, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get(tagStr)
		if tag != "" {
			fields = append(fields, fmt.Sprintf("`%s`", tag))
		}
	}
	return fields, nil
}
