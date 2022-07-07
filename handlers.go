package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
)

type customerHandlers struct {
	sync.Mutex
	store map[string]Customer
}

func (h *customerHandlers) customers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.get(w, r)
		return
	case "POST":
		h.post(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed"))
		return
	}
}

func (h *customerHandlers) post(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	var customer Customer
	err = json.Unmarshal(bodyBytes, &customer)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	h.Lock()
	h.store[customer.ID] = customer
	defer h.Unlock()
}

func (h *customerHandlers) get(w http.ResponseWriter, r *http.Request) {
	customers := make([]Customer, len(h.store))

	h.Lock()
	i := 0
	for _, customer := range h.store {
		customers[i] = customer
		i++
	}
	h.Unlock()

	jsonBytes, err := json.Marshal(customers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)

}

func newLCustomerHandlers() *customerHandlers {
	return &customerHandlers{
		store: map[string]Customer{
			"id1": Customer{
				FirstName:      "Gerin",
				LastName:       "Aryo",
				ID:             "id1",
				Email:          "gerinaryo14@gmail.com",
				Age:            21,
				IdentityNumber: "12345678",
			},
		},
	}
}
