package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)



var dbHandler *sqlx.DB

var customer_mgr *gorm.DB

var order_mgr *gorm.DB

var product_mgr *gorm.DB

func init_database() *sqlx.DB {
	var link string = "postgres://nitzmmysllifcu:0b1616cc1b89676b8cb99ce7dd69deb6bd130d3bcff86f01545b176278413364@ec2-79-125-93-182.eu-west-1.compute.amazonaws.com:5432/d5lgfqm719kl7a"
	dbHandler, err := sqlx.Open("postgres", link)
	if err != nil {

		log.Fatal(err)

	}

	return dbHandler
}


func init_handlers(router *gin.Engine) {
	router.GET("/page_loaded", pageLoaded)
	router.GET("/products/fetch", fetchProducts)
	router.GET("/customers/fetch", fetchData)

	router.POST("/products/add", addProduct)
	router.POST("/products/remove", removeProduct)
	router.POST("/products/edit", editProduct)

	router.POST("/orders/add", addOrder)
	router.POST("/orders/remove", removeOrder)
	router.POST("/orders/edit", editOrder)

	router.POST("/customers/add", addCustomer)
	router.POST("/customers/remove", removeCustomer)
	router.POST("/customers/edit", editCustomer)
}

func main() {
	dbHandler = init_database()

	customer_mgr = makeCustomersManager()

	order_mgr = makeOrderManager()

	product_mgr = makeProductManager()

	defaultRouter := gin.Default()

	/*defaultRouter := Default()*/

	init_handlers(defaultRouter)

	defaultRouter.Run()
	defer dbHandler.Close()
}
