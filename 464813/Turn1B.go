package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Define a struct to represent the API response
type Response struct {
	Items []struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	} `json:"items"`
}

func main() {
	// Assuming you have a function to make an API request and retrieve the response body
	body := makeAPIRequest()

	// Create a variable to store the decoded API response
	var apiResponse Response

	// Decode the JSON response body into the struct
	err := json.Unmarshal(body, &apiResponse)
	if err != nil {
		panic(err)
	}

	// Use the range keyword to process the API response
	for _, item := range apiResponse.Items {
		fmt.Println("Item Name:", item.Name)
		fmt.Println("Item Age:", item.Age)
	}
}

func makeAPIRequest() []byte {
	// Replace this URL with the actual API endpoint you want to use
	url := "https://api.example.com/data"
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return body
}
