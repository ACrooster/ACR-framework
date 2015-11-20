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

    school := "amstelveencollege"
    token := framework.GetToken(school, code)

    fmt.Println(token)
}
