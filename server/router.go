package main

import (
	h "main/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func route() {
	router = gin.Default()

	router.GET("/", h.DefaultHandler)
	router.GET("/imageDisplay", h.ImageDisplayHandler)

	http.Handle("/", router)
}