// rocksdb test
package main

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/satori/go.uuid"
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

func doPut(basedir string, from int, to int) <-chan putStat {
	var totalTime time.Duration
	var maxTime time.Duration
	putCount := 0
	loopCount := 100000
	putStatCh := make(chan putStat)

	createFile := func(key string) {
		value := randBytes(from, to)
		start := time.Now()
		filename := path.Join(basedir, key)
		F, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_SYNC, 0644)
		if err != nil {
			fmt.Println("OpenFile error : ", err)
			os.Exit(1)
		}
		_, err = F.Write(value)
		if err != nil {
			fmt.Println("Write error : ", err)
			os.Exit(1)
		}
		err = F.Close()
		if err != nil {
			fmt.Println("Close error : ", err)
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
			createFile(key)
		}
		for {
			key := uuid.NewV4().String()
			createFile(key)
			if putCount == 30000000 {
				close(putStatCh)
			}
		}
	}()
	return putStatCh
}

func main() {
	basedir := "/data/fstest"
	putStatch := doPut(basedir, 16*1024, 32*1024)
	for p := range putStatch {
		fmt.Println("avgTime ", p.avgTime, "maxTime", p.maxTime)
	}
	fmt.Println("Exiting. putCount == 30000000")
}
