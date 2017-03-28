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

var workers [workerIDMax]chan chan *snowflakeID

func init() {
	var workerID int64
	for ; workerID < workerIDMax; workerID++ {
		workers[workerID] = make(chan chan *snowflakeID)
		go process(workerID)
	}
}

func process(workerID int64) {
	worker := workers[workerID]
	var timestamp, lasttime, sequence int64
	for true {
		context := <-worker
		timestamp = time.Now().UnixNano() / 1000 / 1000
		if timestamp > lasttime {
			lasttime = timestamp
			sequence = 0
		}
		sequence++
		context <- &snowflakeID{timestamp, workerID, sequence - 1}
	}
}

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
		var context = make(chan *snowflakeID)
		workers[workerID] <- context
		return (<-context).ToInt64()
	}
	return -1
}

func (id *snowflakeID) ToInt64() int64 {
	return (id.timestamp-twepoch)<<tsShift | id.workerID<<wkShift | id.sequence%sequenceMax
}
