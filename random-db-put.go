// rocksdb test
package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

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

	createDBEntry := func(key string) {
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
	go func() {
		for i, key := range filenames {
			createDBEntry(key)
		}
		close(putStatCh)
	}()
	return putStatCh
}

func main() {
	dbName := "/data/mydb"
	cache := gorocksdb.NewLRUCache(512 * 1024 * 1024)
	filter := gorocksdb.NewBloomFilter(15)
	to := gorocksdb.NewDefaultBlockBasedTableOptions()
	to.SetBlockSize(256 * 1024)
	to.SetBlockCache(cache)
	to.SetFilterPolicy(filter)
	options := gorocksdb.NewDefaultOptions()
	options.SetBlockBasedTableFactory(to)
	options.SetCreateIfMissing(true)
	options.SetStatsDumpPeriodSec(60 * 1) // dump stats at 10 minute interval
	options.SetCompactionStyle(gorocksdb.UniversalCompactionStyle)
	options.SetWriteBufferSize(512 * 1024 * 1024)
	options.SetMaxWriteBufferNumber(5)
	options.SetMinWriteBufferNumberToMerge(2)
	db, err := gorocksdb.OpenDb(options, dbName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	putStatch := doPut(db, 16*1024, 32*1024)
	for {
		p := <-putStatch
		fmt.Println("avgTime ", p.avgTime, "maxTime", p.maxTime)
	}
}
