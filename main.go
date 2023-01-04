package main

import (
	"Project_go/handler"
	"Project_go/model"
	"Project_go/utils"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	var dbURL = "user:password@tcp(127.0.0.1:3306)/go-project?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dbURL), &gorm.Config{})
	if err != nil {
		fmt.Print(db)
	}
	db.AutoMigrate(&model.Product{}, &model.Payment{})

	router := gin.Default()
	router.POST("/products", func(c *gin.Context) {
		handler.CreateProductHandler(c, db)
	})
	router.PUT("/products/:id", func(c *gin.Context) {
		handler.UpdateProductHandler(c, db)
	})
	router.DELETE("/products/:id", func(c *gin.Context) {
		handler.DeleteProductHandler(c, db)
	})
	router.GET("/products/:id", func(c *gin.Context) {
		handler.GetProductByIDHandler(c, db)
	})
	router.GET("/products", func(c *gin.Context) {
		handler.GetAllProductsHandler(c, db)
	})
	router.POST("/payments", func(c *gin.Context) {
		handler.CreatePaymentHandler(c, db)
	})
	router.PUT("/payments/:id", func(c *gin.Context) {
		handler.UpdatePaymentHandler(c, db)
	})
	router.DELETE("/payments/:id", func(c *gin.Context) {
		handler.DeletePaymentHandler(c, db)
	})
	router.GET("/payments/:id", func(c *gin.Context) {
		handler.GetPaymentByIDHandler(c, db)
	})
	router.GET("/payments", func(c *gin.Context) {
		handler.GetAllPaymentsHandler(c, db)
	})

	router.GET("/payments/stream", func(c *gin.Context) {
		handler.StreamPaymentsHandler(c, utils.GetBroadcaster(), db)
	})

	router.Run(fmt.Sprintf(":%v", 8089))

}
