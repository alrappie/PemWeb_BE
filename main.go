package main

import (
	"PemWeb_BE/Database"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func main() {
	fmt.Println("HI")
	db := Database.Open()
	fmt.Println("Database terinisialisasi")

	r = gin.Default()
	r.Use(cors.Default())
	//router disini

	fmt.Println("Router siap")
	if err := r.Run(":5000"); err != nil {
		fmt.Println("error")
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Server berjalan")
}
