package before

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var Valid *validator.Validate

func RegisterValidate() {
	validate := validator.New()
	validate.RegisterValidation("phone", validatePhone)
	validate.RegisterValidation("chinese", validateChinese)
	validate.RegisterValidation("letter", validateLetter)
	validate.RegisterValidation("stand", validateStand)

	Valid = validate
	fmt.Println("validate.register.success")
}

// 手机号验证
func validatePhone(f validator.FieldLevel) bool {
	field := f.Field().String()
	if len(field) == 0 {
		return true // 为空时不验证
	}
	if len(field) != 11 {
		return false
	}
	reg := regexp.MustCompile(`^1[3-9]\d{9}$`)

	if reg.MatchString(field) {
		return true
	} else {
		return false
	}
}

// 汉字验证
func validateChinese(f validator.FieldLevel) bool {
	field := f.Field().String()
	if len(field) == 0 {
		return true // 为空时不验证
	}
	// 等效于 ^[\u4e00-\u9fa5]+$
	reg := regexp.MustCompile(`^[\x{4e00}-\x{9fa5}]+$`)

	if reg.MatchString(field) {
		return true
	} else {
		return false
	}
}

// 字母验证
func validateLetter(f validator.FieldLevel) bool {
	field := f.Field().String()
	if len(field) == 0 {
		return true // 为空时不验证
	}
	reg := regexp.MustCompile(`^[A-Za-z]+$`)

	if reg.MatchString(field) {
		return true
	} else {
		return false
	}
}

// 标准字符，允许：字母、汉字、中划线、下划线、数字、斜杠、竖线
func validateStand(f validator.FieldLevel) bool {
	field := f.Field().String()
	if len(field) == 0 {
		return true // 为空时不验证
	}
	reg := regexp.MustCompile(`^[\w\x{4e00}-\x{9fa5}-0-9\\/|]+$`)

	if reg.MatchString(field) {
		return true
	} else {
		return false
	}
}
