// rocksdb test
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"
)

type getStat struct {
	maxTime time.Duration
	avgTime time.Duration
}

func doGet(basedir string) <-chan getStat {
	var totalTime time.Duration
	var maxTime time.Duration
	getCount := 0
	loopCount := 10000
	getStatCh := make(chan getStat)

	getFile := func(key string) {
		start := time.Now()
		filename := path.Join(basedir, key)
		_, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Println("ReadFile error : ", err)
			os.Exit(1)
		}
		delta := time.Since(start)
		if delta > maxTime {
			maxTime = delta
		}
		totalTime += delta
		getCount++
		if (getCount % loopCount) == 0 {
			avgTime := totalTime / time.Duration(loopCount)
			getStatCh <- getStat{maxTime, avgTime}
			maxTime = 0
			totalTime = 0
		}
	}

	go func() {
		for i, key := range filenames {
			getFile(key)
		}
	}()
	return getStatCh
}

func main() {
	basedir := "/data/fstest"
	putStatch := doGet(basedir)
	for p := range putStatch {
		fmt.Println("avgTime ", p.avgTime, "maxTime", p.maxTime)
	}
	fmt.Println("Done")
}
