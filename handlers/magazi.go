package handlers

import (
	"log"
	"net/http"
	"time"

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
	start := time.Now()
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
	requesterIP := r.RemoteAddr

	log.Printf(
		"%s\t\t%s\t\t%s\t\t%v",
		r.Method,
		r.RequestURI,
		requesterIP,
		time.Since(start),
	)
}

// add new value to datastore for a given key
func (m *Magazi) addData(rw http.ResponseWriter, r *http.Request) {
	m.l.Println("POST")
	start := time.Now()
	d := &data.Data{}
	err := d.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	} else {
		data.AddData(d)
	}
	requesterIP := r.RemoteAddr

	log.Printf(
		"%s\t\t%s\t\t%s\t\t%s\t\t%v",
		r.Method,
		r.RequestURI,
		d,
		requesterIP,
		time.Since(start),
	)
}

// flush data to file
func (m *Magazi) flushData(rw http.ResponseWriter, r *http.Request) {
	m.l.Println("POST")
	start := time.Now()
	data.UpdateFile(m.file)
	requesterIP := r.RemoteAddr

	log.Printf(
		"%s\t\t%s\t\t%s\t\t%v",
		r.Method,
		r.RequestURI,
		requesterIP,
		time.Since(start),
	)
}

// update file
func (m *Magazi) UpdateFile() {
	data.UpdateFile(m.file)
}

// prepare data store for initialization
func (m *Magazi) PrepareDataStore() {
	data.PrepareDataStore(m.file)
}
