package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB
var e error

type Customer struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

func main() {
	db, e = gorm.Open("sqlite3", "./example.db")
	if e != nil {
		fmt.Println(e)
	}
	defer db.Close()

	db.AutoMigrate(&Customer{})

	r := gin.Default()
	// Get customers
	r.GET("/customers", getCustomers)
	// Get customer by id
	r.GET("/customers/:id", getCustomerById)
	// Insert new customer
	r.POST("/customers", insertCustomer)
	// Update customer
	r.PUT("/customers/:id", updateCustomer)
	// Delete customer
	r.DELETE("/customers/:id", deleteCustomer)
	r.Run(":1991")
}

// Get customers
func getCustomers(c *gin.Context) {
	var customers []Customer
	if e := db.Find(&customers).Error; e != nil {
		c.AbortWithStatus(404)
		fmt.Println(e)
	} else {
		c.JSON(200, customers)
	}
}

// Get customer by id
func getCustomerById(c *gin.Context) {
	var customer Customer
	id := c.Params.ByName("id")
	if e := db.Where("id = ?", id).First(&customer).Error; e != nil {
		c.AbortWithStatus(404)
		fmt.Println(e)
	} else {
		c.JSON(200, customer)
	}
}

// Insert new customer
func insertCustomer(c *gin.Context) {
	var customer Customer
	c.BindJSON(&customer)
	db.Create(&customer)
	c.JSON(200, customer)
}

// Update customer
func updateCustomer(c *gin.Context) {
	var customer Customer
	id := c.Params.ByName("id")
	if e := db.Where("id = ?", id).First(&customer).Error; e != nil {
		c.AbortWithStatus(404)
		fmt.Println(e)
	} else {
		c.BindJSON(&customer)
		db.Save(&customer)
		c.JSON(200, customer)
	}
}

// Delete customer
func deleteCustomer(c *gin.Context) {
	var customer Customer
	id := c.Params.ByName("id")
	d := db.Where("id = ?", id).Delete(&customer)
	fmt.Println(d)
	c.JSON(200, gin.H{"id #" + id: "deleted"})
}
