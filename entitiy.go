package main

type Customer struct {
	ID             string `json:"customer_id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	Age            int    `json:"age"`
	IdentityNumber string `json:"identity_number"`
}

type LoanDocument struct {
	ID              string  `json:"loan_id"`
	CustomerID      int     `json:"customer_id"`
	Status          bool    `json:"status"`
	ApplicationDate string  `json:"application_date"`
	NominalPinjaman float64 `json:"nominal_pinjaman"`
	JenisPinjaman   string  `json:"jenis_pinjaman"`
	Customer        Customer
}

type Pegawai struct {
	ID        string
	FirstName string
	LastName  string
}
