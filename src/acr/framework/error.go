package framework

import (
    "fmt"
)

var errorCode int
var errorMsg string

const (

    ERROR_NONE = iota
    ERROR_SCHOOL
    ERROR_CONNECTION
    ERROR_CODE
    ERROR_UNKNOWN
)

func GetError() int {

    // TODO: Figure out a beter way to print this
    fmt.Println(errorMsg, ": ErrorCode:", errorCode)
    tErrorCode := errorCode

    errorCode = ERROR_NONE
    errorMsg = ""

    return tErrorCode
}

func setError(mErrorCode int, mErrorMsg string) {

    errorCode = mErrorCode
    errorMsg = mErrorMsg
}
