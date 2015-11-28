package framework

import (
    "encoding/json"
    "strings"
    "net/http"
    "net/url"
    "io/ioutil"
)

var school string

func SetSchool(mSchool string) {

    school = mSchool
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
    // TODO: Make this not harcoded
    res, err := http.PostForm("https://" + school + ".zportal.nl/api/v2/oauth/token", values)

    token := ""

    resByte, _ := ioutil.ReadAll(res.Body)

    // Check if an error has occurec
    if err == nil {

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

    return token
}
