package utils

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"path/filepath"
)

/**
 * 递归函数
 * @Method Recursion
 * @Param
 */
func Recursion(function func() int) error {
	if function() > 0 {
		Recursion(function)
	} else {
		return nil
	}

	return nil
}

/**
 * 获取根目录下文件
 */
func PathResolve(url ...string) string {
	root, err := os.Getwd()
	if err != nil {
		return ""
	}

	path := filepath.Join(append([]string{root}, url...)...)
	return path
}

/**
 * 创建MD5加密
 */
func MD5Encode(val string) string {
	md := md5.New()
	md.Write([]byte(val))

	return hex.EncodeToString(md.Sum(nil))
}

/**
 * 比对信息和加密信息是否一致
 */
func MD5Verify(sourceVal, encryptVal string) bool {
	encrypt := MD5Encode(sourceVal)

	return encryptVal == encrypt
}
