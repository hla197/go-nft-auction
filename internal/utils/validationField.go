package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// FieldValidate 接口
type FieldValidateIF interface {
	Validate(err error, targetStruct interface{}) string
}
type FieldValidate struct {
}

func (f *FieldValidate) Validate(err error, targetStruct interface{}) string {
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			return formatValidationError(e, targetStruct)
		}
	}
	return "请求参数无效"
}

func formatValidationError(e validator.FieldError, obj interface{}) string {
	// 使用反射获取字段信息
	fieldType := reflect.TypeOf(obj)
	structField, found := fieldType.FieldByName(e.Field())

	var fieldName string

	if found {
		// 优先从 tag 中获取 label
		label := structField.Tag.Get("label")
		if label != "" {
			fieldName = label
		} else {
			// 如果没有 label，尝试从 json tag 获取并美化
			jsonTag := structField.Tag.Get("json")
			if jsonTag != "" {
				fieldName = strings.Split(jsonTag, ",")[0]
				fieldName = strings.Title(fieldName) // 首字母大写
			} else {
				fieldName = e.Field()
			}
		}
	} else {
		fieldName = e.Field()
	}

	// 4. 根据 Tag 生成提示（这部分逻辑保持不变）
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s 不能为空", fieldName)
	case "email":
		return fmt.Sprintf("%s 格式不正确", fieldName)
	case "min":
		return fmt.Sprintf("%s 长度不能少于 %s 位", fieldName, e.Param())
	case "max":
		return fmt.Sprintf("%s 长度不能超过 %s 位", fieldName, e.Param())
	case "eqfield":
		// 处理确认密码等字段
		// 获取比较字段的 Label
		compareField := e.Param()
		compareStructField, _ := fieldType.FieldByName(compareField)
		compareLabel := compareStructField.Tag.Get("label")
		if compareLabel == "" {
			compareLabel = compareField
		}
		return fmt.Sprintf("%s 与 %s 不一致", fieldName, compareLabel)
	default:
		return fmt.Sprintf("%s 参数校验失败", fieldName)
	}
}
