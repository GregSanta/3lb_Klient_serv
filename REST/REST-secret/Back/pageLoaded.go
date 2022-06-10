package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/jmoiron/sqlx"
)


type Category struct {
	Category_id   int    `json:"id" db:"category_id"`
	Category_name string `json:"name" db:"category_name"`
}

func pageLoaded(ctx *gin.Context) {
	
	categoryList := []Category{}
	err := dbHandler.Select(&categoryList, "SELECT category_id, category_name FROM categories")

	
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	} else {
		ctx.JSON(200, categoryList)
	}
}
