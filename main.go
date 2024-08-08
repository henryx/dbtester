package main

import (
	"dbtest/common"
	"dbtest/dbes"
	"dbtest/dbmongo"
	"dbtest/dbmysql"
	"dbtest/dbpg"
	"dbtest/dbsqlite"
	"log"
	"os"
	"time"

	"github.com/alecthomas/kong"
)

type Test interface {
	Name() string
	New(cli *common.CLI)
	Close()
	Load(size int, filename string)
	CountJSON() int64
	IndexJSON()
	FindJSON() int64
	Url() string
}

func test(dbtype string, c *common.CLI) {
	var test Test
	var start, end time.Time
	var duration time.Duration

	switch dbtype {
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
	case "sqlite":
		test = &dbsqlite.SQLite{}
	default:
		panic("Database not supported: " + dbtype)
	}

	test.New(c)
	defer test.Close()

	rows := c.Rows
	datafile := c.Datafile

	log.Println("Start load data on", test.Name(), "database (host", test.Url()+")")
	start = time.Now()
	test.Load(rows, datafile)
	end = time.Now()
	duration = end.Sub(start)
	log.Println("Finish load after", duration)

	time.Sleep(5 * time.Second)
	log.Println()

	log.Println("------------------")
	log.Println("Start count without index")
	start = time.Now()
	n := test.CountJSON()
	end = time.Now()
	duration = end.Sub(start)
	log.Printf("Counted %d items in %s", n, duration)

	time.Sleep(5 * time.Second)
	log.Println()

	log.Println("------------------")
	log.Println("Start find without index")
	start = time.Now()
	n = test.FindJSON()
	end = time.Now()
	duration = end.Sub(start)
	log.Printf("Found %d items without index in %s", n, duration)

	time.Sleep(5 * time.Second)
	log.Println()

	log.Println("------------------")
	log.Println("Start index field")
	start = time.Now()
	test.IndexJSON()
	end = time.Now()
	duration = end.Sub(start)
	log.Printf("Index field finished in %s", duration)

	time.Sleep(5 * time.Second)
	log.Println()

	log.Println("------------------")
	log.Println("Start find with index")
	start = time.Now()
	n = test.FindJSON()
	end = time.Now()
	duration = end.Sub(start)
	log.Printf("Found %d items using index in %s", n, duration)
}

func main() {
	var c common.CLI
	var err error

	parser, err := kong.New(&c)
	if err != nil {
		panic(err)
	}

	ctx, err := parser.Parse(os.Args[1:])
	if err != nil {
		panic(err)
	}

	test(ctx.Command(), &c)
}
