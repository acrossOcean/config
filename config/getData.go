package config

import "fmt"

func (receiver *Config) Get(key string) (interface{}, bool) {
	return receiver.searchVal(key)
}

func (receiver *Config) GetCurrentCache() map[string]map[string]interface{} {
	dataMap := make(map[string]map[string]interface{})

	receiver.dataLocker.RLock()
	for k, v := range receiver.data {
		data := make(map[string]interface{})
		for kk, vv := range v {
			data[kk] = vv
		}
		dataMap[k] = data
	}
	receiver.dataLocker.RUnlock()

	return dataMap
}

func (receiver *Config) String(key string) (string, bool) {
	if result, ok := receiver.searchVal(key); ok {
		switch result := result.(type) {
		case string:
			return result, true
		default:
			return fmt.Sprint(result), true
		}
	} else {
		return "", false
	}
}

func (receiver *Config) DefaultString(key string, defaultVal string) string {
	if result, ok := receiver.String(key); ok {
		return result
	} else {
		return defaultVal
	}
}

func (receiver *Config) StringList(key string) ([]string, error) {
	if result, ok := receiver.searchVal(key); ok {
		switch result := result.(type) {
		case []interface{}:
			return receiver.convertToStringList(result)
		case string:
			return []string{result}, nil
		default:
			return []string{fmt.Sprint(result)}, nil
		}
	}
	return nil, nil
}

func (receiver *Config) Int(key string) (int, bool) {
	if result, ok := receiver.searchVal(key); ok {
		switch result := result.(type) {
		case int:
			return result, true
		default:
			return 0, false
		}
	} else {
		return 0, false
	}
}

func (receiver *Config) DefaultInt(key string, defaultVal int) int {
	if result, ok := receiver.Int(key); ok {
		return result
	} else {
		return defaultVal
	}
}

func (receiver *Config) Bool(key string) (bool, bool) {
	if result, ok := receiver.searchVal(key); ok {
		switch result := result.(type) {
		case bool:
			return result, true
		default:
			return false, false
		}
	} else {
		return false, false
	}
}

func (receiver *Config) DefaultBool(key string, defaultVal bool) bool {
	if result, ok := receiver.Bool(key); ok {
		return result
	} else {
		return defaultVal
	}
}

func (receiver *Config) Int64(key string) (int64, bool) {
	if result, ok := receiver.searchVal(key); ok {
		switch result := result.(type) {
		case int64:
			return result, true
		default:
			return 0, false
		}
	} else {
		return 0, false
	}
}

func (receiver *Config) DefaultInt64(key string, defaultVal int64) int64 {
	if result, ok := receiver.Int64(key); ok {
		return result
	} else {
		return defaultVal
	}
}

func (receiver *Config) convertToStringList(values []interface{}) ([]string, error) {
	result := make([]string, 0)

	for _, value := range values {
		switch value := value.(type) {
		case string:
			result = append(result, value)
		default:
			result = append(result, fmt.Sprint(value))
		}
	}

	return result, nil
}
