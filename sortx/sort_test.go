package sortx

import (
	"fmt"
	"testing"
)

const TestText1 = "orgId2676410queryconditionsfieldsnullorderByfieldidsortOrderASCENDINGpaginationlimit20offset0"

func TestSortString(t *testing.T) {
	ts := SortString(TestText1, ASC)
	fmt.Printf("result: %s\n", ts)
}
