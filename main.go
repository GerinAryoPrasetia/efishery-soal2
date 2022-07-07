package main

import (
	"net/http"
)

func main() {
	// router := http.NewServeMux()
	pegawai := newPegawaiPortal()
	customerHandlers := newCustomerHandlers()
	loanDocumentHandlers := newLoanDocumentHandlres()
	http.HandleFunc("/customers", customerHandlers.customers)
	http.HandleFunc("/customers/", customerHandlers.getCustomer)
	http.HandleFunc("/loan-documents", loanDocumentHandlers.loanDocuments)
	http.HandleFunc("/pegawai", pegawai.handler)

	err := http.ListenAndServe(":4000", nil)
	if err != nil {
		panic(err)
	}
}
