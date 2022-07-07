package main

import "time"

type Customer struct {
	ID             string `json:"customer_id" gorm:"PRIMARY_KEY"`
	FirstName      string `json:"first_name" gorm:"column:first_name"`
	LastName       string `json:"last_name" gorm:"column:last_name"`
	Email          string `json:"email" gorm:"column:email"`
	Age            int    `json:"age" gorm:"column:age"`
	IdentityNumber string `json:"identity_number" gorm:"column:identity_number"`
	// LoanDocuments  []LoanDocument `gorm:"foreign_key:customer_id"`
}

type LoanDocument struct {
	ID              string    `json:"loan_id"`
	CustomerID      string    `json:"customer_id" gorm:"column:customer_id"`
	Status          bool      `json:"status" gorm:"column:status"`
	ApplicationDate time.Time `json:"application_date" gorm:"column:application_date"`
	NominalPinjaman float64   `json:"nominal_pinjaman" gorm:"column:nominal_pinjaman"`
	JenisPinjaman   string    `json:"jenis_pinjaman"`
	// Customer        Customer
}

type Pegawai struct {
	ID        string
	FirstName string
	LastName  string
}

type pegawaiPortal struct {
	password string
}
