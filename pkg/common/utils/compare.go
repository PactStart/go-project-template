package utils

import (
	"fmt"
	"reflect"
	"strings"
)

type Diff struct {
	FieldName string      `json:"fieldName"`
	OldValue  interface{} `json:"oldValue"`
	NewValue  interface{} `json:"newValue"`
}

func CompareStructs(a, b interface{}, ignoreFields ...string) ([]Diff, error) {
	diffs := []Diff{}
	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)

	if va.Type() != vb.Type() {
		return nil, fmt.Errorf("cannot compare different types")
	}

	for i := 0; i < va.NumField(); i++ {
		field := va.Type().Field(i)
		name := field.Name

		// 检查当前字段是否应该忽略
		ignore := false
		for _, f := range ignoreFields {
			if f == name {
				ignore = true
				break
			}
		}
		if ignore {
			continue
		}

		valueA := va.Field(i).Interface()
		valueB := vb.Field(i).Interface()

		if !reflect.DeepEqual(valueA, valueB) {
			diffs = append(diffs, Diff{
				FieldName: name,
				OldValue:  valueA,
				NewValue:  valueB,
			})
		}
	}
	return diffs, nil
}

func CompareSpecifiedColumns(a, b interface{}, columns ...string) ([]Diff, error) {
	diffs := []Diff{}
	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)

	if va.Type() != vb.Type() {
		return nil, fmt.Errorf("cannot compare different types")
	}

	typeMap := make(map[string]string)
	for i := 0; i < va.NumField(); i++ {
		field := va.Type().Field(i)
		gormTag := field.Tag.Get("gorm")
		if gormTag == "" {
			continue
		}
		tags := strings.Split(gormTag, ";")
		for _, tag := range tags {
			if strings.HasPrefix(tag, "column:") {
				columnName := strings.TrimPrefix(tag, "column:")
				typeMap[columnName] = field.Name
				break
			}
		}
	}

	for _, col := range columns {
		fieldName, ok := typeMap[col]
		if !ok {
			continue // 如果列名不在结构体中定义，则跳过
		}
		fieldA := va.FieldByName(fieldName).Interface()
		fieldB := vb.FieldByName(fieldName).Interface()
		if !reflect.DeepEqual(fieldA, fieldB) {
			diffs = append(diffs, Diff{
				FieldName: fieldName,
				OldValue:  fieldA,
				NewValue:  fieldB,
			})
		}
	}
	return diffs, nil
}
