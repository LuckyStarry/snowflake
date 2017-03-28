package snowflake

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
		var sq, ms int64
		var ts = time.Now().UnixNano() / 1000 / 1000
		for true {
			wid := <-worker
			ms = time.Now().UnixNano() / 1000 / 1000
			if ms > ts {
				ts = ms
				sq = 0
			}
			sq++
			sequence <- &snowflakeID{ms, wid, sq}
		}
	}()
}

var worker = make(chan int64)
var sequence = make(chan *snowflakeID)

type snowflakeID struct {
	timestamp int64
	workerID  int64
	sequence  int64
}

// NextID used to gets an unique id with worker id : 0
func NextID() int64 {
	return NextIDWorker(0)
}

// NextIDWorker used to gets an unique id with explicit worker id
func NextIDWorker(workerID int64) int64 {
	if workerID >= 0 && workerID <= workerIDMax {
		worker <- workerID
		return (<-sequence).ToInt64()
	}
	return -1
}

func (id *snowflakeID) ToInt64() int64 {
	return (id.timestamp-twepoch)<<tsShift | id.workerID<<wkShift | id.sequence%sequenceMax
}
