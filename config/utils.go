package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// FileExists reports whether the named file or directory exists.
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// 如果是结构体, 那么抻平结构体
/*
{
	"a": {
		"b":2
	}
}
变为:
{
	"a": {
		"b":2
	},
	"a>>b":2
}
*/
func flatMap(m map[string]interface{}) map[string]interface{} {
	result := m
	for k, v := range m {
		switch v := v.(type) {
		// 不会有int出现, 只会出现最大的可能性 float64
		case float64:
			s := fmt.Sprint(v)
			if strings.Index(s, ".") > 0 {
				result[k] = v
			} else {
				result[k] = int(v)
			}
		case string:
			result[k] = v
		case map[string]interface{}:
			mm := flatMap(v)
			for kk, vv := range mm {
				m[k+">>"+kk] = vv
			}
		}
	}

	return result
}

// 获取传入str 的真实身份,  "1" -> int(1)  "true" -> bool(true) ....
func getRealV(str string) interface{} {
	// 如果是 int
	i, err := strconv.Atoi(str)
	if err == nil {
		return i
	}
	// 如果是 float
	f, err := strconv.ParseFloat(str, 64)
	if err == nil {
		return f
	}

	b, err := strconv.ParseBool(str)
	if err == nil {
		return b
	}

	return str
}
