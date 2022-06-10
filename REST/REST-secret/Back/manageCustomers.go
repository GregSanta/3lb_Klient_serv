package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type customer struct {
	gorm.Model

	CustomerId  string
	CompanyName string
	ContactName string
	Address     string
	City        string
	Country     string
	Phone       string
	Fax         string
}

func makeCustomersManager() *gorm.DB {
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: dbHandler.DB,
	}), &gorm.Config{})

	if err == nil {
		db.AutoMigrate(&customer{})
		return db
	}

	return nil
}

func addCustomer(ctx *gin.Context) {
	var newCustomer customer
	newCustomer.Address = ctx.Query("address")
	newCustomer.City = ctx.Query("city")
	newCustomer.CompanyName = ctx.Query("company")
	newCustomer.ContactName = ctx.Query("contact")
	newCustomer.Country = ctx.Query("country")
	newCustomer.CustomerId = ctx.Query("id")
	newCustomer.Phone = ctx.Query("phone")
	newCustomer.Fax = ctx.Query("fax")

	if newCustomer.Address == "" ||
		newCustomer.City == "" ||
		newCustomer.CompanyName == "" ||
		newCustomer.ContactName == "" ||
		newCustomer.Country == "" ||
		newCustomer.CustomerId == "" ||
		newCustomer.Phone == "" ||
		newCustomer.Fax == "" {
		ctx.AbortWithError(400, errors.New("Id, company, contact, address, city, country, phone, fax must be included"))
		return
	}

	customer_mgr.Create(&newCustomer)
	ctx.AbortWithStatus(http.StatusAccepted)
}

func editCustomer(ctx *gin.Context) {
	Address := ctx.Query("address")
	City := ctx.Query("city")
	CompanyName := ctx.Query("company")
	ContactName := ctx.Query("contact")
	Country := ctx.Query("country")
	CustomerId := ctx.Query("id")
	Phone := ctx.Query("phone")
	Fax := ctx.Query("fax")

	if CustomerId == "" {
		ctx.AbortWithError(400, errors.New("Id must be included"))
		return
	}

	var customer customer
	customer_mgr.First(&customer, "customer_id = ?", CustomerId)

	if Address != "" {
		customer_mgr.Model(&customer).Update("address", Address)
	}
	if City != "" {
		customer_mgr.Model(&customer).Update("city", City)
	}
	if CompanyName != "" {
		customer_mgr.Model(&customer).Update("company_name", CompanyName)
	}
	if ContactName != "" {
		customer_mgr.Model(&customer).Update("contact_name", ContactName)
	}
	if Country != "" {
		customer_mgr.Model(&customer).Update("country", Country)
	}
	if Phone != "" {
		customer_mgr.Model(&customer).Update("phone", Phone)
	}
	if Fax != "" {
		customer_mgr.Model(&customer).Update("fax", Fax)
	}

	ctx.AbortWithStatus(http.StatusAccepted)
}

func removeCustomer(ctx *gin.Context) {
	CustomerId := ctx.Query("id")
	if CustomerId == "" {
		ctx.AbortWithError(401, errors.New("Id must be included"))
	}
	var customer customer
	customer_mgr.First(&customer, "customer_id = ?", CustomerId)

	customer_mgr.Unscoped().Delete(&customer)
}
