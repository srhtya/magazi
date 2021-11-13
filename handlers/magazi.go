package handlers

import (
	"log"
	"net/http"

	"github.com/serhatyavuzyigit/magazi/data"
)

// magazi handler struct
type Magazi struct {
	l    *log.Logger
	file string
}

// creates new Magazi
func NewMagazi(l *log.Logger, file string) *Magazi {
	return &Magazi{l, file}
}

func (m *Magazi) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// call appropiate functions for incoming request
	if r.Method == http.MethodGet {
		m.getData(rw, r)
		return
	}

	if r.Method == http.MethodPost && r.RequestURI == "/flush" {
		m.flushData(rw, r)
		return
	}

	if r.Method == http.MethodPost {
		m.addData(rw, r)
		return
	}

	// catch all
	// if no method is satisfied return an error
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

// get value from datastore for a given key
func (m *Magazi) getData(rw http.ResponseWriter, r *http.Request) {
	m.l.Println("GET")
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(rw, "Key not provided in query parameters", http.StatusBadRequest)
	} else {
		// fetch the data from the datastore
		v, err := data.GetData(key)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
		} else {
			// serialize the object to JSON
			err = v.ToJSON(rw)
			if err != nil {
				http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
			}
		}

	}
}

// add new value to datastore for a given key
func (m *Magazi) addData(rw http.ResponseWriter, r *http.Request) {
	m.l.Println("POST")

	d := &data.Data{}
	err := d.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	} else {
		data.AddData(d)
	}
}

// flush data to file
func (m *Magazi) flushData(rw http.ResponseWriter, r *http.Request) {
	m.l.Println("POST")
	data.UpdateFile(m.file)
}

// update file
func (m *Magazi) UpdateFile() {
	data.UpdateFile(m.file)
}
