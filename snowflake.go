package snowflake

import (
	"bytes"
	"time"
)

const (
	// tw (Twitter) epoch is 2010-11-04 01:42:54.657
	twepoch     int64 = 1288834974657
	wkShift     uint  = 12
	tsShift     uint  = 12 + 10
	workerIDMax       = 0x3ff
	sequenceMax       = 0xfff

	alphabet62 = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
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

// ISnowflakeID is an entity of SnowflakeID
type ISnowflakeID interface {
	ToInt64() int64
	ToBase62() string
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
	sfid := NextSnowflakeID(workerID)
	if sfid != nil {
		return sfid.ToInt64()
	}
	return -1
}

// NextSnowflakeID used to gets an unique id with explicit worker id
func NextSnowflakeID(workerID int64) ISnowflakeID {
	if workerID >= 0 && workerID <= workerIDMax {
		var context = make(chan *snowflakeID)
		workers[workerID] <- context
		return <-context
	}
	return nil
}

func (id *snowflakeID) ToInt64() int64 {
	return (id.timestamp-twepoch)<<tsShift | id.workerID<<wkShift | id.sequence%sequenceMax
}

func (id *snowflakeID) ToBase62() string {
	buffer := &bytes.Buffer{}
	v := id.ToInt64()
	for ; v > 0; v /= 62 {
		buffer.WriteByte(alphabet62[v%62])
	}
	return buffer.String()
}
