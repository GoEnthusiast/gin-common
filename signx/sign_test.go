package signx

import (
	"encoding/json"
	"fmt"
	"github.com/GoEnthusiast/gin-common/sortx"
	"testing"
)

const (
	TESTPARAM1 = `{"orgId":2676410,"query":{"pagination":{"offset":0,"limit":20},"fields":null,"orderBy":[{"field":"id","sortOrder":"ASCENDING"}],"conditions":[]}}`
	TESTPARAM2 = "{\"orgId\":2676410,\"query\":{\"pagination\":{\"offset\":0,\"limit\":20},\"fields\":null,\"orderBy\":[{\"field\":\"id\",\"sortOrder\":\"ASCENDING\"}],\"conditions\":[]}}"
	TESTPARAM3 = "orgId2676410queryconditionsfieldsnullorderByfieldidsortOrderASCENDINGpaginationlimit20offset0"
)

func TestSign_Encry(t *testing.T) {
	// 加密
	encodeEn := NewSignx(
		WithPath("/api/v1/spider/asa/object/find-campaigns"),
		WithSalt("yyy-mobile:spider-server-go"),
	)
	reqParamStr := sortx.SortString(TESTPARAM3, sortx.ASC)
	signStr, err := encodeEn.Encry(reqParamStr)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return
	}
	fmt.Printf("sign: %s\n", signStr)
}

func TestSign_Decry(t *testing.T) {
	signStr := "dzcxHikxO1Q9IStEK0JrJD5JKz9/XQMoPQdQaiZCLww8QH8YPzUkDihqNQMjLi9BNzg0BDYIYB4qWwAVE3ofFRMhCQJMMAMKGz5VHg0lWBkHV0YFCwlcBE9BXAQUV0taBDczABpDJxYBByYfWUA+EwAtIB0XLSBHASpBUi9aVhhdBEAUWEMWShoUDBZKE14SSh0UDxdOE0AfEBdJQAwDBBwEUxQeGiRGQB9DVkNHUEEeU1g5Wko="
	decodeEn := NewSignx(
		WithPath("/api/v1/spider/asa/object/find-campaigns"), // 要与加密设置一致，否则会解密失败
		WithSalt("yyy-mobile:spider-server-go"),              // 要与加密设置一致，否则会解密失败
		WithTimeout(86400),                                   // 根据实际需求设置sign过期时间，单位秒
	)
	patams, err := decodeEn.Decry(signStr)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return
	}
	if patams != sortx.SortString(TESTPARAM3, sortx.ASC) {
		fmt.Printf("error: %s\n", "解析参数与测试参数不相等")
		return
	}
	fmt.Printf("params: %s\n", patams)
}

// TestSignEnScript 给爬虫用
func TestSignEnScript(t *testing.T) {
	var param map[string]any
	err := json.Unmarshal([]byte(TESTPARAM2), &param)
	if err != nil {
		panic(err)
	}
	signStr, err := SignEnScript("", "yyy-mobile:spider-server-go", param)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("sign: %v\n", signStr)
}
