package config

import "flag"

func (receiver *Config) readFromParam(key string) (interface{}, bool) {
	if !receiver.isParamParsed {
		receiver.parseParam()
	}

	receiver.dataLocker.RLock()
	param := receiver.data[dataParam]
	receiver.dataLocker.RUnlock()

	result, ok := param[key]
	return result, ok
}

func (receiver *Config) parseParam() {
	receiver.paramLocker.Lock()
	defer receiver.paramLocker.Unlock()

	if flag.Parsed() {
		flag.Parse()
	}

	dataMap := make(map[string]interface{})
	flag.Visit(func(f *flag.Flag) {
		dataMap[f.Name] = getRealV(f.Value.String())
	})

	receiver.dataLocker.Lock()
	receiver.data[dataParam] = dataMap
	receiver.dataLocker.Unlock()

	receiver.isParamParsed = true
}
