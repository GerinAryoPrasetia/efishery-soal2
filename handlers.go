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
	"gorm.io/gorm"
)

type customerHandlers struct {
	sync.Mutex
	store map[string]Customer
}

type loanDocumentHandlres struct {
	sync.Mutex
	store map[string]LoanDocument
}

// =========================================================== //
// ===== README FIRST ==== //
// Pada App ini terdapat dua case yaitu :
// 1. Proses CRUD dan terhubung dengan DB PostrgreSQL
// 2. Proses CRUD, tetapi tidak terhubung dengan DB, dan data disimpan ke dalam memory
// ===== README FIRST ==== //
// =========================================================== //

// 1. HANDLERS MENGGUNAKAN DB + GORM

// =========================================================== //

func getLoanWithDB(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value(dbContext).(*gorm.DB)

	if !ok {
		http.Error(w, "no database found", http.StatusInternalServerError)
		return
	}

	var loanDocuments []LoanDocument

	if err := db.Find(&loanDocuments).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(loanDocuments)

}

func getLoanByID(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value(dbContext).(*gorm.DB)

	if !ok {
		http.Error(w, "no database found", http.StatusInternalServerError)
		return
	}

	loan_id := r.URL.Query().Get("loan_id")

	loanDocument := LoanDocument{
		ID: loan_id,
	}

	if err := db.Find(&loanDocument).Error; err != nil {
		w.Write([]byte(fmt.Sprintf("Loan Document Not Found With ID : '%v'", loan_id)))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loanDocument)
}

func updateLoanDocument(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value(dbContext).(*gorm.DB)

	if !ok {
		http.Error(w, "no database found", http.StatusInternalServerError)
		return
	}

	loan_id := r.URL.Query().Get("loan_id")

	var loanDocument LoanDocument
	if err := json.NewDecoder(r.Body).Decode(&loanDocument); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	loanDocument.ID = loan_id
	existLoan := LoanDocument{
		ID: loan_id,
	}

	if err := db.First(&existLoan).Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	existLoan.CustomerID = loanDocument.CustomerID
	existLoan.NominalPinjaman = loanDocument.NominalPinjaman
	existLoan.JenisPinjaman = loanDocument.JenisPinjaman
	existLoan.ApplicationDate = loanDocument.ApplicationDate
	existLoan.Status = loanDocument.Status

	if err := db.Save(&existLoan).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existLoan)
}

func postLoanWithDB(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value(dbContext).(*gorm.DB)

	if !ok {
		http.Error(w, "no database found", http.StatusInternalServerError)
		return
	}

	var loanDocument LoanDocument

	if err := json.NewDecoder(r.Body).Decode(&loanDocument); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user_id := r.URL.Query().Get("user_id")

	customer := Customer{
		ID: user_id,
	}

	if err := db.First(&customer).Error; err != nil {
		w.Write([]byte(fmt.Sprintf("Customer Not Found With ID : '%v'", user_id)))
		return
	}

	uuid_document := uuid.NewV1().String()
	loanDocument.ID = uuid_document
	loanDocument.CustomerID = user_id
	loanDocument.Status = false

	if err := db.Create(&loanDocument).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(loanDocument)

}

func deleteLoanWithDB(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value(dbContext).(*gorm.DB)
	if !ok {
		http.Error(w, "no database found", http.StatusInternalServerError)
		return
	}

	// loan_id, err := strconv.ParseInt(r.URL.Query().Get("loan_id"), 10, 64)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	loan_id := r.URL.Query().Get("loan_id")
	existLoan := LoanDocument{
		ID: loan_id,
	}

	if err := db.First(&existLoan).Error; err != nil {
		w.Write([]byte(fmt.Sprintf("Loan Document Not Found With ID : '%v'", loan_id)))
		return
	}

	if err := db.Delete(&existLoan).Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existLoan)
}

func updateStatusLoan(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value(dbContext).(*gorm.DB)
	if !ok {
		http.Error(w, "no database found", http.StatusInternalServerError)
		return
	}

	loan_id := r.URL.Query().Get("loan_id")

	var loanDocument LoanDocument

	loanDocument.ID = loan_id
	existLoan := LoanDocument{
		ID: loan_id,
	}

	if err := db.First(&existLoan).Error; err != nil {
		w.Write([]byte(fmt.Sprintf("Loan Document Not Found With ID : '%v'", loan_id)))
		return
	}

	existLoan.Status = !loanDocument.Status

	if err := db.Save(&existLoan).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existLoan)
}

func getAllCustomers(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value(dbContext).(*gorm.DB)
	if !ok {
		http.Error(w, "no database found", http.StatusInternalServerError)
		return
	}

	var customers []Customer

	if err := db.Find(&customers).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

func createCustomer(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value(dbContext).(*gorm.DB)

	if !ok {
		http.Error(w, "no database found", http.StatusInternalServerError)
		return
	}

	var customer Customer

	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	uuid_customer := uuid.NewV1().String()
	customer.ID = uuid_customer

	if err := db.Create(&customer).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(customer)
}

// =========================================================== //

// 2. HANDLERS TANPA MENGGUNAKAN DB + GORM (100% Vanilla go)

// =========================================================== //

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

	var loanDocument LoanDocument
	err = json.Unmarshal(bodyBytes, &loanDocument)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fmt.Println(parts[2])

	uuid_document := uuid.NewV1().String()
	loanDocument.ID = uuid_document
	loanDocument.CustomerID = parts[2] //simpan customer id yang diambil dari parameters

	h.Lock()
	h.store[loanDocument.ID] = loanDocument
	defer h.Unlock()
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
