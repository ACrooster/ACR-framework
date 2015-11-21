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

    framework.SetSchool("amstelveencollege")

    token := framework.GetToken(code)

    fmt.Println(token)
}
