package main

import (
	"challange/router"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	router.SetupCustomerRoutes(r)

	r.Run(":8080")
}
