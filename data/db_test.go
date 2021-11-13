package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
)

func Test_ReturnError_From_getDataFromMap(t *testing.T) {
	m := make(map[string]Data)
	key := "test-key"

	_, err := getDataFromMap(key, m)
	if err == nil {
		t.Errorf("error in getDataFromMap: should return an error for non existent key")
	}
}

func Test_ReturnData_From_getDataFromMap(t *testing.T) {
	m := make(map[string]Data)
	key := "test-key"
	value := "test-value"
	m[key] = Data{Key: key, Value: value}
	d, _ := getDataFromMap(key, m)

	if d.Key != key {
		t.Errorf("error in getDataFromMap: keys %s and %s should be equal", key, d.Key)
	}

	if d.Value != value {
		t.Errorf("error in getDataFromMap: values %s and %s should be equal", value, d.Value)
	}
}

func Test_ToJSON(t *testing.T) {
	d := Data{Key: "test-key", Value: "test-value"}
	buf := new(bytes.Buffer)
	err := d.ToJSON(buf)
	if err != nil {
		t.Errorf("error in ToJSON: function should not throw an error")
	}
}

func Test_FromJSON(t *testing.T) {
	r := strings.NewReader("{\"Key\":\"test-key\",\"Value\":\"test-value\"}")
	d := Data{}
	err := d.FromJSON(r)
	if err != nil {
		t.Errorf("error in FromJSON: function should not throw an error")
	}
}

func Test_updateFileWithMap(t *testing.T) {
	m := make(map[string]Data)
	key := "test-key"
	value := "test-value"
	m[key] = Data{Key: key, Value: value}
	file := "test-file"
	updateFileWithMap(file, m)

	var newMap map[string]Data
	loc := fmt.Sprintf("%s%s%s", "/tmp/", file, ".json")
	newFile, err := ioutil.ReadFile(loc)
	if err == nil {
		json.Unmarshal([]byte(newFile), &newMap)
		if !reflect.DeepEqual(m, newMap) {
			t.Errorf("error in updateFileWithMap: maps %s and %s should be equal", m, newMap)
		}
	} else {
		t.Errorf("error in updateFileWithMap: file did not created")
	}
}
