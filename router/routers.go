package router

import (
	"challange/handler"

	"github.com/gin-gonic/gin"
)

func SetupCustomerRoutes(router *gin.Engine) {
	customerRoutes := router.Group("/customers")
	{
		customerRoutes.GET("/", handler.GetCustomers)
		customerRoutes.GET("/:id", handler.GetCustomerByID)
		customerRoutes.POST("/create", handler.CreateCustomer)
		customerRoutes.PUT("/:id", handler.UpdateCustomer)
		customerRoutes.DELETE("/:id", handler.DeleteCustomer)
		// Tambahkan rute lainnya seperti POST, PUT, dan DELETE
	}
	orderRoutes := router.Group("/orders")
	{
		orderRoutes.GET("/", handler.GetOrders)
		orderRoutes.GET("/:id", handler.GetOrderById)
		orderRoutes.POST("/create", handler.CreateOrder)
		orderRoutes.PUT("/:id", handler.UpdateOrder)
		orderRoutes.DELETE("/:id", handler.DeleteOrder)
		// }
	}
}
