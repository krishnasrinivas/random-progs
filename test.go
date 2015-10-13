package main
// rocksdb test
import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/satori/go.uuid"
	"github.com/tecbot/gorocksdb"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes() []byte {
	n := 4*1024 + rand.Intn(28*1024)
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return (b)
}

func main() {
	dbName := "/data/mydb"
	options := gorocksdb.NewDefaultOptions()
	options.SetCreateIfMissing(true)
	options.SetStatsDumpPeriodSec(30)
	db, err := gorocksdb.OpenDb(options, dbName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(db.GetProperty("rocksdb.stats"))
	wo := gorocksdb.NewDefaultWriteOptions()
	i := 0
	start := time.Now()
	for {
		key := uuid.NewV4().String()
		value := RandStringBytes()

		err = db.Put(wo, []byte(key), value)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		i++
		fmt.Printf("count : %d\r", i)
		if (i % 500000) == 0 {
			elapsed := time.Since(start)
			fmt.Printf("time taken : %s\n", elapsed)
			start = time.Now()
		}
	}
}
