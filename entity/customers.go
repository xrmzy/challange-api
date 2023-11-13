package entity

import (
	"challange/config"
	"fmt"
	"log"
)

type Customer struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
}

// generate customer ID
func GenerateCustomerID() string {
	// randomNumber := GenerateRandomNumber(1000)
	// codePrefix := "CS"
	// return fmt.Sprintf("%s%03d", codePrefix, randomNumber)

	db, err := config.ConnectDb()
	if err != nil {
		log.Fatalf("Failed to Connect to the dabase: %v", err)
	}
	defer db.Close()

	var existingIDs []string
	rows, err := db.Query("SELECT cust_id FROM customers;")
	if err != nil {
		log.Fatalf("Failed to query customer IDs: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			log.Fatalf("Failed to scan customer IDs: %v", err)
		}
		existingIDs = append(existingIDs, id)
	}

	randomNumber := GenerateRandomNumber(1000)
	codePrefix := "CS"
	newID := fmt.Sprintf("%s%03d", codePrefix, randomNumber)

	// Pastikan ID yang dihasilkan adalah unik
	for {
		if !contains(existingIDs, newID) {
			break
		}
		randomNumber := GenerateRandomNumber(1000)
		newID = fmt.Sprintf("%s%03d", codePrefix, randomNumber)
	}

	return newID
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
