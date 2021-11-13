package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// data struct for datastore
type Data struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var dataMap map[string]Data

// get value from datastore for a given key
func GetData(key string) (Data, error) {
	return getDataFromMap(key, dataMap)
}

func getDataFromMap(key string, data map[string]Data) (Data, error) {
	if data[key].Value == "" {
		return data[key], errors.New("given key is not present at this moment")
	}
	return data[key], nil
}

// add new value to datastore for a given key
func AddData(d *Data) {
	if dataMap == nil {
		dataMap = make(map[string]Data)
	}
	addDataToMap(d, dataMap)
}

func addDataToMap(d *Data, data map[string]Data) {
	data[d.Key] = Data{d.Key, d.Value}
}

// convert data to json
func (d *Data) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(d)
}

// convert json to data
func (d *Data) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(d)
}

// update backup file current datastore
func UpdateFile(file string) {
	updateFileWithMap(file, dataMap)
}

func updateFileWithMap(file string, data map[string]Data) {
	loc := fmt.Sprintf("%s%s%s", "/tmp/", file, ".json")
	jsonMapStr, err := json.Marshal(data)
	if err == nil {
		ioutil.WriteFile(loc, jsonMapStr, os.ModePerm)
	}
}

// reload data from backup file if it exists
func PrepareDataStore(file string) {
	loc := fmt.Sprintf("%s%s%s", "/tmp/", file, ".json")
	if _, err := os.Stat(loc); errors.Is(err, os.ErrNotExist) {
		dataMap = make(map[string]Data)
	} else {
		backupFile, err := ioutil.ReadFile(loc)
		if err == nil {
			json.Unmarshal([]byte(backupFile), &dataMap)
		} else {
			dataMap = make(map[string]Data)
		}
	}
}
