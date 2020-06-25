package core

import "os"

func (receiver *Config) readFromEnv(key string) (interface{}, bool) {
	str, ok := os.LookupEnv(key)
	return getRealV(str), ok
}
