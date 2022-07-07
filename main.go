package main

import "net/http"

func main() {
	// router := http.NewServeMux()
	customerHandlers := newLCustomerHandlers()
	http.HandleFunc("/customers", customerHandlers.customers)

	err := http.ListenAndServe(":4000", nil)
	if err != nil {
		panic(err)
	}
}
