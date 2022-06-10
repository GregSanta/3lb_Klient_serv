package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type order struct {
	gorm.Model

	OrderId      string
	CustomerId   string
	OrderDate    time.Time
	DeliveryDate time.Time `gorm:"column:shipped_date"`
}

var dateFormat string = "2006-01-02"

func makeOrderManager() *gorm.DB {
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: dbHandler.DB,
	}), &gorm.Config{})

	if err == nil {
		db.AutoMigrate(&order{})
		return db
	}

	return nil
}

func addOrder(ctx *gin.Context) {
	orderId := ctx.Query("order_id")
	customerId := ctx.Query("customer_id")
	orderDate, deliveryDate := ctx.Query("order_date"), ctx.Query("delivery_date")
	orderParsed, orderErr := time.Parse(dateFormat, orderDate)
	deliveryParsed, deliveryErr := time.Parse(dateFormat, deliveryDate)

	if orderErr == nil && deliveryErr == nil {
		order_mgr.Create(&order{
			OrderId:      orderId,
			CustomerId:   customerId,
			OrderDate:    orderParsed,
			DeliveryDate: deliveryParsed,
		})
		ctx.AbortWithStatus(http.StatusAccepted)
	} else {
		ctx.AbortWithError(400, errors.New("Date format: YYYY-MM-DD"))
		return
	}
}

func editOrder(ctx *gin.Context) {
	orderId := ctx.Query("order_id")
	customerId := ctx.Query("customer_id")
	orderDate, deliveryDate := ctx.Query("order_date"), ctx.Query("delivery_date")
	orderParsed, orderErr := time.Parse(dateFormat, orderDate)
	deliveryParsed, deliveryErr := time.Parse(dateFormat, deliveryDate)

	var order order
	order_mgr.First(&order, "order_id = ?", orderId)

	if customerId != "" {
		order_mgr.Model(&order).Update("customer_id", customerId)
	}
	if orderDate != "" && orderErr == nil {
		order_mgr.Model(&order).Update("order_date", orderParsed)
	}
	if deliveryDate != "" && deliveryErr == nil {
		order_mgr.Model(&order).Update("shipped_date", deliveryParsed)
	}

	ctx.AbortWithStatus(http.StatusAccepted)
}

func removeOrder(ctx *gin.Context) {
	orderId := ctx.Query("order_id")
	var order order
	order_mgr.First(&order, "order_id = ?", orderId)

	order_mgr.Unscoped().Delete(&order)
}
