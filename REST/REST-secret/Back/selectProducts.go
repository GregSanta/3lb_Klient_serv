package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jmoiron/sqlx"
)

type Product struct {
	Product_id       int      `json:"id" xml:"id" db:"product_id"`
	Product_name     string   `json:"name" xml:"name" db:"product_name"`
	Product_category *string  `json:"category" xml:"category" db:"category_name"`
	Product_quantity *int     `json:"quantity" xml:"quantity" db:"units_in_stock"`
	Product_price    *float64 `json:"price" xml:"price" db:"unit_price"`

	Order_date     *time.Time `json:"date" xml:"date" db:"order_date"`
	Order_quantity *int       `json:"order_quantity" xml:"order_quantity" db:"quantity"`
	MonthSale      *int       `json:"month" xml:"month" db:"month"`
}

type TableResponse struct {
	Headers []string    `json:"headers" xml:"headers"`
	Content interface{} `json:"content" xml:"content"`
}

func sendResult(err error, ctx *gin.Context, productList interface{}, headers []string, format string) {
	if err == nil {
		var wholeTable TableResponse
		wholeTable.Content, wholeTable.Headers = productList, headers
		switch strings.ToLower(format) {
		case "xml":
			ctx.XML(200, wholeTable)
			return

		default:
			ctx.JSON(200, wholeTable)
		}

	} else {
		ctx.AbortWithError(http.StatusBadRequest, err)
	}
}

func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func makeQuery(category string, id string, year string, dateFrom string, dateTo string, sort string) (string, []string) {
	selects := []string{"product_name", "product_id", "category_name", "products.unit_price"}
	headers := []string{"Название", "Идентификатор", "Категория", "Цена"}
	joins := []string{"categories"}
	usings := []string{"category_id"}
	terms := []string{}
	groups := []string{}
	sorts := []string{}
	var way string

	if category != "" {
		headers = append(headers, "На складе")
		selects = append(selects, "units_in_stock")
		terms = append(terms, fmt.Sprintf("category_id = %s", category))
	}
	if id != "" {
		selects = append(selects, "units_in_stock")
		headers = append(headers, "На складе")
		terms = append(terms, fmt.Sprintf("product_id = %s", id))
	}
	if dateFrom != "" && dateTo != "" {
		headers = append(headers, "Заказанное количество")
		selects = append(selects, "quantity")
		terms = append(terms, fmt.Sprintf("order_date >= to_date('%s', 'YYYY-MM-DD')", dateFrom), fmt.Sprintf("order_date < to_date('%s', 'YYYY-MM-DD')", dateTo))
		joins = append(joins, "order_details", "orders")
		usings = append(usings, "product_id", "order_id")
	}
	if year != "" {
		joins = append(joins, "order_details", "orders")
		usings = append(usings, "product_id", "order_id")
		selects = append(selects, "extract(month from order_date) as month")
		for _, elem := range selects {
			word := strings.Split(elem, " ")
			groups = append(groups, word[len(word)-1])
		}
		selects = append(selects, "sum(quantity) as quantity")
		terms = append(terms, fmt.Sprintf("extract(year from order_date) = %s", year))
		headers = append(headers, "Месяц", "Количество")
	}
	if match, _ := regexp.MatchString("[+|-]\\w+", sort); match {
		if sort[0] == '+' {
			way = "Asc"
		} else {
			way = "Desc"
		}

		sorts = append(sorts, sort[1:])
	}

	var resultQuery string
	resultQuery += "select " + strings.Join(selects, ", ") + " from products"

	for i := 0; i < len(joins); i++ {
		joins[i] = " join " + joins[i] + fmt.Sprintf(" using(%s) ", usings[i])
	}
	resultQuery += strings.Join(joins, " ")

	if len(terms) != 0 {
		resultQuery += " where " + strings.Join(terms, " and ")
	}

	if len(groups) != 0 {
		resultQuery += " group by " + strings.Join(groups, ", ")
	}

	if len(sorts) != 0 {
		resultQuery += " order by " + strings.Join(sorts, ", ") + " " + way
	}

	return resultQuery, headers
}

func fetchProducts(ctx *gin.Context) {

	productList := []Product{}
	format := ctx.Query("format")
	yearToFetch := ctx.Query("year")
	productIdToFetch := ctx.Query("id")
	categoryToFetch := ctx.Query("category")
	dateFrom, dateTo := ctx.Query("date_from"), ctx.Query("date_to")
	sort := ctx.Query("sort")
	query, headers := makeQuery(categoryToFetch, productIdToFetch, yearToFetch, dateFrom, dateTo, sort)

	err := dbHandler.Select(&productList, query)
	fmt.Print(err)
	fmt.Print(query)
	sendResult(err, ctx, productList, headers, format)
}
