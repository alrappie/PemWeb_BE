package main

import (
	"PemWeb_BE/Database"
	"PemWeb_BE/User"
	"fmt"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func main() {
	fmt.Println("HI")
	db := Database.Open()
	fmt.Println("Database terinisialisasi")

	r = gin.Default()
	r.Use(func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    })

	//router disini
	User.Routes(db, r)

	fmt.Println("Router siap")
	if err := r.Run(":5000"); err != nil {
		fmt.Println("error")
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Server berjalan")
}
