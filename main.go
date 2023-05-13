package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"localhost/labback/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	err := models.ConnectDatabase()
	checkErr(err)

	r := gin.Default()

	// CORS for https://foo.com and https://github.com origins, allowing:
	// - PUT and PATCH methods
	// - Origin header
	// - Credentials share
	// - Preflight requests cached for 12 hours
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		// AllowOriginFunc: func(origin string) bool {
		// 	return origin == "https://github.com"
		// },
		MaxAge: 12 * time.Hour,
	}))

	// API v1
	v1 := r.Group("/api/v1")
	{
		v1.GET("/", alive)
		v1.GET("bank", getBanks)
		v1.GET("bank/:id", getBankByID)
		v1.POST("bank", addBank)
		v1.PUT("bank/:id", updateBank)
		v1.DELETE("bank/:id", deleteBank)
		v1.OPTIONS("bank", options)
	}

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	log.Fatal(r.Run(":9090"))
}

func alive(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "alive"})
}

func getBanks(c *gin.Context) {

	banks, err := models.GetBanks(10)

	checkErr(err)

	if banks == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": banks})
	}
}

func getBankByID(c *gin.Context) {

	// grab the Id of the record want to retrieve
	id := c.Param("id")

	bank, err := models.GetBankByID(id)

	checkErr(err)
	// if the name is blank we can assume nothing is found
	if bank.Bank == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": bank})
	}
}

func addBank(c *gin.Context) {

	var json models.Bank

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	success, err := models.AddBank(json)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func updateBank(c *gin.Context) {

	var json models.Bank

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bankId, err := strconv.Atoi(c.Param("id"))

	fmt.Printf("Updating id %d", bankId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	success, err := models.UpdateBank(json, bankId)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func deleteBank(c *gin.Context) {

	bankId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	success, err := models.DeleteBank(bankId)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func options(c *gin.Context) {

	ourOptions := "HTTP/1.1 200 OK\n" +
		"Allow: GET,POST,PUT,DELETE,OPTIONS\n" +
		"Access-Control-Allow-Origin: http://127.0.0.1:8080\n" +
		"Access-Control-Allow-Methods: GET,POST,PUT,DELETE,OPTIONS\n" +
		"Access-Control-Allow-Headers: Content-Type\n"

	c.String(200, ourOptions)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
