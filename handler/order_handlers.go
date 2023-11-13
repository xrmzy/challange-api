package handler

import (
	"challange/config"
	"challange/entity"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetOrders(c *gin.Context) {
	db, err := config.ConnectDb()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal terhubung ke database"})
		return
	}
	defer db.Close()

	var orders []entity.Orders
	rows, err := db.Query("SELECT order_id, cust_id, cust_name, service, unit, outlet_name, order_date, status FROM orders;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data Transaksi"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var order entity.Orders
		err := rows.Scan(&order.OrderId, &order.CustomerId, &order.CustomerName, &order.Service, &order.Unit, &order.OutletName, &order.OrderDate, &order.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca data Transaksi"})
			return
		}
		orders = append(orders, order)
	}
	c.JSON(http.StatusOK, orders)
}

func GetOrderById(c *gin.Context) {
	db, err := config.ConnectDb()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal terhubung ke database"})
		return
	}
	defer db.Close()

	id := c.Param("id")
	fmt.Println("ID:", id)

	var order entity.Orders
	sqlStatement := ("SELECT order_id, cust_id, cust_name, service, unit, outlet_name, order_date, status FROM orders WHERE order_id = $1;")
	fmt.Println("Query:", sqlStatement)

	err = db.QueryRow(sqlStatement, id).Scan(
		&order.OrderId, &order.CustomerId, &order.CustomerName, &order.Service, &order.Unit, &order.OutletName, &order.OrderDate, &order.Status)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pelanggan tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func CreateOrder(c *gin.Context) {
	db, err := config.ConnectDb()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal terhubung ke database"})
		return
	}
	defer db.Close()

	var order entity.Orders
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Dapatkan cust_name dari data pesanan
	customerName := order.CustomerName

	// Cari cust_id berdasarkan cust_name dalam tabel customers
	var customerID string
	err = db.QueryRow("SELECT cust_id FROM customers WHERE cust_name = $1", customerName).Scan(&customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencari ID Pelanggan"})
		return
	}

	// generate ID Order
	order.OrderId = entity.GenerateOrderID()

	queryCheckID := "SELECT COUNT(*) FROM orders WHERE order_id = $1"
	var count int
	err = db.QueryRow(queryCheckID, order.OrderId).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memeriksa ID Transaksi"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID Transaksi sudah ada dalam database"})
		return
	}

	queryInsert := `
   INSERT INTO orders (order_id, cust_id, cust_name, service, unit, outlet_name, order_date, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
    `

	_, err = db.Exec(queryInsert, &order.OrderId, &customerID, &order.CustomerName, &order.Service, &order.Unit, &order.OutletName, &order.OrderDate, &order.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan Transaksi"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Transaksi berhasil ditambahkan", "order_id": order.OrderId})
}

func UpdateOrder(c *gin.Context) {
	db, err := config.ConnectDb()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal terhubung ke database"})
		return
	}
	defer db.Close()

	id := c.Param("id")

	var updateOrder entity.Orders
	if err := c.ShouldBindJSON(&updateOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sqlStatement := `
        UPDATE orders SET cust_name = $1, service = $2, unit = $3, outlet_name = $4, order_date = $5, status = $6 WHERE order_id= $7;
    `
	_, err = db.Exec(sqlStatement, updateOrder.CustomerName, updateOrder.Service, updateOrder.Unit, updateOrder.OutletName, updateOrder.OrderDate, updateOrder.Status, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data Transaksi"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data Transaksi berhasil diperbarui"})
}

func DeleteOrder(c *gin.Context) {
	db, err := config.ConnectDb()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal terhubung ke database"})
		return
	}
	defer db.Close()

	id := c.Param("id")

	var existingID string
	err = db.QueryRow("SELECT order_id FROM orders WHERE order_id = $1", id).Scan(&existingID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaksi tidak ditemukan"})
		return
	}

	_, err = db.Exec("DELETE FROM orders WHERE order_id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus Transaksi"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Transaksi berhasil dihapus"})
}
