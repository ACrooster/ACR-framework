package framework

import (
    "github.com/jeffail/gabs"
    "strings"
    "net/http"
    "io/ioutil"
    "time"
    "strconv"
)

const (
    STATUS_NORMAL = iota
    STATUS_CHANGED
    STATUS_CANCELLED
    STATUS_ACTIVITY
    STATUS_EXAM
    STATUS_FREE
    STATUS_DATE
)

var scheduleData []*gabs.Container
var classCount float64
var year int
var week int

func RequestScheduleData(weekUnix int64) {

    year, week = time.Unix(weekUnix, 0).ISOWeek()

    user := "~me"
    start := FirstDayOfISOWeek()
    end := start + 604800

    // Execute the get request
    res, err := http.Get("https://" + "amstelveencollege" + ".zportal.nl/api/v2/appointments?user=" + user + "&start=" + strconv.Itoa(start) + "&end=" + strconv.Itoa(end) + "&valid=true&access_token=" + access_token)

    // Check if an error has occured
    if err == nil {

	resByte, _ := ioutil.ReadAll(res.Body)
	// fmt.Println("DATA:")
	// fmt.Println(string(resByte))

	// Cleanup
	defer res.Body.Close()

	jsonParsed, err := gabs.ParseJSON(resByte)

	scheduleData, _ = jsonParsed.Path("response.data").Children()
	classCount, _ = jsonParsed.Path("response.totalRows").Data().(float64)

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

func GetWeek() int {

    return week
}

func GetClassCount() int {

    return int(classCount)
}

func GetClassStartUnix(index int) int64 {

    return int64(scheduleData[index].Path("start").Data().(float64))
}

func GetClassTimeSlot(index int) int {

    data := scheduleData[index].Path("startTimeSlot").Data()

    if data != nil {

	return int(data.(float64))
    } else {

	return 0
    }
}

// TODO: This function is propably really crappy
func IsClassValid(index int) bool {

    bStart := scheduleData[index].Path("start").Data().(float64)
    bEnd := scheduleData[index].Path("end").Data().(float64)

    for i := 0; i < GetClassCount(); i++ {
	if (i != index) {
	    iStart := scheduleData[i].Path("start").Data().(float64)
	    iEnd := scheduleData[i].Path("end").Data().(float64)

	    if (bStart >= iStart && bStart < iEnd && GetClassStatus(index) == STATUS_CANCELLED) {

		return false
	    }

	    if (bEnd > iStart && bEnd <= iEnd && GetClassStatus(index) == STATUS_CANCELLED) {

		return false
	    }
	}
    }

    return true
}

func GetClassName(index int) string {

    if index < int(classCount) {

	tmp, _ := scheduleData[index].Path("subjects").Children()
	return tmp[0].Data().(string)
    } else {

	// TODO: This should throw an error
	return ""
    }
}

// TODO: These two functions could be more generic
func GetClassTeacher(index int) string {

	tmp, _ := scheduleData[index].Path("teachers").Children()
	var tmp2 string
	for i := 0; i < len(tmp); i++ {
	    tmp2 += strings.ToUpper(tmp[i].Data().(string))

	    if i < len(tmp)-1 {
		tmp2 += " & "
	    }
	}

	return tmp2
}

func GetClassRoom(index int) string {

	tmp, _ := scheduleData[index].Path("locations").Children()
	var tmp2 string
	for i := 0; i < len(tmp); i++ {
	    tmp2 += strings.ToUpper(tmp[i].Data().(string))

	    if i < len(tmp)-1 {
		tmp2 += " & "
	    }
	}

	return tmp2
}

func GetClassStatus(index int) int {

    var status int
    status = STATUS_NORMAL

    if scheduleData[index].Path("type").Data().(string) != "lesson" {

	status = STATUS_ACTIVITY
    }

    if scheduleData[index].Path("type").Data().(string) == "exam" {

	status = STATUS_EXAM
    }

    if scheduleData[index].Path("modified").Data().(bool) || scheduleData[index].Path("moved").Data().(bool) {

	status = STATUS_CHANGED
    }

    if scheduleData[index].Path("cancelled").Data().(bool) {

	status = STATUS_CANCELLED
    }

    return status
}

func GetClassStartTime(index int) string {

    unixTimeStamp, _ := scheduleData[index].Path("start").Data().(float64)
    return formatTime(unixTimeStamp)
}

func GetClassEndTime(index int) string {

    unixTimeStamp, _ := scheduleData[index].Path("end").Data().(float64)
    return formatTime(unixTimeStamp)
}

func formatTime(unixTimeStamp float64) string {

    unixIntValue := int64(unixTimeStamp)
    timeStamp := time.Unix(unixIntValue, 0)
    timeStampUTC := timeStamp.UTC()
    hr, min, _ := timeStampUTC.Clock()

    hr += timeDiff

    if (hr > 23) {
	hr -= 24
    }

    if (hr < 0) {
	hr += 24
    }

    if (min >= 10) {

	return strconv.Itoa(hr) + ":" + strconv.Itoa(min)
    } else {

	return strconv.Itoa(hr) + ":0" + strconv.Itoa(min)
    }
}

func FirstDayOfISOWeek() int {
    date := time.Date(year, 0, 0, 0, 0, 0, 0, time.UTC)
    isoYear, isoWeek := date.ISOWeek()

    // iterate back to Monday
    for date.Weekday() != time.Monday {
	date = date.AddDate(0, 0, -1)
	isoYear, isoWeek = date.ISOWeek()
    }

    // iterate forward to the first day of the first week
    for isoYear < year {
	date = date.AddDate(0, 0, 7)
	isoYear, isoWeek = date.ISOWeek()
    }

    // iterate forward to the first day of the given week
    for isoWeek < week {
	date = date.AddDate(0, 0, 7)
	isoYear, isoWeek = date.ISOWeek()
    }

    return int(date.Unix())
}

func GetDayUnix(index int) int64 {

    return int64(FirstDayOfISOWeek() + 3600*24*index)
}

func GetDayNumber(index int) int {

    _, _, day := time.Unix(GetDayUnix(index), 0).Date()
    return day
}

func GetDayMonth(index int) int {

    _, month, _ := time.Unix(GetDayUnix(index), 0).Date()
    return int(month)
}
