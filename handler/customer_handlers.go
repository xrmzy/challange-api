package handler

import (
	// "enigma-laundry-console-api/config"
	// "enigma-laundry-console-api/entity"
	"challange/config"
	"challange/entity"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCustomers(c *gin.Context) {
	db, err := config.ConnectDb()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal terhubung ke database"})
		return
	}
	defer db.Close()

	// Logika untuk mengambil semua pelanggan dari database
	var customers []entity.Customer

	rows, err := db.Query("SELECT cust_id, cust_name, phone_number, address, email FROM customers;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pelanggan"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var customer entity.Customer
		err := rows.Scan(&customer.Id, &customer.Name, &customer.PhoneNumber, &customer.Address, &customer.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca data pelanggan"})
			return
		}
		customers = append(customers, customer)
	}

	c.JSON(http.StatusOK, customers)
}

func GetCustomerByID(c *gin.Context) {
	// Implementasi logika untuk mengambil data pelanggan berdasarkan ID dari database
	db, err := config.ConnectDb()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal terhubung ke database"})
		return
	}
	defer db.Close()

	id := c.Param("id")
	fmt.Println("ID:", id)
	var customer entity.Customer

	sqlStatement := ("SELECT cust_id, cust_name, phone_number, address, email FROM customers WHERE cust_id = $1;")
	fmt.Println("Query:", sqlStatement)

	err = db.QueryRow(sqlStatement, id).Scan(
		&customer.Id, &customer.Name, &customer.PhoneNumber, &customer.Address, &customer.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pelanggan tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

func CreateCustomer(c *gin.Context) {
	db, err := config.ConnectDb()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal terhubung ke database"})
		return
	}
	defer db.Close()

	// Mendapatkan data pelanggan dari permintaan (request body)
	var customer entity.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate ID pelanggan
	customer.Id = entity.GenerateCustomerID()

	// Periksa apakah ID pelanggan sudah ada dalam database
	queryCheckID := "SELECT COUNT(*) FROM customers WHERE cust_id = $1"
	var count int
	err = db.QueryRow(queryCheckID, customer.Id).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memeriksa ID pelanggan"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID pelanggan sudah ada dalam database"})
		return
	}

	// Lakukan penambahan data pelanggan ke dalam database
	queryInsert := `
    INSERT INTO customers (cust_id, cust_name, phone_number, address, email)
    VALUES ($1, $2, $3, $4, $5);
    `

	_, err = db.Exec(queryInsert, customer.Id, customer.Name, customer.PhoneNumber, customer.Address, customer.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan pelanggan"})
		return
	}

	// Setelah berhasil menambahkan pelanggan, Anda bisa memberikan respons sukses
	c.JSON(http.StatusCreated, gin.H{"message": "Pelanggan berhasil ditambahkan", "cust_id": customer.Id})
}

func UpdateCustomer(c *gin.Context) {
	// Implementasi logika untuk mengubah data pelanggan di database
	db, err := config.ConnectDb()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal terhubung ke database"})
		return
	}
	defer db.Close()

	id := c.Param("id")

	// Mendapatkan data pelanggan dari permintaan (request body)
	var updatedCustomer entity.Customer
	if err := c.ShouldBindJSON(&updatedCustomer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sqlStatement := `
        UPDATE customers
        SET cust_name = $2, phone_number = $3, address = $4, email = $5
        WHERE cust_id = $1;
    `
	_, err = db.Exec(sqlStatement, id, updatedCustomer.Name, updatedCustomer.PhoneNumber, updatedCustomer.Address, updatedCustomer.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data pelanggan"})
		return
	}

	// Setelah berhasil mengubah data pelanggan, Anda bisa memberikan respons sukses
	c.JSON(http.StatusOK, gin.H{"message": "Data pelanggan berhasil diperbarui"})
}

func DeleteCustomer(c *gin.Context) {
	db, err := config.ConnectDb()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal terhubung ke database"})
		return
	}
	defer db.Close()

	id := c.Param("id")

	// Periksa apakah pelanggan dengan ID yang diberikan ada dalam basis data
	var existingID string
	err = db.QueryRow("SELECT cust_id FROM customers WHERE cust_id = $1", id).Scan(&existingID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pelanggan tidak ditemukan"})
		return
	}

	// Lakukan penghapusan data pelanggan dari database
	_, err = db.Exec("DELETE FROM customers WHERE cust_id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus pelanggan"})
		return
	}

	// Setelah berhasil menghapus pelanggan, Anda bisa memberikan respons sukses
	c.JSON(http.StatusOK, gin.H{"message": "Pelanggan berhasil dihapus"})
}
