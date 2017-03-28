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
		var workers [workerIDMax]chan chan *snowflakeID
		for true {
			ctx := <-worker
			wc := workers[ctx.wid]
			if wc == nil {
				wc = make(chan chan *snowflakeID)
				workers[ctx.wid] = wc
				go func(n int64, c chan chan *snowflakeID) {
					var ms, ts, sq int64
					for true {
						sfc := <-c
						ms = time.Now().UnixNano()
						if ms > ts {
							ts = ms
							sq = 0
						}
						sq++
						sfc <- &snowflakeID{ms / 1000 / 1000, n, sq - 1}
					}
				}(ctx.wid, wc)
			}
			go func() {
				wc <- ctx.sfc
			}()
		}
	}()
}

var worker = make(chan *context)

type context struct {
	wid int64
	sfc chan *snowflakeID
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
		var ctx = &context{workerID, make(chan *snowflakeID)}
		worker <- ctx
		return (<-ctx.sfc).ToInt64()
	}
	return -1
}

func (id *snowflakeID) ToInt64() int64 {
	return (id.timestamp-twepoch)<<tsShift | id.workerID<<wkShift | id.sequence%sequenceMax
}
