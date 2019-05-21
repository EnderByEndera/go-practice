package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// Add Route here to create RESTful API
func initRoute() {
	router = gin.Default()

	v1 := router.Group("/go")
	v1.POST("/getdata", getData)
}

func getData(c *gin.Context) {
	param1 := c.PostForm("ratio")
	param2 := c.PostForm("csvchoice")
	r, err := strconv.ParseFloat(param1, 32)
	choice, err := strconv.ParseBool(param2)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": err,
			"ratio": param1})
	} else {
		resources, answer := realize(r, choice)
		c.JSON(http.StatusOK, gin.H{
			"path":   resources,
			"length": answer,
			"ratio":  param1})
	}
}
