package framework

import (
    "github.com/jeffail/gabs"
    "strings"
    "net/http"
    "io/ioutil"
)

var userData []*gabs.Container

func RequestUserData() {

    // Execute the get request
    res, err := http.Get("https://" + school + ".zportal.nl/api/v3/users/~me?access_token=" + access_token)

    // Check if an error has occurec
    if err == nil {

	resByte, _ := ioutil.ReadAll(res.Body)

	// Cleanup
	defer res.Body.Close()

	jsonParsed, err := gabs.ParseJSON(resByte)

	userData, _ = jsonParsed.Path("response.data").Children()

	// TODO: Do more error checking
	if err != nil {

	    setError(ERROR_UNKNOWN, err.Error())
	}

    } else {

	if strings.Contains(err.Error(), "No address associated with hostname") {

	    setError(ERROR_CONNECTION, err.Error())
	} else {

	    setError(ERROR_UNKNOWN, err.Error())
	}
    }
}

func GetId() string {

    return userData[0].Path("code").Data().(string)
}

func GetName() string {

    p := userData[0].Path("prefix").Data().(string)
    if p == "" {

	return userData[0].Path("firstName").Data().(string) + " " + userData[0].Path("lastName").Data().(string)
    } else {

	return userData[0].Path("firstName").Data().(string) + " " + p + " " + userData[0].Path("lastName").Data().(string)
    }
}

func IsEmployee() bool {

    return userData[0].Path("isEmployee").Data().(bool)
}
