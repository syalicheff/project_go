package handler

import (
	"Project_go/model"
	"Project_go/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

func CreatePaymentHandler(c *gin.Context, db *gorm.DB) {
	var p model.Payment
	if err := c.BindJSON(&p); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var product model.Product
	fmt.Print(p.ProductID)
	if err := db.First(&product, p.ProductID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Print(err)
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if p.PricePaid < product.Price {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("price paid (%v) must be greater than or equal to product price (%v)", p.PricePaid, product.Price))
		return
	}

	p.CreatedAt = time.Now()
	if err := db.Create(&p).Error; err != nil {
		fmt.Print(err)
		return
	}
	utils.GetBroadcaster().Broadcast(&p)

	c.JSON(http.StatusCreated, p)
}
func UpdatePaymentHandler(c *gin.Context, db *gorm.DB) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid payment ID: %v", err))
		return
	}

	var p model.Payment
	if err := c.BindJSON(p); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var product model.Product
	if err := db.First(&product, p.ProductID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if p.PricePaid < product.Price {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("price paid (%v) must be greater than or equal to product price (%v)", p.PricePaid, product.Price))
		return
	}

	p.UpdatedAt = time.Now()
	if err := db.Model(&model.Payment{ID: uint(id)}).Updates(p).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	utils.GetBroadcaster().Broadcast(&p)

	c.JSON(http.StatusOK, p)
}
func DeletePaymentHandler(c *gin.Context, db *gorm.DB) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid payment ID: %v", err))
		return
	}
	if err := db.Delete(&model.Payment{ID: uint(id)}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}
func GetPaymentByIDHandler(c *gin.Context, db *gorm.DB) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid payment ID: %v", err))
		return
	}
	var p model.Payment
	if err := db.First(&p, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, p)
}
func GetAllPaymentsHandler(c *gin.Context, db *gorm.DB) {
	var payments []model.Payment
	if err := db.Find(&payments).Error; err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, payments)
}
func StreamPaymentsHandler(c *gin.Context) {
	client := utils.GetBroadcaster().Subscribe()
	defer utils.GetBroadcaster().Unsubscribe(client)

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	// Send broadcast events to the client
	for {
		select {
		case event := <-client:
			if p, ok := event.(*model.Payment); ok {
				c.SSEvent("payment", p)
			}
		case <-c.Request.Context().Done():
			return
		}
	}
}
