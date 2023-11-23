package corebase

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/tzRex/freely-handle/before"
)

/**
 * 验证结构体内的字段数据是否符合要求
 * @Param stu 需要验证的结构体
 * @Any errInfo 在结构体的tag中定义的字段
 */
func ValidStuct(stu interface{}) error {
	err := before.Valid.Struct(stu)
	if err != nil {
		return processErr(stu, err)
	}
	return nil
}

/**
 * 自定义错误信息
 * @Param stu interface{} 对应的结构体或结构体指针
 */
func processErr(stu interface{}, err error) error {
	if err == nil {
		return nil
	}

	invalid, ok := err.(*validator.InvalidValidationError)
	if !ok {
		validErr := err.(validator.ValidationErrors)
		rType := reflect.TypeOf(stu)

		if rType.Kind() == reflect.Ptr {
			rType = rType.Elem()
		}

		for _, err := range validErr {
			fieldName := err.Field()
			field, ok := rType.FieldByName(fieldName)

			if ok {
				errInfo := field.Tag.Get("errInfo")
				if errInfo != "" {
					return fmt.Errorf(errInfo)
				} else {
					return err
				}
			} else {
				return err
			}
		}
	}
	return invalid
}
