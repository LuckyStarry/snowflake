# snowflake4go [![Build Status](https://travis-ci.org/LuckyStarry/snowflake4go.svg)](https://travis-ci.org/LuckyStarry/snowflake4go)

## About snowflake4go
snowflake4go is an implementation in golang without lock.

## How to use
```go
uid := snowflake4go.NextID()
```

or you have an explicit worker id (machine id) like this:
```go
uid := snowflake4go.NextIDWorker(workerId)
```

## Copyright and license
Code and documentation copyright 2017 Sun Bo. Code released under [the MIT license](https://github.com/LuckyStarry/snowflake4go/blob/master/LICENSE).