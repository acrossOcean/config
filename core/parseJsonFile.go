package core

import (
	"encoding/json"
	"io/ioutil"
)

func (receiver *Config) parseJsonFile(path string) map[string]interface{} {
	result := make(map[string]interface{})

	info, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Error("读取json配置文件:", path, "错误:", err.Error())
		return result
	}

	err = json.Unmarshal(info, &result)
	if err != nil {
		logger.Error("解析json配置文件:", path, "错误:", err.Error())
		return result
	}

	result = flatMap(result)

	return result
}
