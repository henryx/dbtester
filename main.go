package main

import (
	"dbtest/dbcouch"
	"dbtest/dbes"
	"dbtest/dbmongo"
	"dbtest/dbmysql"
	"dbtest/dbpg"
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/ini.v1"
)

type Test interface {
	Name() string
	New(cfg *ini.Section)
	Close()
	Load(size int, filename string)
	Count() int64
	Index()
	Find() int64
	Url() string
}

func test(cfg *ini.File) {
	var test Test
	var start, end time.Time
	var duration time.Duration

	switch cfg.Section("default").Key("database").String() {
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
	case "couchdb":
		test = &dbcouch.CouchDB{}
	default:
		panic("Database not supported: " + cfg.Section("default").Key("database").String())
	}

	test.New(cfg.Section(cfg.Section("default").Key("database").String()))
	defer test.Close()

	rows := cfg.Section("default").Key("rows").MustInt(1)
	datafile := cfg.Section("default").Key("datafile").String()

	log.Println("Started load data on", test.Name(), "database (host", test.Url()+")")
	start = time.Now()
	test.Load(rows, datafile)
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
	var inifile string

	if len(os.Args) <= 1 {
		inifile = "load.ini"
	} else {
		inifile = os.Args[1]
	}

	cfg, err := ini.Load(inifile)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	test(cfg)
}
