package config

import (
	"github.com/acrossOcean/config/config"
)

var defConf *config.Config

func init() {
	defConf = config.NewConfig()
}

func AddPath(path string, morePath ...string) {
	defConf.AddPath(path, morePath...)
}

func SetReadOrder(readOrder config.ReadOrder, moreOrder ...config.ReadOrder) {
	defConf.SetReadOrder(readOrder, moreOrder...)
}

func WatchChange(isOpen bool) {
	defConf.WatchChange(isOpen)
}

func String(key string) (string, bool) {
	return defConf.String(key)
}

func DefaultString(key string, defaultVal string) string {
	return defConf.DefaultString(key, defaultVal)
}

func StringList(key string) ([]string, error) {
	return defConf.StringList(key)
}

func Int(key string) (int, bool) {
	return defConf.Int(key)
}

func DefaultInt(key string, defaultVal int) int {
	return defConf.DefaultInt(key, defaultVal)
}

func Bool(key string) (bool, bool) {
	return defConf.Bool(key)
}

func DefaultBool(key string, defaultVal bool) bool {
	return defConf.DefaultBool(key, defaultVal)
}

func Int64(key string) (int64, bool) {
	return defConf.Int64(key)
}

func DefaultInt64(key string, defaultVal int64) int64 {
	return defConf.DefaultInt64(key, defaultVal)
}

func Get(key string) (interface{}, bool) {
	return defConf.Get(key)
}

func GetCurrentCache() map[string]map[string]interface{} {
	return defConf.GetCurrentCache()
}
