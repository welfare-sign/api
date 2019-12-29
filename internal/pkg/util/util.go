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

	"github.com/chenhg5/collection"
	"github.com/gofrs/uuid"
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

//NewV4 添加重试机制
func NewV4() (uuid.UUID, error) {
	c := 0
look:
	uid, err := uuid.NewV4()
	if err == nil {
		return uid, nil
	}
	if c < 3 {
		c++
		goto look
	}
	return [16]byte{}, errors.WithMessage(err, "NewV4 retry 3 time error")
}

// GetRecommandNum 获取推荐的数字
func GetRecommandNum(oldNum uint64, nums []uint64) ([]uint64, error) {
	if len(nums) == 0 {
		return []uint64{}, errors.New("不存在的数字列表")
	}
	maxNum := getPropeNum(oldNum+1, "+", nums)
	minNum := getPropeNum(oldNum-1, "-", nums)
	return []uint64{minNum, maxNum}, errors.New("您填入的数字已被占用")
}

func getPropeNum(num uint64, tag string, nums []uint64) uint64 {
	if !collection.Collect(nums).Contains(num) || num == 0 {
		return num
	}
	if tag == "+" {
		return getPropeNum(num+1, tag, nums)
	} else {
		return getPropeNum(num-1, tag, nums)
	}

}

// SortRecord .
type SortRecord struct {
	Tag string
	ID  uint64
	Num uint64
}

// SortRecordList .
type SortRecordList []SortRecord

// Len .
func (s SortRecordList) Len() int {
	return len(s)
}

// Less .
func (s SortRecordList) Less(i, j int) bool {
	if s[i].Num == s[j].Num {
		return s[i].Tag == "+"
	}
	return s[i].Num < s[j].Num
}

// Swap .
func (s SortRecordList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Abs .
func Abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}

// GetDefineFriday 获取周五12点
func GetDefineFriday() time.Time {
	now := time.Now()

	offset := int(time.Friday - now.Weekday())
	if offset == 5 {
		offset = -2
	}

	return time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, time.Local).AddDate(0, 0, offset)
}

// GetDefineSaturday 获取周六零点
func GetDefineSaturday() time.Time {
	now := time.Now()

	offset := int(time.Saturday - now.Weekday())
	if offset == 6 {
		offset = -1
	}
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
}

// IsCurrentTimeBetweenD1AndD2 计算当前时间是否在两个时间段内
func IsCurrentTimeBetweenD1AndD2(d1, d2 time.Time) bool {
	now := time.Now()
	return now.After(d1) && now.Before(d2)
}
