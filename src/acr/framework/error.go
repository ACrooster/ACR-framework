package framework

var ErrorCode int

const ERROR_NONE = 0
const ERROR_CONNECTION = 1
const ERROR_UNKNOWN = 256

func GetError() int {

    return ErrorCode
}

