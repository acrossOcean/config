package config

import (
	"github.com/fsnotify/fsnotify"
)

func (receiver *Config) wait() {
	go func() {
		stopChan := make(chan struct{}, 1)
		for {
			switch <-receiver.watcherChangeChan {
			case true:
				if !receiver.isWatcher {
					receiver.isWatcher = true
					go receiver.startWatch(stopChan)
				}
			case false:
				if receiver.isWatcher {
					receiver.isWatcher = false
					stopChan <- struct{}{}
				}
			}
		}
	}()
}

func (receiver *Config) startWatch(stopChan chan struct{}) {
	// 监听文件变化,并更新
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Error("开启文件监听失败:", err.Error())
		<-stopChan
		return
	}
	defer watcher.Close()

	for _, path := range receiver.filePaths {
		for _, p := range receiver.getFilePath(path) {
			_ = watcher.Add(p)
		}
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op == fsnotify.Remove || event.Op == fsnotify.Write || event.Op == fsnotify.Create {
				receiver.parseFile()
			}
		case err := <-watcher.Errors:
			logger.Error("文件监听错误:", err)
		case <-stopChan:
			return
		}
	}
}
