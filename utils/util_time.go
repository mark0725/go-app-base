package utils

import "time"

func GetCurDate8() string {
	curDate := time.Now().Format("20060102")
	return curDate
}

func GetCurDate10() string {
	curDate := time.Now().Format("2006-01-02")
	return curDate
}

func GetCurDateTime14() string {
	curTime := time.Now().Format("20060102150405")
	return curTime
}

func GetCurDateTime19() string {
	curTime := time.Now().Format("2006-01-02T15:04:05")
	return curTime
}

func GetCurTime6() string {
	curTime := time.Now().Format("150405")
	return curTime
}

func GetCurTime8() string {
	curTime := time.Now().Format("15:04:05")
	return curTime
}
