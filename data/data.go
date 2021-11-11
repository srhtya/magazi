package data

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type Data struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ReturnValue struct {
	Value string `json:"value"`
}

var dataMap map[string]ReturnValue

func GetValue(key string) (ReturnValue, error) {
	if dataMap[key].Value == "" {
		return dataMap[key], errors.New("given key is not present at this moment")
	}
	return dataMap[key], nil
}

func AddValue(d *Data) {
	dataMap[d.Key] = ReturnValue{d.Value}
}

func (r *ReturnValue) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(r)
}

func (d *Data) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(d)
}

func GetCurrentDataMap() map[string]ReturnValue {
	return dataMap
}

func init() {
	dataMap = make(map[string]ReturnValue)
	dataMap["serhat"] = ReturnValue{"value1"}
	dataMap["omer"] = ReturnValue{"value2"}
	dataMap["yavuz"] = ReturnValue{"value3"}
	dataMap["yigit"] = ReturnValue{"value4"}
}

func UpdateBackup() {
	dataMap := GetCurrentDataMap()
	jsonMapStr, err := json.Marshal(dataMap)
	if err == nil {
		ioutil.WriteFile("/tmp/backup-data.json", jsonMapStr, os.ModePerm)
	}
	getMapFromBackup()
}

func getMapFromBackup() {
	log.Println("getMapFromBackup")
	backupFile, err := os.OpenFile("/tmp/backup-data.json", os.O_CREATE, os.ModePerm)
	log.Println(err.Error())
	if err != nil {
		encoder := json.NewEncoder(backupFile)
		var backupData []ReturnValue
		encoder.Encode(backupData)
		// dataMap := make(map[string]data.ReturnValue)
		for _, d := range backupData {
			log.Println(d)
		}
	}
	defer backupFile.Close()
}
