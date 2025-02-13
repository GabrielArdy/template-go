package commons

import "time"

var location *time.Location

func GenerateCustomUID() string {
	currentTime := time.Now()
	return currentTime.Format("20060102150405")
}

func GetLocalTime() time.Time {
	location, _ := time.LoadLocation("Asia/Makassar")
	return time.Now().In(location)
}
