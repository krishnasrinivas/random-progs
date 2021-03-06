// rocksdb test
package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/satori/go.uuid"
	"github.com/tecbot/gorocksdb"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randBytes(from int, to int) []byte {
	n := from + rand.Intn(to-from)
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return (b)
}

type putStat struct {
	maxTime time.Duration
	avgTime time.Duration
}

func doPut(db *gorocksdb.DB, from int, to int) <-chan putStat {
	var totalTime time.Duration
	var maxTime time.Duration
	putCount := 0
	loopCount := 100000
	putStatCh := make(chan putStat)

	wo := gorocksdb.NewDefaultWriteOptions()

	go func() {
		for {
			key := []byte(uuid.NewV4().String())
			value := randBytes(from, to)
			start := time.Now()
			err := db.Put(wo, key, value)
			if err != nil {
				fmt.Println("dbPut error : ", err)
				os.Exit(1)
			}
			delta := time.Since(start)
			if delta > maxTime {
				maxTime = delta
			}
			totalTime += delta
			putCount++
			if (putCount % loopCount) == 0 {
				avgTime := totalTime / time.Duration(loopCount)
				putStatCh <- putStat{maxTime, avgTime}
				maxTime = 0
				totalTime = 0
			}
		}
	}()
	return putStatCh
}

func main() {
	dbName := "/data/mydb"
	compactionThreads := 4
	options := gorocksdb.NewDefaultOptions()
	options.SetCreateIfMissing(true)
	options.SetStatsDumpPeriodSec(60 * 10) // dump stats at 10 minute interval
	options.SetCompactionStyle(gorocksdb.UniversalCompactionStyle)
	options.IncreaseParallelism(compactionThreads)
	options.IncreaseParallelism(compactionThreads)
	db, err := gorocksdb.OpenDb(options, dbName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	putStatch := doPut(db, 64*1024, 128*1024)
	for {
		p := <-putStatch
		fmt.Println("avgTime ", p.avgTime, "maxTime", p.maxTime)
	}
}
