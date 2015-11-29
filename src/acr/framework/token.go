package framework

import (
    "encoding/json"
    "strings"
    "net/http"
    "net/url"
    "io/ioutil"
)

func IsStillValid() bool {

    // TODO: Check if the acces token is still valid
    return true
}


func IsValid(mCode string) bool {

    return len(mCode) == 12
}

func GetToken(mCode string) string {

    // Check if the school has been set
    if school == "" {
	setError(ERROR_SCHOOL, "School name not set")
	return ""
    }

    // Variable that stores post data
    values := url.Values{}
    // Set the type of autorisation requested
    values.Set("grant_type", "authorization_code")
    // Set the verification code
    values.Set("code", mCode)

    // Execute the post request
    res, err := http.PostForm("https://" + school + ".zportal.nl/api/v2/oauth/token", values)

    token := ""

    // Check if an error has occurec
    if err == nil {

	resByte, _ := ioutil.ReadAll(res.Body)

	// Cleanup
	defer res.Body.Close()

	var v map[string]interface{}
	err := json.Unmarshal(resByte, &v)

	if err == nil {

	    // Decode the json response and store it in an array
	    switch r := v["access_token"].(type) {

	    case string:
		token = r
	    }
	    // TODO: There should also be an error in case the json does not contain the correct item
	// I assume status 400 is only send on a wrong code, as all the other syntax should be correct
	} else if strings.Contains(string(resByte), "Status 400") {

	    setError(ERROR_CODE, "Invalid code")
	    return ""
	} else {

	    setError(ERROR_UNKNOWN, err.Error())
	    return ""
	}
    } else {

	if strings.Contains(err.Error(), "No address associated with hostname") {

	    setError(ERROR_CONNECTION, err.Error())
	    return ""
	} else {

	    setError(ERROR_UNKNOWN, err.Error())
	    return ""
	}
    }

    access_token = token
    return token
}
