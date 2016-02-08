package main

import (
    "acr/framework"
    "strconv"
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
    framework.SetToken("rdnn2ng1c8t2ok3jd4ss3977v9")
    framework.SetTimeDiff(1)

    framework.RequestUserData()
    framework.GetError()
    fmt.Println(framework.GetName())
    fmt.Println(framework.GetId())

    framework.RequestScheduleData(1448918611)
    framework.GetError()
    classCount := framework.GetClassCount()
    fmt.Println(classCount)
    for i := 0; i < classCount; i++ {
	fmt.Println(framework.GetClassName(i) + " " + framework.GetClassStartTime(i) + " - " + framework.GetClassEndTime(i) + " " + framework.GetClassTeacher(i) + " " + framework.GetClassRoom(i) + " " + strconv.Itoa(framework.GetClassStatus(i)))

    }
}
