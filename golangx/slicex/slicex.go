package slicex

import (
	"errors"
	"math/rand"
	"reflect"
)

var TypeError = errors.New("Parameter type error")
var LenthError = errors.New("Slice is empty")

// RemoveStringFields
/*
	从string类型slice中删除字段
*/
func RemoveStringFields(sli []string, opts ...string) []string {
	var out []string
	for _, v := range sli {
		found := false
		for _, opt := range opts {
			if v == opt {
				found = true
				break
			}
		}
		if !found {
			out = append(out, v)
		}
	}
	return out
}

// GetSliceValue
/*
	根据下标获取slice的元素
	return any, bool
*/
func GetSliceValue(slice interface{}, index int) (any, bool) {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return nil, false
	}
	for i := 0; i < s.Len(); i++ {
		if i == index {
			return s.Index(i).Interface(), true
		}
	}
	return nil, false
}

// GetSliceValueWithRand
/*
	随机获取slice的某一个元素
*/
func GetSliceValueWithRand(slice interface{}) (any, error) {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return nil, TypeError
	}
	if s.Len() == 0 {
		return nil, LenthError
	}
	return s.Index(rand.Intn(s.Len())).Interface(), nil
}
