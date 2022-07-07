package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"

	uuid "github.com/satori/go.uuid"
)

type customerHandlers struct {
	sync.Mutex
	store map[string]Customer
}

type loanDocumentHandlres struct {
	sync.Mutex
	store map[string]LoanDocument
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

func (h *loanDocumentHandlres) loanDocuments(w http.ResponseWriter, r *http.Request) {
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

func (h *customerHandlers) post(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("Need content-type : application/json, but got '%s'", ct)))
		return
	}

	var customer Customer
	err = json.Unmarshal(bodyBytes, &customer)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	uuid_customer := uuid.NewV1().String()
	customer.ID = uuid_customer

	h.Lock()
	h.store[customer.ID] = customer
	defer h.Unlock()
}

func (h *customerHandlers) getCustomer(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fmt.Println(parts[2])

	h.Lock()

	customer, ok := h.store[parts[2]]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	h.Unlock()

	jsonBytes, err := json.Marshal(customer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *loanDocumentHandlres) get(w http.ResponseWriter, r *http.Request) {
	loanDocuments := make([]LoanDocument, len(h.store))

	h.Lock()
	i := 0
	for _, loandDocument := range h.store {
		loanDocuments[i] = loandDocument
		i++
	}
	h.Unlock()

	jsonBytes, err := json.Marshal(loanDocuments)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *loanDocumentHandlres) post(w http.ResponseWriter, r *http.Request) {

}

func (a pegawaiPortal) handler(w http.ResponseWriter, r *http.Request) {
	user, pass, ok := r.BasicAuth()
	if !ok || user != "pegawai" || pass != a.password {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - Unauthorized"))
		return
	}

	w.Write([]byte("Authorized"))
}

func newCustomerHandlers() *customerHandlers {
	return &customerHandlers{
		store: map[string]Customer{},
	}
}

func newLoanDocumentHandlres() *loanDocumentHandlres {
	return &loanDocumentHandlres{
		store: map[string]LoanDocument{},
	}
}

func newPegawaiPortal() *pegawaiPortal {
	password := os.Getenv("ADMIN_PASSWORD")
	if password != "" {
		panic("required env ADMIN_PASSWORD not set")
	}
	return &pegawaiPortal{password: password}
}
