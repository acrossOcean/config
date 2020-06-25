package core

import (
	"os"
	"strings"
)

func (receiver *Config) readFromFile(key string) (interface{}, bool) {
	if !receiver.isFileParsed {
		receiver.parseFile()
	}

	receiver.dataLocker.RLock()
	data := receiver.data[dataFile]
	receiver.dataLocker.RUnlock()

	result, ok := data[key]
	return result, ok
}

func (receiver *Config) parseFile() {
	receiver.fileLocker.Lock()
	defer receiver.fileLocker.Unlock()

	// 读取文件内容, 区分不同文件
	valueMap := make(map[string]interface{})

	for _, path := range receiver.filePaths {
		for _, p := range receiver.getFilePath(path) {
			var dataMap map[string]interface{}
			switch strings.Split(p, ".")[len(strings.Split(p, "."))-1] {
			case "json":
				dataMap = receiver.parseJsonFile(p)
			case "ini":
				dataMap = receiver.parseIniFile(p)
			case "yaml":
				dataMap = receiver.parseYamlFile(p)
			default:
				f, _ := os.Stat(p)
				if !f.IsDir() {
					logger.Error("不支持的配置文件格式:%s", p)
				}
			}

			for k, v := range dataMap {
				valueMap[k] = v
			}
		}
	}

	receiver.data[dataFile] = valueMap
}
