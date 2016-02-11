package main

import (
    "acr/framework"
    "github.com/gopherjs/gopherjs/js"
    "honnef.co/go/js/dom"
    "github.com/fabioberger/cookie"
    "sort"
    "strconv"
    "time"
    "fmt"
)

var days = []string {
    "Maandag",
    "Dinsdag",
    "Woensdag",
    "Donderdag",
    "Vrijdag",
}
var months = []string {
    "januari",
    "februari",
    "maart",
    "april",
    "mei",
    "juni",
    "juli",
    "augustus",
    "september",
    "oktober",
    "november",
    "december",
}

var token string
const school = "amstelveencollege"

var user = framework.MY_SCHEDULE

func main() {

    framework.SetSchool(school)
    framework.SetTimeDiff(1)

    d := dom.GetWindow().Document()
    // auth := d.GetElementByID("auth_form").(*dom.HTMLFormElement)
    // auth.SetClass("")

    authButton := d.GetElementByID("auth_button")
    authButton.AddEventListener("click", false, func(event dom.Event) {

	go func() {

	    token = framework.GetToken(d.GetElementByID("auth").(*dom.HTMLInputElement).Value)
	    fmt.Println("Received token: " + token)
	    expires := time.Now().Add(time.Minute * 60 * 24 * 365)
	    cookie.Set("token", token, &expires, "/")
	    js.Global.Get("location").Call("reload", false)
	}()
    })

    token2, ok := cookie.Get("token")
    token = token2

    if ok {

	// fmt.Println("Token: " + token)
	framework.SetToken(token)
        //
	// code := d.GetElementByID("code_form").(*dom.HTMLFormElement)
	// code.SetClass("")
        //
	// codeButton := d.GetElementByID("code_button")
	// codeButton.AddEventListener("click", false, func(event dom.Event) {
        //
	//     go func() {
        //
	// 	user = d.GetElementByID("code").(*dom.HTMLInputElement).Value
	// 	fmt.Println(user)
        //
	// 	showSchedule()
	// 	fmt.Println(token)
	//     }()
	// })
	showSchedule()
    }


}

// TODO: Add error checking
func showSchedule() {

    fmt.Println(token)

    framework.RequestScheduleData(time.Now().Unix(), user)
    // framework.RequestScheduleData(1454948725, strconv.Itoa(108890))

    count := framework.GetClassCount()
    var classes = Classes{}

    for i := 0; i < count; i++ {

	class := classInfo{}

	class.name = framework.GetClassName(i)
	class.start = framework.GetClassStartTime(i)
	class.end = framework.GetClassEndTime(i)
	class.teacher = framework.GetClassTeacher(i)
	class.status = framework.GetClassStatus(i)
	class.unixStart = framework.GetClassStartUnix(i)
	class.timeSlot = framework.GetClassTimeSlot(i)

	classes = append(classes, class)
    }

    sort.Sort(classes)

    endOfDay := false
    for i := 0; i < count; i++ {

	free := 0

	if (i == 0 || endOfDay) {

	    free = classes[i].timeSlot - 1;

	    for j := 0; j < free; j++ {

		empty := classInfo{};
		empty.status = framework.STATUS_FREE
		empty.unixStart = classes[i].unixStart - int64(j + 1)
		empty.timeSlot = classes[i].timeSlot - j - 1

		classes = append(classes,empty)
	    }
	}

	if (i+1 < count) {

	    free = classes[i + 1].timeSlot - classes[i].timeSlot - 1
	    endOfDay = (classes[i+1].unixStart - classes[i].unixStart) > 10 * 3600
	} else {

	    free = 0
	    endOfDay = false
	}

	for j := 0; j < free; j++ {

		empty := classInfo{};
		empty.status = framework.STATUS_FREE
		empty.unixStart = classes[i].unixStart + int64(j + 1)
		empty.timeSlot = classes[i].timeSlot + j + 1

		classes = append(classes,empty)
	}
    }

    for i := 0; i < 5; i++ {

	day := classInfo{}

	day.name = days[i] + " " + strconv.Itoa(framework.GetDayNumber(i)) + " " + months[framework.GetDayMonth(i)-1]
	day.status = framework.STATUS_DATE
	day.unixStart = int64(framework.FirstDayOfISOWeek() + i * 24 * 3600)

	classes = append(classes, day)
    }

    sort.Sort(classes)

    // Actually output to the page
    d := dom.GetWindow().Document()
    schedule := d.GetElementByID("schedule").(*dom.HTMLDivElement)
    schedule.SetClass("mdl-grid")

    day := d.CreateElement("div").(*dom.HTMLDivElement)

    for _, v := range classes {

	if v.status == framework.STATUS_DATE {

	    schedule.AppendChild(day)
	    day = d.CreateElement("div").(*dom.HTMLDivElement)
	    // day.SetClass("mdl-cell")
	}

	container := d.CreateElement("div").(*dom.HTMLDivElement)

	if v.status == framework.STATUS_DATE {

	    date := d.CreateElement("h4").(*dom.HTMLHeadingElement)
	    date.SetTextContent(v.name)

	    container.AppendChild(date)
	} else {

	    class := d.CreateElement("div").(*dom.HTMLDivElement)
	    class.SetClass("class-card mdl-card mdl-shadow--4dp " + getColor(v.status))

	    information := d.CreateElement("div").(*dom.HTMLDivElement)
	    information.SetClass("mdl-card__supporting-text")

	    name := d.CreateElement("div").(*dom.HTMLDivElement)
	    name.SetClass("class-name")
	    name.SetTextContent(v.name)
	    information.AppendChild(name)

	    start := d.CreateElement("div").(*dom.HTMLDivElement)
	    start.SetClass("class-start")
	    start.SetTextContent(v.start)
	    information.AppendChild(start)

	    end := d.CreateElement("div").(*dom.HTMLDivElement)
	    end.SetClass("class-end")
	    end.SetTextContent(v.end)
	    information.AppendChild(end)

	    teacher := d.CreateElement("div").(*dom.HTMLDivElement)
	    teacher.SetClass("class-teacher")
	    teacher.SetTextContent(v.teacher)
	    information.AppendChild(teacher)

	    class.AppendChild(information)
	    container.AppendChild(class)
	}
	day.AppendChild(container)
    }
    schedule.AppendChild(day)
}

func getColor(status int) string {
    if status == framework.STATUS_CHANGED {
	return "class-changed"
    } else if status == framework.STATUS_CANCELLED {
	return "class-cancelled"
    } else if status == framework.STATUS_FREE {
	return "class-free"
    } else if status == framework.STATUS_ACTIVITY {
	return "class-activity"
    } else if status == framework.STATUS_EXAM {
	return "class-exam"
    }

    return ""
}

type classInfo struct {

    name string
    start string
    end string
    teacher string
    timeSlot int
    status int
    unixStart int64
}

type Classes []classInfo

func (slice Classes) Len() int {

    return len(slice)
}

func (slice Classes) Less(i, j int) bool {
    return slice[i].unixStart < slice[j].unixStart
}

func (slice Classes) Swap(i, j int) {
    slice[i], slice[j] = slice[j], slice[i]
}
