package main

import (
    "fine-service/handlers"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    r.POST("/fines", handlers.CreateFine)
    r.GET("/fines", handlers.ListFines)

    r.Run(":3000")
}
