package util

import (
	"reflect"
	"strings"
)

func GetGormFields(stc interface{}) []string {
	value := reflect.ValueOf(stc)
	typ := value.Type()
	if typ.Kind() == reflect.Ptr {
		if value.IsNil() {
			return nil
		} else {
			typ = typ.Elem()
			value = value.Elem()
		}
	}
	if typ.Kind() == reflect.Struct {
		columns := make([]string, 0, value.NumField())
		for i := 0; i < value.NumField(); i++ {
			fieldType := typ.Field(i)
			if fieldType.IsExported() {
				if fieldType.Tag.Get("gorm") == "-" {
					continue
				}
				name := Camel2Snake(fieldType.Name)
				if len(fieldType.Tag.Get("gorm")) > 0 {
					content := fieldType.Tag.Get("gorm")
					if strings.HasPrefix(content, "column:") {
						content = content[7:]
						pos := strings.Index(content, ";")
						if pos > 0 {
							name = content[0:pos]
						} else if pos < 0 {
							name = content
						}
					}
				}
				columns = append(columns, name)
			}
		}
		return columns
	} else {
		return nil
	}
}
