package entity

import (
	"database/sql"
	"fmt"
	"log"
)

func InsertCustomerSQL(db *sql.DB, customer *Customer) {
	sqlStatement := "INSERT INTO customers (cust_id, cust_name, phone_number, address, email) VALUES ($1, $2, $3, $4, $5);"
	_, err := db.Exec(sqlStatement, customer.Id, customer.Name, customer.PhoneNumber, customer.Address, customer.Email)
	if err != nil {
		log.Fatalf("Failed to Insert Customer data: %v", err)
	}
}

func UpdateCustomerSQL(db *sql.DB, customer *Customer) {
	sqlStatement := "UPDATE customers SET cust_name = $2, phone_number = $3, address = $4, email = $5 WHERE cust_id = $1;"
	_, err := db.Exec(sqlStatement, customer.Id, customer.Name, customer.PhoneNumber, customer.Address, customer.Email)
	if err != nil {
		log.Fatalf("Failed to Update Customer Data: %v", err)
	}
}

func DeleteCustomerSQL(db *sql.DB, customer *Customer) {
	sqlStatement := "DELETE FROM customers WHERE cust_id = $1;"
	_, err := db.Exec(sqlStatement, customer.Id)
	if err != nil {
		log.Fatalf("Failed to Delete Customer Data: %v", err)
	} else {
		fmt.Println("Successfully Delete Data!")
	}
}

func AddOrderSQL(tx *sql.Tx, order *Orders) {
	sqlStatement := "INSERT INTO orders (order_id, cust_id, cust_name, service, unit, outlet_name, order_date, status) VALUES ($1, (SELECT cust_id FROM customers WHERE cust_name = $2), $3, $4, $5, $6, $7, $8);"

	_, err := tx.Exec(sqlStatement, order.OrderId, order.CustomerName, order.CustomerName, order.Service, order.Unit, order.OutletName, order.OrderDate, order.Status)

	if err != nil {
		log.Fatalf("Failed to Add Order: %v \n", err)
		tx.Rollback()
	} else {
		fmt.Println("Successfully Added Order!")
		err = tx.Commit()
		if err != nil {
			log.Fatalf("Failed to Commit Transaction: %v\n", err)
		} else {
			fmt.Println("Transaction Committed Successfully")
		}
	}
}

func UpdateOrderSQL(tx *sql.Tx, order *Orders) {
	sqlStatement := "UPDATE orders SET cust_id =$2 ,cust_name = $3, service = $4, unit = $5, outlet_name = $6, order_date = $7, status = $8 WHERE order_id= $1;"

	_, err := tx.Exec(sqlStatement, order.OrderId, order.CustomerId, order.CustomerName, order.Service, order.Unit, order.OutletName, order.OrderDate, order.Status)

	if err != nil {
		log.Fatalf("Failed to Update Order: %v\n", err)
		tx.Rollback()
	} else {
		fmt.Println("Successfully Update Order!")
		err = tx.Commit()
		if err != nil {
			log.Fatalf("Failed to Commit Transaction: %v\n", err)
		} else {
			fmt.Println("Transaction Committed Successfully")
		}
	}
}

func DeleteOrderSQL(tx *sql.Tx, order *Orders) error {
	sqlStatement := "DELETE FROM orders WHERE order_id = $1;"
	_, err := tx.Exec(sqlStatement, order.OrderId)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func SearchOrderBySQL(db *sql.DB, order *Orders) error {
	sqlStatement := "SELECT order_id, cust_id, cust_name, service, unit, outlet_name, order_date, status FROM orders WHERE order_id = $1;"

	row := db.QueryRow(sqlStatement, order.OrderId)
	err := row.Scan(&order.OrderId, &order.CustomerId, &order.CustomerName, &order.Service, &order.Unit, &order.OutletName, &order.OrderDate, &order.Status)

	if err == sql.ErrNoRows {
		fmt.Println("Transaksi tidak ditemukan.")
		return nil
	} else if err != nil {
		return err
	}

	fmt.Println("Order ID :", order.OrderId)
	fmt.Println("Customer ID :", order.CustomerId)
	fmt.Println("Nama Customer :", order.CustomerName)
	fmt.Println("Service :", order.Service)
	fmt.Println("Unit :", order.Unit)
	fmt.Println("Outlet :", order.OutletName)
	fmt.Println("Tanggal Order:", order.OrderDate)
	fmt.Println("Status :", order.Status)

	return nil
}

func SearchCustByID(db *sql.DB, customer *Customer) error {
	sqlStatement := "SELECT cust_id, cust_name, phone_number, address, email FROM customers WHERE cust_id = $1;"
	row := db.QueryRow(sqlStatement, customer.Id)

	err := row.Scan(&customer.Id, &customer.Name, &customer.PhoneNumber, &customer.Address, &customer.Email)
	if err == sql.ErrNoRows {
		return fmt.Errorf("Pelanggan tidak ditemukan")
	} else if err != nil {
		return err
	}

	return nil
}
