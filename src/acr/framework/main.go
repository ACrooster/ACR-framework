package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
)

const api string = "https://amstelveencollege.zportal.nl/api/v2/oauth/token"

func main() {

    // Ask user for verification code
    fmt.Print("Enter code: ")
    var code string
    fmt.Scanln(&code)

    // Variable that stores post data
    values := url.Values{}
    // Set the type of autorisation requested
    values.Set("grant_type", "authorization_code")
    // Set the verification code
    values.Set("code", code)

    // Execute the post request
    // TODO: Make this not harcoded
    res, err:= http.PostForm(api, values)

    // Check if an error has occurec
    if err != nil {

	// Print the error
	fmt.Println(err)
    } else {

	// Cleanup
	defer res.Body.Close()

	// Decode the json response and store it in an array
	var v map[string]interface{}
	err := json.NewDecoder(res.Body).Decode(&v)

	// Check for errors
	if err != nil {

	    // Print the error
	    fmt.Println(err)
	}

	// Print the received acces_token
	fmt.Println(v["access_token"])
    }


}
