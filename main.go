package main

import (

	// library to build HTTP client and server
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	// helps building web applications
	"github.com/gin-gonic/gin"

	// package to handle uuid
	"github.com/google/uuid"
)

/*
Dictionary

	Key[UUID]: Value[Receipt]
*/
var receipt_map map[string]Receipt

/*
Data models
	- Receipt json data
*/

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`

	// better to define it separatly but ok for this example
	Items []struct {
		ShortDescription string `json:"shortDescription"`
		Price            string `json:"price"`
	} `json:"items"`

	Total string `json:"total"`
}

/*
Function to generate UUID and store the receipt json data into dictionary
*/

func postProcessReceipts(context *gin.Context) {
	var newReceipt Receipt

	if error := context.BindJSON(&newReceipt); error != nil {
		return
	}

	// generate a new uuid, convert it to string and use it as a key to store the data
	id := (uuid.New()).String()
	receipt_map[id] = newReceipt

	// response
	context.IndentedJSON(http.StatusOK, gin.H{
		"id": id,
	})

}

/*
Function to get the data from the dictionary provided UUID.
Calculate and return the points based on the recieved data
*/

func getPoints(context *gin.Context) {
	id := context.Param("id")

	if _, ok := receipt_map[id]; ok {
		// id exists in the map
		context.IndentedJSON(http.StatusOK, gin.H{"points": calculate(receipt_map[id])})
	} else {
		// id does not exist in the map
		context.IndentedJSON(http.StatusNotFound, gin.H{"error": "Id not found in the database"})
	}
}

/*
calculate the points accoring to the respective json to the passed id
*/

func calculate(receipt_map Receipt) int {
	// total 7 conditions

	fmt.Printf("alphanumeric: %d \n roundDollarAmount: %d \n multiple: %d \n everyTwo: %d \n trimmedLength: %d \n oddDate: %d \n afterTime: %d \n", alphanumeric(receipt_map), roundDollarAmount(receipt_map), multiple(receipt_map), everyTwo(receipt_map), trimmedLength(receipt_map), oddDate(receipt_map), afterTime(receipt_map))

	return alphanumeric(receipt_map) + roundDollarAmount(receipt_map) + multiple(receipt_map) + everyTwo(receipt_map) + trimmedLength(receipt_map) + oddDate(receipt_map) + afterTime(receipt_map)
}

// One point for every alphanumeric character in the retailer name.
func alphanumeric(receipt_map Receipt) int {

	reg, err := regexp.Compile("[^a-zA-Z0-9]+")

	if err != nil {
		fmt.Println(err)
	}
	processedString := reg.ReplaceAllString(receipt_map.Retailer, "")

	return len(processedString)

	// return len(strings.ReplaceAll(receipt_map.Retailer, " ", ""))

}

// 50 points if the total is a round dollar amount with no cents.
func roundDollarAmount(receipt_map Receipt) int {
	var total, err = strconv.ParseFloat(receipt_map.Total, 64)

	if err != nil {
		return 0
	}

	if math.Mod(total, 1) == 0 {
		return 50
	}

	return 0
}

// 25 points if the total is a multiple of `0.25`.
func multiple(receipt_map Receipt) int {

	var total, err = strconv.ParseFloat(receipt_map.Total, 64)

	if err != nil {
		return 0
	}

	if math.Mod(total, 0.25) == 0 {
		return 25
	}

	return 0
}

// 5 points for every two items on the receipt.
func everyTwo(receipt_map Receipt) int {
	return len(receipt_map.Items) / 2 * 5
}

// If the trimmed length of the item description is a multiple of 3, multiply the price by `0.2` and round up to the nearest integer. The result is the number of points earned.
func trimmedLength(receipt_map Receipt) int {
	var final int = 0
	for _, item := range receipt_map.Items {
		trimmedLength := len(strings.TrimSpace(item.ShortDescription))
		if math.Mod(float64(trimmedLength), 3) == 0 {
			var total, err = strconv.ParseFloat(item.Price, 64)
			if err != nil {
				continue
			}
			final += int(math.Ceil(total * 0.2))
		}
	}

	return final
}

// 6 points if the day in the purchase date is odd.
func oddDate(receipt_map Receipt) int {
	purchaseDate, _ := time.Parse("2006-01-02", receipt_map.PurchaseDate)
	if math.Mod(float64(purchaseDate.Day()), 2) != 0 {
		return 6
	}

	return 0
}

// 10 points if the time of purchase is after 2:00pm and before 4:00pm.
func afterTime(receipt_map Receipt) int {
	purchaseTime, _ := time.Parse("15:04", receipt_map.PurchaseTime)
	if purchaseTime.After(time.Date(0, time.January, 1, 14, 0, 0, 0, time.UTC)) && purchaseTime.Before(time.Date(0, time.January, 1, 16, 0, 0, 0, time.UTC)) {
		return 10
	}

	return 0
}

// main function
func main() {

	// define map into local memory
	receipt_map = make(map[string]Receipt)

	// create the api links
	router := gin.Default()

	router.POST("/receipts/process", postProcessReceipts)

	router.GET("/receipts/:id/points", getPoints)

	// start the server
	router.Run("localhost:9091")

}
