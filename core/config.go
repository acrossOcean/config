package core

import (
	"math"
	"os"
	"path/filepath"
	"sync"

	"github.com/acrossOcean/log"
)

const (
	dataParam = "param"
	dataFile  = "file"
)

var (
	logger *log.Logger
)

type Config struct {
	// 是否监控文件变化
	isWatcher bool
	// 最大文件监听数
	maxFileNum        int8
	watcherChangeChan chan bool

	isFileParsed bool
	fileLocker   sync.RWMutex

	// 是否已经分析过参数了
	isParamParsed bool
	paramLocker   sync.Mutex

	// 配置获取顺序
	filePaths []string
	readOrder []ReadOrder

	data       map[string]map[string]interface{}
	dataLocker sync.RWMutex

	// 分割器
	splitStr string
}

func init() {
	logger = log.DefaultLogger()
	logger.SetStaticTags(map[string]interface{}{
		"_appName":    PkgName,
		"_appVersion": PkgVersion,
	})
}

func defaultConfig() *Config {
	dataMap := make(map[string]map[string]interface{})
	dataMap[dataFile] = make(map[string]interface{})
	dataMap[dataParam] = make(map[string]interface{})
	var result = &Config{
		isWatcher:         false,
		maxFileNum:        100,
		watcherChangeChan: make(chan bool, 1),
		readOrder:         DefaultReadOrder(),
		data:              dataMap,
		splitStr:          ">>",
	}

	result.wait()

	return result
}

func NewConfig() *Config {
	return defaultConfig()
}

func NewConfigWithPath(path string, morePath ...string) *Config {
	cfg := defaultConfig()
	cfg.addPath(path, morePath...)

	return cfg
}

func (receiver *Config) AddPath(path string, morePath ...string) {
	receiver.addPath(path, morePath...)
}

func (receiver *Config) SetReadOrder(first ReadOrder, more ...ReadOrder) {
	receiver.readOrder = []ReadOrder{first}

	receiver.readOrder = append(receiver.readOrder, more...)
}

func (receiver *Config) WatchChange(isOpen bool) {
	receiver.watcherChangeChan <- isOpen
}

func (receiver *Config) SetDefaultSplitStr(str string) {
	receiver.splitStr = str
}

func (receiver *Config) SetMaxFileWatch(num int8) {
	if num < 0 {
		num = math.MaxInt8
	}

	receiver.maxFileNum = num
}

func (receiver *Config) searchVal(key string) (interface{}, bool) {
	for _, searchType := range receiver.readOrder {
		switch searchType {
		case ReadFromEnv:
			result, ok := receiver.readFromEnv(key)
			if ok {
				return result, true
			}
		case ReadFromParam:
			result, ok := receiver.readFromParam(key)
			if ok {
				return result, true
			}
		case ReadFromFile:
			result, ok := receiver.readFromFile(key)
			if ok {
				return result, true
			}
		default:
			logger.Error("未知的搜索方式:%v", searchType)
		}
	}

	return nil, false
}

func (receiver *Config) addPath(path string, morePath ...string) {
	receiver.filePaths = append(receiver.filePaths, path)
	for _, path := range morePath {
		receiver.filePaths = append(receiver.filePaths, receiver.getFilePath(path)...)
	}

	receiver.parseFile()
}

func (receiver *Config) getFilePath(path string) []string {
	result := make([]string, 0)
	workPath, err := os.Getwd()
	if err != nil {
		logger.Error("获取文件当前目录失败1:", err.Error())
	}

	appPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logger.Error("获取文件当前目录失败2:", err.Error())
	}

	workPath = filepath.Join(workPath, path)
	appPath = filepath.Join(appPath, path)

	finalPath := workPath
	if !FileExists(workPath) {
		if FileExists(appPath) {
			finalPath = appPath
		} else {
			finalPath = ""
		}
	}

	if finalPath == "" {
		return result
	}

	f, _ := os.Stat(finalPath)
	if f.IsDir() && f.Name() != filepath.Base(finalPath) {
		result = append(result, receiver.getDirFile(finalPath)...)
	} else {
		result = append(result, finalPath)
	}

	return result
}

func (receiver *Config) getDirFile(path string) []string {
	f, _ := os.Stat(path)
	if f.IsDir() {
		result := make([]string, 0)

		_ = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if len(result) >= 100 {
				return nil
			}

			if err != nil {
				logger.Error("读取文件", path, "错误:", err.Error())
			}

			if !info.IsDir() {
				result = append(result, path)
			} else {
				result = append(result, receiver.getDirFile(path)...)
			}

			return nil
		})

		return result
	} else {
		return []string{path}
	}
}
