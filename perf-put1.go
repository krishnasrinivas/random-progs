package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"

	"time"

	"github.com/minio/minio-go"
)

func usage() {
	fmt.Println("perf <object-size-in-MB> <parallel-upload-count>")
	os.Exit(0)
}

func main() {
	bucket := "testbucket"
	objectPrefix := "testobject"
	if len(os.Args) != 3 {
		usage()
	}
	length, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Print(err)
		usage()
	}
	length = length * 1024 * 1024
	nr, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Print(err)
		usage()
	}
	f, err := os.Open("bigfile")
	if err != nil {
		log.Fatal(err)
	}
	ch := make(chan struct{})
	var wg = &sync.WaitGroup{}
	for i := 0; i < nr; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			client, err := minio.NewCore(os.Getenv("MINIO_ENDPOINT"), os.Getenv("MINIO_ACCESS_KEY"), os.Getenv("MINIO_SECRET_KEY"), false)
			if err != nil {
				log.Fatal(err)
			}
			// Start all the goroutines at the same time
			<-ch
			fmt.Println("starting", i)
			_, err = client.PutObject(bucket, fmt.Sprintf("%s.%d", objectPrefix, i), int64(length), io.NewSectionReader(f, 0, int64(length)), nil, nil, nil)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("ending", i)
		}(i)
	}
	t1 := time.Now()
	close(ch)
	wg.Wait() // Wait till all go routines finish
	delta := time.Since(t1).Seconds()
	bandwidth := float64(length*nr) / delta / 1024 / 1024 // in MBps
	objPerSec := float64(nr) / delta
	fmt.Printf("data=%d bytes, time=%f seconds, obj/sec=%f, bandwidth=%f MBps\n", length*nr, delta, objPerSec, bandwidth)
}
