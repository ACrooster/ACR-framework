package framework

import (
    "encoding/json"
    "fmt"
    "strings"
    "net/http"
    "net/url"
)

var school string

func SetSchool(m_school string) {

    school = m_school
}

func GetToken(m_code string) string {

    ErrorCode = ERROR_NONE

    // Variable that stores post data
    values := url.Values{}
    // Set the type of autorisation requested
    values.Set("grant_type", "authorization_code")
    // Set the verification code
    values.Set("code", m_code)

    // Execute the post request
    // TODO: Make this not harcoded
    res, err:= http.PostForm("https://" + school + ".zportal.nl/api/v2/oauth/token", values)

    token := ""

    // Check if an error has occurec
    if err == nil {

	// Cleanup
	defer res.Body.Close()

	var v map[string]interface{}
	err := json.NewDecoder(res.Body).Decode(&v)

	if err == nil {

	    // Decode the json response and store it in an array
	    switch r := v["access_token"].(type) {

	    case string:
		token = r
	    }
	} else {

	    fmt.Print(err)

	    ErrorCode = ERROR_UNKNOWN
	}
    } else {

	if strings.Contains(err.Error(), "No address associated with hostname") {

	    ErrorCode = ERROR_CONNECTION
	} else {

	    ErrorCode = ERROR_UNKNOWN
	}
    }

    return token
}
