package main

import (
	"fmt"
	"load/dbes"
	"load/dbmongo"
	"load/dbmysql"
	"load/dbpg"
	"log"
	"os"
	"strconv"
	"time"
)

type Test interface {
	Name() string
	New(host string)
	Close()
	Load(size int, filename string)
	Count() int64
	Index()
	Find() int64
}

func test(host, db string, size int, filename string) {
	var test Test
	var start, end time.Time
	var duration time.Duration

	switch db {
	case "mongo":
		test = &dbmongo.Mongo{}
		break
	case "postgres":
		test = &dbpg.Postgres{}
		break
	case "mysql":
		test = &dbmysql.MySQL{}
	case "elasticsearch":
		test = &dbes.Elasticsearch{}
	default:
		panic("Database not supported: " + db)
	}

	test.New(host)
	defer test.Close()

	log.Println("Started load data on", test.Name(), "database (host", host+")")
	start = time.Now()
	test.Load(size, filename)
	end = time.Now()
	duration = end.Sub(start)
	log.Println("Finished load after", duration)

	time.Sleep(5 * time.Second)
	log.Println()

	log.Println("------------------")
	log.Println("Start count without index")
	start = time.Now()
	c := test.Count()
	end = time.Now()
	duration = end.Sub(start)
	log.Printf("Counted %d items in %s", c, duration)

	time.Sleep(5 * time.Second)
	log.Println()

	log.Println("------------------")
	log.Println("Start find without index")
	start = time.Now()
	c = test.Find()
	end = time.Now()
	duration = end.Sub(start)
	log.Printf("Found %d items without index in %s", c, duration)

	time.Sleep(5 * time.Second)
	log.Println()

	log.Println("------------------")
	log.Println("Start index field")
	start = time.Now()
	test.Index()
	end = time.Now()
	duration = end.Sub(start)
	log.Printf("Indexed field in %s", duration)

	time.Sleep(5 * time.Second)
	log.Println()

	log.Println("------------------")
	log.Println("Start find with index")
	start = time.Now()
	c = test.Find()
	end = time.Now()
	duration = end.Sub(start)
	log.Printf("Found %d items using index in %s", c, duration)
}

func main() {

	if len(os.Args) <= 1 {
		fmt.Println("Usage:", os.Args[0], "<database> <host> <number of records before commit> <filename>")
		os.Exit(1)
	}

	dbtype := os.Args[1]
	host := os.Args[2]
	size, err := strconv.Atoi(os.Args[3])
	if err != nil {
		panic("Cannot parse argument: " + err.Error())
	}
	filename := os.Args[4]

	test(host, dbtype, size, filename)
}
