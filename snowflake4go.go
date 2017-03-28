package snowflake4go

import "time"

const (
	/* 2015-01-01 */
	twepoch     int64 = 1420041600000
	wkShift     uint  = 12
	tsShift     uint  = 12 + 10
	workerIDMax       = 0x3ff
	sequenceMax       = 0xfff
)

func init() {
	go func() {
		var sq int64
		var ms int64
		var ts int64 = -1
		for true {
			ms = <-milliseconds
			if ms > ts {
				ts = ms
				sq = 0
			}
			sq++
			sequence <- sq
		}
	}()
}

var milliseconds = make(chan int64)
var sequence = make(chan int64)

// NextID used to gets an unique id with worker id : 0
func NextID() int64 {
	return NextIDWorker(0)
}

// NextIDWorker used to gets an unique id with explicit worker id
func NextIDWorker(workerID int64) int64 {
	if workerID >= 0 && workerID <= workerIDMax {
		now := time.Now().UnixNano() / 1000 / 1000
		milliseconds <- now
		return <-sequence%sequenceMax | (now-twepoch)<<tsShift | workerID<<wkShift
	}
	return -1
}
