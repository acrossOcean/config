package core

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func (receiver *Config) parseIniFile(path string) map[string]interface{} {
	result := make(map[string]interface{})
	data, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Error("读取ini文件:", path, "错误:", err.Error())
		return result
	}

	result = receiver.parseIniData(filepath.Dir(path), data)
	result = flatMap(result)

	return result
}

func (receiver *Config) parseIniData(path string, data []byte) map[string]interface{} {
	result := make(map[string]interface{})

	bNumComment := []byte{'#'}  // number signal
	bSemComment := []byte{';'}  // semicolon signal
	bEqual := []byte{'='}       // equal signal
	bDQuote := []byte{'"'}      // quote signal
	sectionStart := []byte{'['} // section start signal
	sectionEnd := []byte{']'}   // section end signal

	// 按行读
	buf := bufio.NewReader(bytes.NewReader(data))
	// check the BOM
	head, err := buf.Peek(3)
	if err == nil && head[0] == 239 && head[1] == 187 && head[2] == 191 {
		for i := 1; i <= 3; i++ {
			_, _ = buf.ReadByte()
		}
	}

	currentSection := ""
	sectionMap := make(map[string]interface{})
	for {
		line, _, err := buf.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			logger.Error("获取ini文件:", path, " 错误:", err.Error())
			continue
		}

		// 去除前后无用空格
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			// 空行
			continue
		}

		if bytes.HasPrefix(line, bNumComment) || bytes.HasPrefix(line, bSemComment) {
			// 本行是注释
			continue
		}

		if bytes.HasPrefix(line, sectionStart) && bytes.HasSuffix(line, sectionEnd) {
			// 保留老段落的信息
			result[currentSection] = sectionMap
			// 新段落开始
			currentSection = string(line[1 : len(line)-1])
			sectionMap = map[string]interface{}{}
			continue
		}

		kvPair := bytes.SplitN(line, bEqual, 2)
		key := string(string(bytes.TrimSpace(kvPair[0])))

		// 如果是 include xxx.ini 这种, 再解析一下 xxx.ini
		if len(kvPair) == 1 && strings.HasPrefix(key, "include") {
			files := strings.Fields(key)
			if len(files) == 2 && files[0] == "include" {
				filePath := strings.Trim(files[1], "\"")
				if !filepath.IsAbs(filePath) {
					filePath = filepath.Join(path, filePath)
				}

				v := receiver.parseIniFile(filePath)

				for kk, vv := range v {
					result[kk] = vv
				}
			}
		}

		if len(kvPair) != 2 {
			logger.Error("ini配置文件:", path, "错误, 没有 k=v 的格式:", string(line))
			continue
		}

		value := bytes.TrimSpace(kvPair[1])
		if bytes.HasPrefix(value, bDQuote) {
			result[key] = string(bytes.Trim(value, string(bDQuote)))
			continue
		}

		// 判断v 的真实类型
		result[key] = getRealV(string(value))
	}

	return result
}
