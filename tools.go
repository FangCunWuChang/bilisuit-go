package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func GetSettingFilePath() string {
	var FilePath string
	if len(os.Args) == 1 {
		FilePath = "./setting.json"
	} else {
		FilePath = os.Args[len(os.Args)-1]
	}
	_, err := os.Lstat(FilePath)
	if err != nil {
		fmt.Printf("[%v]不存在\n", FilePath)
		os.Exit(1)
	}
	fmt.Printf("配置文件:[%v]\n", FilePath)
	return FilePath
}

type SettingContent struct {
	StartTime int64  `json:"start_time"`
	DelayTime int64  `json:"delay_time"`
	ItemId    string `json:"item_id"`
}

type SettingFile struct {
	Setting  SettingContent    `json:"setting"`
	FormData string            `json:"form_data"`
	Headers  map[string]string `json:"headers"`
}

func ReaderSetting(filePath string) (map[string]string, int64, int64, string) {
	var SettingData, _ = os.ReadFile(filePath)
	var settingContent = SettingFile{}

	_ = json.Unmarshal(SettingData, &settingContent)

	var headers = settingContent.Headers
	var formData = settingContent.FormData
	var startTime = settingContent.Setting.StartTime

	if startTime <= time.Now().Unix() {
		fmt.Printf("%v\n", "启动时间小于当前时间")
		os.Exit(2)
	}

	var delayTime = settingContent.Setting.DelayTime

	fmt.Printf("装扮id:[%v]\n", settingContent.Setting.ItemId)
	fmt.Printf("启动时间:[%v]\n", startTime)
	fmt.Printf("延时:[%vms]\n", delayTime)

	return headers, startTime, delayTime, formData
}
