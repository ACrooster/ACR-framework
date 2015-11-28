package framework

import (
    "encoding/json"
    "strings"
    "net/http"
    "io/ioutil"
)

var jsontype jsonobject

func RequestUser() {

    // Execute the get request
    res, err := http.Get("https://" + school + ".zportal.nl/api/v2/users/~me?access_token=" + access_token)

    // Check if an error has occurec
    if err == nil {

	resByte, _ := ioutil.ReadAll(res.Body)

	// Cleanup
	defer res.Body.Close()

	err := json.Unmarshal(resByte, &jsontype)

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

    return jsontype.Response.Data[0].Code
}

func GetName() string {

    p := jsontype.Response.Data[0].Prefix
    if p == "" {

	return jsontype.Response.Data[0].FirstName + " " + jsontype.Response.Data[0].LastName
    } else {

	return jsontype.Response.Data[0].FirstName + " " + jsontype.Response.Data[0].Prefix + " " + jsontype.Response.Data[0].LastName
    }
}

type jsonobject struct {
    Response ResponseType
}

type ResponseType struct {
    Status int
    Message string
    Details string
    EventId int
    StartRow int
    EndRow int
    TotalRows int
    Data []DataType
}

type DataType struct {
    Code string
    Roles []string
    IsStudent bool
    IsEmployee bool
    IsFamilyMember bool
    FirstName string
    Prefix string
    LastName string
    Archived bool
}
