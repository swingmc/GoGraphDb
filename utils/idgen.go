package utils

import "time"

var (
	timeChan = make(chan int64, 100)
	timeBarrier = func() int64{
		time.Sleep(time.Nanosecond)
		return <-timeChan
	}

	idChan = make(chan int64, 100)
)

func init() {
	idChan <- 0
}

func GenId() int64{
	<-idChan
	id := time.Now().UnixNano()
	time.Sleep(time.Nanosecond)
	idChan <- 0
	return id
}

func GenTimeStamp() int64{
	timeChan <- time.Now().UnixNano()
	return timeBarrier()
}

