/**
 * 感谢前辈们的开源，让自己得以学习
 * 摘抄至strcase库：https://github.com/iancoleman/strcase
 */
package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// 假设一个负载的字符串："AnyKind of_string"

// 转换为：any_kind_of_string
func StrToSnake(s string) string {
	return ToScreamingDelimited(s, '_', "", false)
}

// 转换为：any_kind.of_string
func StrToSnakeIgnore(s, ignore string) string {
	return ToScreamingDelimited(s, '_', ignore, true)
}

// 转换为：any-kind-of-string
func StrToKebab(s string) string {
	return ToScreamingDelimited(s, '-', "", false)
}

// 转换为：AnyKindOfString
func StrToCamel(s string) string {
	return ToCamelInitCase(s, true)
}

// 转换为：anyKindOfString  首字母小写
func StrToLowerCamel(s string) string {
	return ToCamelInitCase(s, false)
}

// 将斜号转为冒号：/admin/base  -->  admin:base
func StrToColonBase(s string) string {
	var ignore = []string{}
	return ToBevelColon(s, ignore)
}

/**
 * 将长字符串转为hash值
 */
func StrToHash(s string) string {
	hasher := sha256.New()
	hasher.Write([]byte(s))
	hashStr := hex.EncodeToString(hasher.Sum(nil))
	return hashStr
}
