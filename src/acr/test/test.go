package main

import (
    "acr/framework"
    "fmt"
)

func main() {

    // Ask user for verification code
    fmt.Print("Enter code: ")
    var code string
    // NOTE: This does not work with gopherjs
    fmt.Scanln(&code)

    data := new(framework.Data)
    data.School = "amstelveencollege"
    framework.GetToken(code, data)

    fmt.Println(data.School)
    fmt.Println(data.Token)
}
