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
    STATUS_FREE
    STATUS_DATE
)

var scheduleData []*gabs.Container
var classCount float64

func RequestScheduleData() {

    user := "~me"
    start := "1448060400"
    end :=   "1448665200"

    // Execute the get request
    res, err := http.Get("https://" + "amstelveencollege" + ".zportal.nl/api/v2/appointments?user=" + user + "&start=" + start + "&end=" + end + "&access_token=ucrer3dmolfjsl846lt58pji56")

    // Check if an error has occurec
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

func GetClassCount() int {

    return int(classCount)
}

func IsClassValid(index int) bool {

    return scheduleData[index].Path("valid").Data().(bool)
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

    if scheduleData[index].Path("type").Data().(string) == "activity" {

	status = STATUS_ACTIVITY
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

    return strconv.Itoa(hr) + ":" + strconv.Itoa(min)
}

