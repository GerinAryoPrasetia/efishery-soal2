package main

import (
	"context"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type key int

const (
	dbContext key = iota
)

func DBMiddleware(next http.Handler, db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), dbContext, db)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func main() {
	// router := http.NewServeMux()
	dsn := "host=localhost user=postgres password=whoami00 dbname=loan_application port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(Customer{}, Pegawai{}, LoanDocument{})
	if err != nil {
		panic(err)
	}

	router := http.NewServeMux()

	pegawai := newPegawaiPortal()
	customerHandlers := newCustomerHandlers()
	loanDocumentHandlers := newLoanDocumentHandlres()
	//route dengan handlers 100% vanilla go, tanpa gorm atau DB (penyimpanan data via memory)
	http.HandleFunc("/customers", customerHandlers.customers)
	http.HandleFunc("/customers/", customerHandlers.getCustomer)
	http.HandleFunc("/loan-documents", loanDocumentHandlers.loanDocuments)
	http.HandleFunc("/loan-documents/", loanDocumentHandlers.post)
	http.HandleFunc("/pegawai", pegawai.handler)

	//route dengan GORM dan DB
	router.HandleFunc("/loan-documents", getLoanWithDB)
	router.HandleFunc("/loan-document", getLoanByID)
	router.HandleFunc("/loan-document/create", postLoanWithDB)
	router.HandleFunc("/loan-document/delete", deleteLoanWithDB)
	router.HandleFunc("/loan-document/update", updateLoanDocument)
	router.HandleFunc("/loan-document/change-status", updateStatusLoan)

	router.HandleFunc("/customers", getAllCustomers)
	router.HandleFunc("/customer/create", createCustomer)

	routerMiddleware := DBMiddleware(router, db)

	// Inisiasi server tanpa database
	// err = http.ListenAndServe(":4000", nil)
	// if err != nil {
	// 	panic(err)
	// }

	err = http.ListenAndServe(":4000", routerMiddleware)
	if err != nil {
		panic(err)
	}
}
