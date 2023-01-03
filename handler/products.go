package handler

import (
	"Project_go/model"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB

func CreateProductHandler(c *gin.Context, db *gorm.DB) {

	var p model.Product
	if err := c.BindJSON(&p); err != nil {
		fmt.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := db.Create(&p).Error; err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, p)
}

func UpdateProductHandler(c *gin.Context, db *gorm.DB) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid product ID: %v", err))
		return
	}

	// Récupérez les données de la requête
	var p model.Product
	if err := c.BindJSON(&p); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Mise à jour du produit dans la base de données
	if err := db.Model(&p).Where("id = ?", id).Updates(p).Error; err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Renvoyez le produit mis à jour à l'utilisateur
	c.JSON(http.StatusOK, p)
}

func DeleteProductHandler(c *gin.Context, db *gorm.DB) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid product ID: %v", err))
		return
	}

	if err := db.Where("id = ?", id).Delete(&model.Product{}).Error; err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func GetProductByIDHandler(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var product model.Product
	if err := db.First(&product, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, product)
}

func GetAllProductsHandler(c *gin.Context, db *gorm.DB) {
	var products []model.Product
	if err := db.Find(&products).Error; err != nil {
		fmt.Print(err)
		return
	}
	c.JSON(http.StatusOK, products)
}
