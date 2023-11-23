package utils

import (
	"errors"
	"fmt"
	"reflect"
)

// 查找uint数组，找到了返回对应下标，没有则返回 -1
func IndexOfToUint(list []uint, target uint) int {
	for i, v := range list {
		if v == target {
			return i
		}
	}
	return -1
}

// 查找int数组，找到了返回对应下标，没有则返回 -1
func IndexOfToInt(list []int, target int) int {
	for i, v := range list {
		if v == target {
			return i
		}
	}
	return -1
}

// 查找字符串数组，找到了返回对应下标，没有则返回 -1
func IndexOfToString(list []string, target string) int {
	for i, v := range list {
		if v == target {
			return i
		}
	}
	return -1
}

/**
 * 获取对应数组的合适项的下标
 * @Method ListIIndexOf
 * @Param src interface{} 数组，可以是指针类型
 * @Param function func(i interface{}) bool 查询条件
 */
func ListFindIndex(src interface{}, function func(i interface{}) bool) int {
	srcV := reflect.ValueOf(src)
	if srcV.Type().Kind() == reflect.Ptr {
		srcV = srcV.Elem()
	}

	if srcV.Kind() != reflect.Slice {
		return -1
	}

	for i := 0; i < srcV.Len(); i++ {
		item := srcV.Index(i).Interface()
		if function(item) {
			return i
		}
	}

	return -1
}

/**
 * 对比两个结构体切片的差异，将变更的部分标记出来。（用于gorm中间表的处理）
 * @param src interface{} 旧的切片
 * @param src interface{} 新的切片
 * @param compareKeys []string 数据变更项的key
 * @param only string 新值需要附带的旧值的唯一字段，例如：ID
 */
func ListToDiff(src interface{}, dst interface{}, compareKeys []string, only string) (StructListDiffResult, error) {
	var result = StructListDiffResult{
		Add:     []interface{}{},
		Mod:     []interface{}{},
		ModToId: []interface{}{},
		Del:     []interface{}{},
		Same:    []interface{}{},
	}

	srcV := reflect.ValueOf(src)
	if srcV.Type().Kind() == reflect.Ptr {
		srcV = srcV.Elem()
	}

	dstV := reflect.ValueOf(dst)
	if dstV.Type().Kind() == reflect.Ptr {
		dstV = dstV.Elem()
	}

	if srcV.Kind() != reflect.Slice {
		return result, errors.New("src.type.err")
	}

	if dstV.Kind() != reflect.Slice {
		return result, errors.New("dst.type.err")
	}

	srcLen := srcV.Len()
	dstLen := dstV.Len()

	// 获取对应的hash值
	getHashKey := func(stu reflect.Value) string {
		key := ""
		if stu.Kind() == reflect.Ptr {
			stu = stu.Elem()
		}
		for _, col := range compareKeys {
			key += fmt.Sprintf("%v:", stu.FieldByName(col).Interface())
		}
		return key
	}

	index := 0

	sameOperate := func(srcI, dstI interface{}, hashSrcI, hashDstI string) {

		if hashSrcI != hashDstI {
			itV := dstV.Index(index)
			if itV.Kind() == reflect.Ptr {
				itV = itV.Elem()
			}

			sV := srcV.Index(index)
			if sV.Kind() == reflect.Ptr {
				sV = sV.Elem()
			}
			result.Mod = append(result.Mod, dstI)
			if only != "" && itV.IsValid() && itV.CanSet() {
				itV.FieldByName(only).Set(sV.FieldByName(only))
				result.ModToId = append(result.ModToId, itV.Interface())
			}
		} else {
			result.Same = append(result.Same, srcI)
		}
	}

	// 当传入值的长度大于现有值时，说明有新增；反之说明有删减项
	if dstLen >= srcLen {
		for index < dstLen {
			nowItem := dstV.Index(index).Interface()
			nowHash := getHashKey(dstV.Index(index))
			if index < srcLen {
				oldItem := srcV.Index(index).Interface()
				oldHash := getHashKey(srcV.Index(index))
				sameOperate(oldItem, nowItem, oldHash, nowHash)
			} else {
				result.Add = append(result.Add, nowItem)
			}

			index++
		}
	} else {
		for index < srcLen {
			oldItem := srcV.Index(index).Interface()
			oldHash := getHashKey(srcV.Index(index))
			if index < dstLen {
				nowItem := dstV.Index(index).Interface()
				nowHash := getHashKey(dstV.Index(index))
				sameOperate(oldItem, nowItem, oldHash, nowHash)
			} else {
				result.Del = append(result.Del, oldItem)
			}

			index++
		}
	}

	return result, nil
}

/**
 * 过滤数组（切片）
 * @Method ListFiltter 返回过滤后的切片
 * @Param slice 切片，可以是指针类型
 * @param function 如何过滤的具体操作
 */
func ListFiltter(slice interface{}, function func(interface{}) bool) interface{} {
	sV := reflect.ValueOf(slice)

	if sV.Type().Kind() == reflect.Ptr {
		sV = sV.Elem()
	}
	if sV.Kind() != reflect.Slice {
		return nil
	}

	len := sV.Len()

	result := reflect.MakeSlice(sV.Type(), 0, len)

	for i := 0; i < len; i++ {
		svI := sV.Index(i)
		if function(svI.Interface()) {
			result = reflect.Append(result, svI)
		}
	}

	return result.Interface()
}
