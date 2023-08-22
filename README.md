# Receipt_Process_Challenge

## Description

This project is part of an assessment to test candidate's coding skills. This project has two api endpoints and it is written in GO language.
The following rules are being followed:
* One point for every alphanumeric character in the retailer name.
* 50 points if the total is a round dollar amount with no cents.
* 25 points if the total is a multiple of `0.25`.
* 5 points for every two items on the receipt.
* If the trimmed length of the item description is a multiple of 3, multiply the price by `0.2` and round up to the nearest integer. The result is the number of points earned.
* 6 points if the day in the purchase date is odd.
* 10 points if the time of purchase is after 2:00pm and before 4:00pm.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Contributing](#contributing)

## Installation

- The “fmt” library needs to be installed which can be installed using "go install fmt" OR please refer to https://golang.org/pkg/fmt/ 
- The other libraries are part of the standard Go library and do not require installation.

## Usage

### This project is built to work on local system.
- The base URL is http://localhost:9091/
- Base URL can be changed by changing the string in the last line in the main function of main.go (Usually second last line of the main.go file)

### How to run the project.
- First make sure the fmt dependency is installed. It is used to print out errors if any occured.
- Execute the command 'go run main.go' in the project directory.

### Api endpoints
#### POST request to send the json data - Path: `/receipts/process`
- This api requires json data as following:

```json
{
  "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    },{
      "shortDescription": "Emils Cheese Pizza",
      "price": "12.25"
    },{
      "shortDescription": "Knorr Creamy Chicken",
      "price": "1.26"
    },{
      "shortDescription": "Doritos Nacho Cheese",
      "price": "3.35"
    },{
      "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
      "price": "12.00"
    }
  ],
  "total": "35.35"
}
```
- With valid data, this endpoint will return a UUID with the following structure of json data.
```json
{
    "id": "0059b9a5-2c83-4859-8846-e15d128e8923"
}
```
  

#### GET request to send the json data - Path: `/receipts/{id}/points`
- This endpoint is used to pass in the id and get the total points for that specific ID
- On passing a valid id points will be returned in json format as below
```json
{
    "points": 28
}
```



## Contributing

The baseline code was taken from https://github.com/fetch-rewards/receipt-processor-challenge.

