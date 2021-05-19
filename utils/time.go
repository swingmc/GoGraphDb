package utils

import (
	"math/rand"
	"time"
	"GoGraphDb/conf"
)

func LogTime() string{
		return string(ParseInt64ToString(time.Now().UnixNano())) + conf.Splitor + time.Now().String()
}

func LogTimeByVersion(versionId int64) string{
		return string(ParseInt64ToString(versionId)) + conf.Splitor + time.Now().String()
}

func GenRetryTime(n int) float64{
	return float64(rand.Intn(n) + conf.BaseRetryTime)/1000
}