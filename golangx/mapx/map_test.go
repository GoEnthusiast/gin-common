package mapx

import (
	"encoding/json"
	"fmt"
	"testing"
)

const (
	TESTPARAM1 = `{"orgId":2676410,"query":{"pagination":{"offset":0,"limit":20},"fields":null,"orderBy":[{"field":"id","sortOrder":"ASCENDING"}],"conditions":[]}}`
	TESTPARAM2 = "{\"orgId\":2676410,\"query\":{\"pagination\":{\"offset\":0,\"limit\":20},\"fields\":null,\"orderBy\":[{\"field\":\"id\",\"sortOrder\":\"ASCENDING\"}],\"conditions\":[]}}"
)

func TestMapToSignStr(t *testing.T) {
	var m map[string]any
	err := json.Unmarshal([]byte(TESTPARAM1), &m)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
	}
	result := MapToString(m, "[^a-zA-Z0-9_]")
	fmt.Printf("result: %s\n", result)
}
