package sortx

import (
	"math/rand"
	"sort"
)

const (
	ASC  = "adc"  // 升序
	DESC = "desc" // 降序
	RAND = "rand" // 随机
)

// SortString
/*
	对字符串进行排序
	sortingType: 排序类型;
		可能的值-升序: asc
		可能的值-降序: desc
		可能的值-随机: rand
	return string
*/
func SortString(s string, sortingType string) string {
	switch sortingType {
	case ASC:
		paramChars := []rune(s)
		sort.Slice(paramChars, func(i, j int) bool {
			return paramChars[i] < paramChars[j]
		})

		return string(paramChars)
	case DESC:
		paramChars := []rune(s)
		sort.Slice(paramChars, func(i, j int) bool {
			return paramChars[i] > paramChars[j]
		})

		return string(paramChars)
	case RAND:
		paramChars := []rune(s)
		sort.Slice(paramChars, func(i, j int) bool {
			return rand.Intn(2) == 0
		})

		return string(paramChars)
	default:

		return s
	}
}
