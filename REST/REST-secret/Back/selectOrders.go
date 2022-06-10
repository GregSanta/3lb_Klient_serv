package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cznic/mathutil"
	"github.com/gin-gonic/gin"
)

type Order struct {
	Order_id      int       `json:"id" xml:"id" db:"order_id"`
	Discount      float64   `json:"dicount" xml:"dicount" db:"discount"`
	Order_date    time.Time `json:"order_date" xml:"order_date" db:"order_date"`
	Delivery_date time.Time `json:"shipped_date" xml:"shipped_date" db:"shipped_date"`
	Products      []Product
}

type Customer struct {
	Id           string  `json:"id" xml:"id" db:"customer_id"`
	Company      string  `json:"company" xml:"company" db:"company_name"`
	Contact_name string  `json:"contact" xml:"contact" db:"contact_name"`
	Address      *string `json:"address" xml:"address" db:"address"`
	City         *string `json:"city" xml:"city" db:"city"`
	Phone        string  `json:"phone" xml:"phone" db:"phone"`
	Fax          *string `json:"fax" xml:"fax" db:"fax"`

	Orders *[]Order
}

func getMainInfo(clientId string) Customer {
	var customer Customer
	dbHandler.Get(&customer, "SELECT customer_id, company_name, contact_name, phone FROM customers WHERE customer_id = $1", clientId)

	return customer
}

func getFullInfo(clientId string) Customer {
	var customer Customer
	fmt.Print(dbHandler.Get(&customer, "select customer_id, company_name, contact_name, address, city, phone, fax from customers where customer_id = $1", clientId))

	var orders []Order
	fmt.Print(dbHandler.Select(&orders, "select order_id from orders where customer_id = $1", clientId))
	for i, unit := range orders {
		fmt.Print(dbHandler.Get(&orders[i], `select order_date, shipped_date from orders where order_id = $1`, orders[i].Order_id))
		var products []Product
		fmt.Print(dbHandler.Select(&products, `select product_id, product_name, units_in_stock, unit_price, sum(quantity) as quantity from products natural join order_details where product_id in (select product_id from order_details where order_id = $1) group by product_id, product_name, units_in_stock, unit_price`, unit.Order_id))
		orders[i].Products = products
	}
	customer.Orders = &orders

	return customer
}
/*func getFullInfo(clientId string) Customer {
	var customer Customer
	fmt.Print(dbHandler.Get(&, ", fax from customers where customer_id = $1", clientId))

	var orders []Order
	fmt.Print(dbHandler.Select(&orders, "select order_id from orders where customer_id = $1", clientId))
	for i, unit := range orders {
		fmt.Print(dbHandler.Get(&orders[i], `select order_date, shipped_date from orders where order_id = $1`, orders[i].Order_id))
		var products []Product
	
	}
	customer.Orders = &orders

	return customer
}*/

func getPriceList(page, pagesize int) []Product {
	var products []Product
	dbHandler.Select(&products, "SELECT product_id, product_name, unit_price FROM products limit $1 offset $2", pagesize, pagesize*(page-1))

	return products
}

func getOrders(date_from, date_to string, page, pagesize int) []Product {
	var products []Product

	if date_from == "" || date_to == "" {
		return nil
	}

	dbHandler.Select(&products, "SELECT product_id, product_name, quantity FROM products NATURAL JOIN order_details NATURAL JOIN orders where order_date >= to_date($1, 'YYYY-MM-DD') and order_date < to_date($2, 'YYYY-MM-DD') limit $3 offset $4", date_from, date_to, pagesize, pagesize*(page-1))

	return products
}

func fetchData(ctx *gin.Context) {
	requestedInfo := ctx.Query("info")
	format := ctx.Query("format")
	page, pageerr := strconv.Atoi(ctx.Query("page"))
	pagesize, pgsizeerr := strconv.Atoi(ctx.Query("pagesize"))

	if pageerr != nil {
		page = 1
	}
	if pgsizeerr != nil {
		pagesize = 20
	} else {
		pagesize = mathutil.Clamp(pagesize, 1, 50)
	}

	switch requestedInfo {
	case "main":
		clientId := ctx.Query("id")
		customer := getMainInfo(clientId)
		sendResult(nil, ctx, customer, nil, format)
		return

	case "full":
		clientId := ctx.Query("id")
		customer := getFullInfo(clientId)
		sendResult(nil, ctx, customer, nil, format)
		return

	case "price":
		products := getPriceList(page, pagesize)
		sendResult(nil, ctx, products, nil, format)
		return

	case "order":
		date_from := ctx.Query("date_from")
		date_to := ctx.Query("date_to")
		orders := getOrders(date_from, date_to, page, pagesize)
		sendResult(nil, ctx, orders, nil, format)
		return

	default:
		ctx.AbortWithStatus(404)
	}
}
