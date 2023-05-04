package signx

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/GoEnthusiast/gin-common/golangx/slicex"
	"github.com/GoEnthusiast/gin-common/sortx"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type (
	Sign interface {
		Encry(string) (string, error)
		Decry(string) (string, error)
	}
	sign struct {
		timeout int64
		path    string
		salt    string
	}
	Option func(encry *sign)
)

func (s sign) Encry(arg string) (string, error) {
	// 只保留参数中的字母、数字、下划线
	reg := regexp.MustCompile("[^a-zA-Z0-9_]")
	newArg := reg.ReplaceAllString(arg, "")
	if newArg != arg {
		return "", errors.New("Arg must contain only 'letters',' numbers', and '_'")
	}

	// 进行b64编码,组装数据
	argB64 := base64.StdEncoding.EncodeToString([]byte(newArg))
	nowDate := time.Now().UnixNano()/1e6 - 1661224081041
	aStr := fmt.Sprintf("%v@#%v@#%v@#3", argB64, s.path, nowDate)
	eSplic := strings.Split(aStr, "")
	// 生成加密数组
	a := len(eSplic)
	n := len(s.salt)
	result := []string{}
	for r := 0; r < a; r++ {
		result = append(result, string([]rune(eSplic[r])[0]^[]rune(string(s.salt[(r+10)%n]))[0]))
	}
	// 生成加密后的字符串
	resultJoin := strings.Join(result, "")
	// 再次进行b64编码
	return base64.StdEncoding.EncodeToString([]byte(resultJoin)), nil
}

func (s sign) Decry(arg string) (string, error) {
	// b64解码
	decode, _ := base64.StdEncoding.DecodeString(arg)
	// 转换成切片
	decodeSlice := strings.Split(string(decode), "")
	// 解密数据
	n := len(s.salt)
	a := len(decodeSlice)
	decrypted := make([]string, a)
	for r := 0; r < a; r++ {
		decrypted[r] = string([]rune(decodeSlice[r])[0] ^ []rune(string(s.salt[(r+10)%n]))[0])
	}
	aStr := strings.Join(decrypted, "")
	eSplic := strings.Split(aStr, "@#")
	// 提取解析出的数据
	var aStrParam, aStrPath, aStrTimestamp string
	if val, ok := slicex.GetSliceValue(eSplic, 0); ok {
		aStrParam = val.(string)
	}
	if val, ok := slicex.GetSliceValue(eSplic, 1); ok {
		aStrPath = val.(string)
	}
	if val, ok := slicex.GetSliceValue(eSplic, 2); ok {
		aStrTimestamp = val.(string)
	}
	// 判断提取数据与加密数据是否吻合
	if aStrPath != s.path {
		return "", errors.New("sign错误")
	}
	dateInt, err := strconv.ParseInt(aStrTimestamp, 10, 64)
	if err != nil {
		return "", err
	}
	if s.timeout > 0 {
		if time.Now().UnixNano()/1e6-(dateInt+1661224081041) > s.timeout*1000 {
			return "", errors.New("sign过期")
		}
	}
	// 解析参数
	param, err := base64.StdEncoding.DecodeString(aStrParam)
	if err != nil {
		return "", err
	}
	return string(param), nil
}

func NewSignx(opts ...Option) Sign {
	var s sign
	for _, opt := range opts {
		opt(&s)
	}
	if s.salt == "" {
		s.salt = "xyz517cda96abcd"
	}
	return &s
}

func WithTimeout(arg int64) Option {
	return func(encry *sign) {
		encry.timeout = arg
	}
}

func WithPath(arg string) Option {
	return func(encry *sign) {
		encry.path = arg
	}
}

func WithSalt(arg string) Option {
	return func(encry *sign) {
		encry.salt = arg
	}
}

func SignEnScript(path, salt string, param map[string]any) (string, error) {
	// 参数json化
	paramJsonByts, err := json.Marshal(param)
	if err != nil {
		return "", err
	}
	// 参数字符去除
	reg := regexp.MustCompile("[^a-zA-Z0-9_]")
	newParam := reg.ReplaceAllString(string(paramJsonByts), "")
	// 参数排序
	paramSort := sortx.SortString(newParam, sortx.ASC)
	s := NewSignx(WithPath(path), WithSalt(salt))
	// 生成加密
	return s.Encry(paramSort)
}
