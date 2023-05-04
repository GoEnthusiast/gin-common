package mapx

import (
	"encoding/json"
	"regexp"
)

// MapToString
/*
	将map转换成字符串
	rsRegeox: 通过正则表达式设置允许保留的字符，例如"[^a-zA-Z0-9_]"，只保留字母、数字、下划线
	return string
*/
func MapToString(m map[string]any, rsRegeox string) string {
	if m == nil || len(m) == 0 {
		return ""
	}

	var str string
	if mByte, err := json.Marshal(m); err != nil {
		return ""
	} else {
		str = string(mByte)
	}

	reg := regexp.MustCompile(rsRegeox)

	return reg.ReplaceAllString(str, "")
}
