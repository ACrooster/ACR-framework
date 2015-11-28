package main

import (
    "acr/framework"
    "fmt"
)

func main() {

    // // Ask user for verification code
    // fmt.Print("Enter code: ")
    // var code string
    // // NOTE: This does not work with gopherjs
    // fmt.Scanln(&code)
    //
    // framework.SetSchool("amstelveencollege")
    //
    // token := framework.GetToken(code)
    // framework.GetError()
    //
    // fmt.Println(token)

    framework.SetSchool("amstelveencollege")
    framework.SetToken("ucrer3dmolfjsl846lt58pji56")
    framework.RequestUserData()
    fmt.Println(framework.GetName())
    fmt.Println(framework.GetId())
    framework.GetError()
}
