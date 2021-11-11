package handlers

import (
	"log"
	"net/http"

	"github.com/serhatyavuzyigit/magazi/data"
)

type Magazi struct {
	l *log.Logger
}

func NewMagazi(l *log.Logger) *Magazi {
	return &Magazi{l}
}

func (m *Magazi) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		m.getValue(rw, r)
		return
	}

	if r.Method == http.MethodPost {
		m.addValue(rw, r)
		return
	}

	// catch all
	// if no method is satisfied return an error
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (m *Magazi) getValue(rw http.ResponseWriter, r *http.Request) {
	m.l.Println("Handle GET")
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(rw, "Key not provided in query parameters", http.StatusBadRequest)
	} else {
		// fetch the products from the datastore
		v, err := data.GetValue(key)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
		} else {
			// serialize the list to JSON
			err = v.ToJSON(rw)
			if err != nil {
				http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
			}
		}

	}
}

func (m *Magazi) addValue(rw http.ResponseWriter, r *http.Request) {
	m.l.Println("Handle POST")

	d := &data.Data{}
	err := d.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	data.AddValue(d)
}
