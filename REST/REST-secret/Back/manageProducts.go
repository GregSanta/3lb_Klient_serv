package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type product struct {
	gorm.Model

	ProductId    string
	ProductName  string
	ProductPrice float64 `gorm:"column:unit_price"`
}

func makeProductManager() *gorm.DB {
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: dbHandler.DB,
	}), &gorm.Config{})

	if err == nil {
		db.AutoMigrate(&product{})
		return db
	}

	return nil
}

func addProduct(ctx *gin.Context) {
	productId := ctx.Query("id")
	productName := ctx.Query("name")
	productPrice := ctx.Query("price")

	if productId == "" || productName == "" || productPrice == "" {
		ctx.AbortWithError(400, errors.New("Id, name and price must be included"))
		return
	}

	if price, err := strconv.ParseFloat(productPrice, 64); err == nil {
		product_mgr.Create(&product{ProductId: productId,
			ProductName:  productName,
			ProductPrice: price})
	} else {
		ctx.AbortWithError(400, errors.New("'price' field must be floating point number"))
		return
	}

	ctx.AbortWithStatus(http.StatusAccepted)
}

func editProduct(ctx *gin.Context) {
	productId := ctx.Query("id")
	if productId == "" {
		ctx.AbortWithError(400, errors.New("Id must be included"))
	}

	productName := ctx.Query("name")
	productPrice := ctx.Query("price")
	price, err := strconv.ParseFloat(productPrice, 64)

	var product product
	product_mgr.First(&product, "product_id = ?", productId)

	if productName != "" {
		product_mgr.Model(&product).Update("product_name", productName)
	}
	if productPrice != "" && err == nil {
		product_mgr.Model(&product).Update("unit_price", price)
	}

	ctx.AbortWithStatus(http.StatusAccepted)
}

func removeProduct(ctx *gin.Context) {
	productId := ctx.Query("id")
	if productId == "" {
		ctx.AbortWithError(400, errors.New("Id must be included"))
	}

	var product product
	product_mgr.First(&product, "product_id = ?", productId)

	product_mgr.Unscoped().Delete(&product)

	ctx.AbortWithStatus(http.StatusAccepted)
}
