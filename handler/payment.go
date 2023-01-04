package handler

import (
	"Project_go/model"
	"Project_go/utils"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
func StreamPaymentsHandler(c *gin.Context, broadcaster *utils.Broadcaster, db *gorm.DB) {
	client := broadcaster.Subscribe()
	defer broadcaster.Unsubscribe(client)

	c.Header("Content-Type", "text/html")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	// how to serve style.css file ?

	// Send header of table
	c.Writer.Write([]byte(`
			<style>
			table {
				border-collapse: collapse;

			}
			th, td {
				border: 1px solid black;
				padding: 5px;
			}

			tr:nth-child(even) {
				background-color: #eee;
			}
			tr:nth-child(odd) {
				background-color: #fff;
			}

			th {
				background-color: black;
				color: white;
			}
			
		</style>
	`))
	c.Writer.Write([]byte("<table><th>Product Name</th><th>Total Price</th><th>Payment Date</th></tr>"))

	// Send broadcast events to the client
	for {
		select {
		case event := <-client:
			// Write row of table for each payment event
			dateTimeFormat := "2006-01-02 at 15:04"
			// replace the event.ProductID by its name
			var product model.Product
			if err := db.First(&product, event.ProductID).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					c.AbortWithStatus(http.StatusNotFound)
					return
				}
				_ = c.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			c.Writer.Write([]byte(fmt.Sprintf("<tr><td>%s</td><td>%d</td><td>%s</td></tr>", product.Name, event.PricePaid, event.CreatedAt.Format(dateTimeFormat))))
			c.Writer.Flush()
		case <-c.Request.Context().Done():
			// Close table when done
			c.Writer.Write([]byte("</table>"))
			return
		}
	}
}
