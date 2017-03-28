# snowflake [![Build Status](https://travis-ci.org/LuckyStarry/snowflake.svg)](https://travis-ci.org/LuckyStarry/snowflake)

## About

This package is a snowflake implementation in Golang without lock.

## Getting Started

### Installing
``` go get github.com/LuckyStarry/snowflake ```

### Usage
```go
uid := snowflake.NextID()
```

or you have an explicit worker id (machine id) like this:
```go
uid := snowflake.NextIDWorker(workerId)
```

## Copyright and license
Code and documentation copyright 2017 Sun Bo. Code released under [the MIT license](https://github.com/LuckyStarry/snowflake/blob/master/LICENSE).