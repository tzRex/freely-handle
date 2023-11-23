package utils

import (
	"errors"
	"fmt"
	"reflect"
)

type StructListDiffResult struct {
	Add     []interface{} // 新增项
	Del     []interface{} // 删除项
	Mod     []interface{} // 变更项
	ModToId []interface{} // 附带原始唯一值的变更项
	Same    []interface{} // 相同项
}

/**
 * 相似结构体中相同类型的值进行赋值
 * @method StructCopy
 * @param {[type]} interface{} src 源结构体
 * @param {[type]} interface{} dst 需要被赋值的结构体
 */
func StructCopy(src, dst interface{}) error {
	srcV, err := srcFilter(src)
	if err != nil {
		return err
	}
	dstV, err := dstFilter(dst)
	if err != nil {
		return err
	}
	srcKeys := make(map[string]bool)
	for i := 0; i < srcV.NumField(); i++ {
		srcKeys[srcV.Type().Field(i).Name] = true
	}
	for i := 0; i < dstV.Elem().NumField(); i++ {
		fName := dstV.Elem().Type().Field(i).Name
		if _, ok := srcKeys[fName]; ok {
			v := srcV.FieldByName(fName)
			dstField := dstV.Elem().Field(i)
			if v.CanInterface() {
				if dstField.Kind() != reflect.Ptr {
					dstV.Elem().Field(i).Set(v)
				}
			}
		}
	}

	return nil
}

func srcFilter(src interface{}) (reflect.Value, error) {
	v := reflect.ValueOf(src)
	if v.Type().Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return reflect.Zero(v.Type()), errors.New("src type error: not a struct or a pointer to struct")
	}
	return v, nil
}

func dstFilter(src interface{}) (reflect.Value, error) {
	v := reflect.ValueOf(src)
	if v.Type().Kind() != reflect.Ptr {
		return reflect.Zero(v.Type()), errors.New("dst type error: not a pointer to struct")
	}
	if v.Elem().Kind() != reflect.Struct {
		return reflect.Zero(v.Type()), errors.New("dst type error: not point to struct")
	}
	return v, nil
}

/**
 * 将结构体转为map对象，用于省略字段使用
 * @method StructToMap
 * @param {[type]} interface{} src 源结构体
 * @param {[type]} bool filterZero 是否过滤空值
 * @param {[type]} bool isSnake 是否转为蛇形字符串
 */
func StructToMap(src interface{}, filterZero bool, isSnake bool) (map[string]interface{}, error) {
	val := reflect.ValueOf(src)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil, fmt.Errorf("point is nil")
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("src.type.not.struct")
	}

	valT := val.Type()
	resultMap := make(map[string]interface{})

	for i := 0; i < valT.NumField(); i++ {
		field := val.Field(i)
		if !field.CanInterface() {
			return nil, fmt.Errorf("field %q fial", valT.Field(i).Name)
		}

		var keyName string
		if isSnake {
			keyName = StrToSnake(valT.Field(i).Name)
		} else {
			keyName = valT.Field(i).Name
		}

		if field.Kind() == reflect.Struct && valT.Field(i).Anonymous {
			subStu, err := StructToMap(field.Interface(), filterZero, isSnake)
			if err != nil {
				return nil, err
			}
			for k, v := range subStu {
				resultMap[k] = v
			}
		} else {
			if filterZero {
				if !field.IsZero() {
					resultMap[keyName] = field.Interface()
				}
			} else {
				resultMap[keyName] = field.Interface()
			}
		}
	}

	return resultMap, nil
}

/**
 * 将结构体中的字段删除并返回新的结构体
 * @method StructFilter
 * @param {[type]} interface{} src 源结构体
 * @param {[type]} []string ignores 需要忽略的字段
 */
func StructFilter(src interface{}, ignores []string) (interface{}, error) {
	colT := reflect.TypeOf(src)

	if colT.Kind() != reflect.Struct {
		return nil, fmt.Errorf("src.must.be.struct")
	}

	length := colT.NumField()
	newFields := make([]reflect.StructField, length)

	for i := 0; i < length; i++ {
		field := colT.Field(i)
		if IndexOfToString(ignores, field.Name) != -1 {
			newFields[i] = field
		}
	}

	newType := reflect.StructOf(newFields)
	newVal := reflect.New(newType).Elem()

	for i := 0; i < newType.NumField(); i++ {
		newCol := newType.Field(i)
		colV, ok := colT.FieldByName(newCol.Name)
		if !ok {
			return nil, fmt.Errorf("field %q not found", newCol.Name)
		}
		newVal.Field(i).Set(reflect.ValueOf(reflect.ValueOf(src).FieldByIndex(colV.Index).Interface()))
	}

	return newVal.Addr().Interface(), nil
}

/**
 * 获取结构体中的属性
 * @method StructGetField
 * @param {[type]} interface{} src 源结构体
 * @param {[type]} string key 需要查找的字段
 *
 */
func StructGetField(src interface{}, key string) (interface{}, error) {
	val, err := srcFilter(src)
	if err != nil {
		return nil, err
	}

	field := val.FieldByName(key)
	if !field.IsValid() {
		return nil, errors.New("field.not.found")
	}

	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return nil, errors.New("field.is.nil")
		}
		field = field.Elem()
	}

	return field.Interface(), nil
}
