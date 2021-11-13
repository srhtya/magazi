package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func Test_NewMagazi(t *testing.T) {
	l := log.New(os.Stdout, "magazi-test-api ", log.LstdFlags)
	file := "test-file"

	m := NewMagazi(l, file)
	if l != m.l {
		t.Errorf("error in NewMagazi: logs should be same")
	}

	if file != m.file {
		t.Errorf("error in NewMagazi: files %s and %s should be same", file, m.file)
	}
}

func Test_False_GetData(t *testing.T) {
	l := log.New(os.Stdout, "magazi-test-api ", log.LstdFlags)
	file := "test-file"
	m := NewMagazi(l, file)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	m.ServeHTTP(w, req)
	res := w.Result()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("error in getData: expected error to be nil got %v", err)
	}
	expected := "Key not provided in query parameters"
	if !strings.Contains(string(data), expected) {
		t.Errorf("error in getData: expected to be %s got %s", expected, string(data))
	}
}

func Test_False_AddData(t *testing.T) {
	l := log.New(os.Stdout, "magazi-test-api ", log.LstdFlags)
	file := "test-file"
	m := NewMagazi(l, file)

	bodyReader := strings.NewReader(`{"Key": "test-key", Value: "test-value"}`)
	req := httptest.NewRequest(http.MethodPost, "/", bodyReader)
	w := httptest.NewRecorder()

	m.ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("error in getData: expected error to be nil got %v", err)
	}
	if !strings.Contains(string(data), "Unable to unmarshal json") {
		t.Errorf("error in getData: expected response to be \"Unable to unmarshal json\" got %s", string(data))
	}
}

func Test_AddData_GetData(t *testing.T) {
	l := log.New(os.Stdout, "magazi-test-api ", log.LstdFlags)
	file := "test-file"
	m := NewMagazi(l, file)

	bodyReader := strings.NewReader(`{"Key": "test-key", "Value": "test-value"}`)
	req := httptest.NewRequest(http.MethodPost, "/", bodyReader)
	w := httptest.NewRecorder()
	m.ServeHTTP(w, req)

	req = httptest.NewRequest(http.MethodGet, "/?key=test-key", nil)
	m.ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("error in getData: expected error to be nil got %v", err)
	}
	expected := `{"key":"test-key","value":"test-value"}`
	if !strings.Contains(string(data), expected) {
		t.Errorf("error in getData: expected response to be %s got %s", expected, string(data))
	}
}

func Test_Flush(t *testing.T) {
	l := log.New(os.Stdout, "magazi-test-api ", log.LstdFlags)
	file := "test-file"
	m := NewMagazi(l, file)

	bodyReader := strings.NewReader(`{"Key": "test-key", "Value": "test-value"}`)
	req := httptest.NewRequest(http.MethodPost, "/", bodyReader)
	w := httptest.NewRecorder()
	m.ServeHTTP(w, req)

	req = httptest.NewRequest(http.MethodPost, "/flush", nil)
	m.ServeHTTP(w, req)

	loc := fmt.Sprintf("%s%s%s", "/tmp/", file, ".json")
	backupFile, err := ioutil.ReadFile(loc)
	if err != nil {
		t.Errorf("error in flushData: expected error to be nil got %v", err)
	}
	expected := `{"test-key":{"key":"test-key","value":"test-value"}}`
	if !strings.Contains(string(backupFile), expected) {
		t.Errorf("error in getData: expected response to be %s got %s", expected, string(backupFile))
	}
}
