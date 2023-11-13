package entity

import (
	"bufio"
	"database/sql"

	// "enigma-laundry-console-api/entity"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

// genarete random ID
func GenerateRandomNumber(max int) int {
	return rand.Intn(max)
}

// validasi email
func IsEmailExists(db *sql.DB, email string) bool {
	var exists bool
	row := db.QueryRow("SELECT EXISTS (SELECT 1 FROM customers WHERE email = $1)", email)
	err := row.Scan(&exists)
	if err != nil {
		panic(err)
	}
	return exists
}

// read input
func ReadInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func IsCustomerExists(db *sql.DB, customerID string) bool {
	var exists bool
	row := db.QueryRow("SELECT EXISTS (SELECT 1 FROM customers WHERE cust_id = $1)", customerID)
	err := row.Scan(&exists)
	if err != nil {
		panic(err)
	}
	return exists
}

func GetCustomerIdByName(db *sql.DB, customerName string) (string, error) {
	var customerId string
	row := db.QueryRow("SELECT cust_id FROM customers WHERE cust_name = $1", customerName)
	err := row.Scan(&customerId)
	if err != nil {
		return "", err
	}
	return customerId, nil
}

func GetCustomerNameById(db *sql.DB, customerId string) (string, error) {
	var customerName string
	row := db.QueryRow("SELECT cust_name FROM customers WHERE cust_id = $1", customerId)
	err := row.Scan(&customerName)
	if err != nil {
		return "", err
	}
	return customerName, nil
}

func FindCustomerID(db *sql.DB, order *Orders) error {
	customerId, err := GetCustomerIdByName(db, order.CustomerName)
	if err != nil {
		fmt.Println("Data Pelanggan tidak ditemukan, Silahkan Tambah Pelanggan Terlebih Dahulu!")
		return err
	}
	order.CustomerId = customerId

	customerName, err := GetCustomerNameById(db, order.CustomerId)
	if err != nil {
		return err
	}
	order.CustomerName = customerName
	return nil
}
