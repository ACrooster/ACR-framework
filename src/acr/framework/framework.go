package framework

import (
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
)

func GetToken(school string, code string) string {

    // Variable that stores post data
    values := url.Values{}
    // Set the type of autorisation requested
    values.Set("grant_type", "authorization_code")
    // Set the verification code
    values.Set("code", code)

    // Execute the post request
    // TODO: Make this not harcoded
    res, err:= http.PostForm("https://" + school + ".zportal.nl/api/v2/oauth/token", values)

    token := ""

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

	switch r := v["access_token"].(type) {

	    case string:
		token = r

	    default:
		fmt.Println("Something went wrong")
	}
    }

    return token
}
