package util

import (
	"math/rand"
	"reflect"
	"regexp"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"welfare-sign/internal/pkg/config"
)

// StructCopy 结构体赋值
func StructCopy(dst interface{}, src interface{}) error {
	srcv := reflect.ValueOf(src)
	dstv := reflect.ValueOf(dst)
	srct := reflect.TypeOf(src)
	dstt := reflect.TypeOf(dst)
	if srct.Kind() != reflect.Ptr || dstt.Kind() != reflect.Ptr ||
		srct.Elem().Kind() == reflect.Ptr || dstt.Elem().Kind() == reflect.Ptr {
		return errors.New("StructCopy: type of parameters must be ptr of value")
	}
	if srcv.IsNil() || dstv.IsNil() {
		return errors.New("StructCopy: value of parameters should not be nil")
	}
	srcV := srcv.Elem()
	dstV := dstv.Elem()
	srcfields := DeepFields(reflect.ValueOf(src).Elem().Type())
	for _, v := range srcfields {
		if v.Anonymous {
			continue
		}
		dst := dstV.FieldByName(v.Name)
		src := srcV.FieldByName(v.Name)
		if !dst.IsValid() {
			continue
		}
		if src.Type() == dst.Type() && dst.CanSet() {
			dst.Set(src)
			continue
		}
		if src.Kind() == reflect.Ptr && !src.IsNil() && src.Type().Elem() == dst.Type() {
			dst.Set(src.Elem())
			continue
		}
		if dst.Kind() == reflect.Ptr && dst.Type().Elem() == src.Type() {
			dst.Set(reflect.New(src.Type()))
			dst.Elem().Set(src)
			continue
		}
	}
	return nil
}

// DeepFields .
func DeepFields(eleType reflect.Type) []reflect.StructField {
	var fields []reflect.StructField

	for i := 0; i < eleType.NumField(); i++ {
		v := eleType.Field(i)
		if v.Anonymous && v.Type.Kind() == reflect.Struct {
			fields = append(fields, DeepFields(v.Type)...)
		} else {
			fields = append(fields, v)
		}
	}

	return fields
}

// IsMobile 正则判断是否是手机号
func IsMobile(mobile string) bool {
	reg := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	compile := regexp.MustCompile(reg)
	return compile.MatchString(mobile)
}

// GenerateCode 生成验证码随机数
func GenerateCode() (code string) {
	s := rand.NewSource(time.Now().UnixNano())
	for i := 0; i < viper.GetInt(config.KeyYuanpianLength); i++ {
		code += strconv.Itoa(rand.New(s).Intn(10))
	}
	return
}

// TimeFormat 时间格式化
func TimeFormat(time time.Time) string {
	return time.Format("2006-01-02 15:04:05")
}
