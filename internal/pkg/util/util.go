package util

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
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
	for i := 0; i < viper.GetInt(config.KeySMSLength); i++ {
		code += strconv.Itoa(rand.New(s).Intn(10))
	}
	return
}

// TimeFormat 时间格式化
func TimeFormat(time time.Time) string {
	return time.Format("2006-01-02 15:04:05")
}

// GetWXSignString 获取微信JS API 签名
func GetWXSignString(values map[string]string) string {
	var (
		i    = 0
		size = len(values)
	)

	//第一步：把Key按字典的字母顺序排序
	sortedParams := make([]string, size)
	for key := range values {
		sortedParams[i] = key
		i++
	}
	sort.Strings(sortedParams)

	//第二步：把所有参数名和参数值串在一起
	buff := bytes.Buffer{}
	for _, key := range sortedParams {
		val, ok := values[key]
		if !ok || len(val) == 0 {
			continue
		}

		buff.WriteString(key)
		buff.WriteString("=")
		buff.WriteString(val)
		buff.WriteString("&")
	}
	buff.Truncate(buff.Len() - 1)

	h := sha1.New()
	h.Write(buff.Bytes())
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

// GetIP .
func GetIP() (ip string) {
	addr, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}

	for _, address := range addr {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				return
			}
		}
	}
	return
}

// Post .
func Post(xml string, url string) (string, error) {
	resp, err := http.Post(url, "text/xml", strings.NewReader(xml))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// GenerateNonceStr .
func GenerateNonceStr(length int) string {
	stem := "HsZPdYr0KeXxI1WcJL2CtDpNfBa3AuMob4v5TgUqOnVk6Ghj7iR8FzwQ9lyEmS"
	buf := make([]byte, length)

	stemSize := len(stem)
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < length; i++ {
		buf[i] = byte(stem[rand.Int()%stemSize])
	}
	return string(buf)
}
