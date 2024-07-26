package main

import (
	"dbtest/cli"
	"dbtest/dbcouch"
	"dbtest/dbes"
	"dbtest/dbmongo"
	"dbtest/dbmysql"
	"dbtest/dbpg"
	"log"
	"time"

	"github.com/alecthomas/kong"
)

type Test interface {
	Name() string
	New(cli *cli.CLI)
	Close()
	Load(size int, filename string)
	CountJSON() int64
	IndexJSON()
	FindJSON() int64
	Url() string
}

func test(dbtype string, cli *cli.CLI) {
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
	case "couchdb":
		test = &dbcouch.CouchDB{}
	default:
		panic("Database not supported: " + dbtype)
	}

	test.New(cli)
	defer test.Close()

	rows := cli.Rows
	datafile := cli.Datafile

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
	c := test.CountJSON()
	end = time.Now()
	duration = end.Sub(start)
	log.Printf("Counted %d items in %s", c, duration)

	time.Sleep(5 * time.Second)
	log.Println()

	log.Println("------------------")
	log.Println("Start find without index")
	start = time.Now()
	c = test.FindJSON()
	end = time.Now()
	duration = end.Sub(start)
	log.Printf("Found %d items without index in %s", c, duration)

	time.Sleep(5 * time.Second)
	log.Println()

	log.Println("------------------")
	log.Println("Start index field")
	start = time.Now()
	test.IndexJSON()
	end = time.Now()
	duration = end.Sub(start)
	log.Printf("Indexed field in %s", duration)

	time.Sleep(5 * time.Second)
	log.Println()

	log.Println("------------------")
	log.Println("Start find with index")
	start = time.Now()
	c = test.FindJSON()
	end = time.Now()
	duration = end.Sub(start)
	log.Printf("Found %d items using index in %s", c, duration)
}

func main() {
	var cli cli.CLI

	ctx := kong.Parse(&cli)

	test(ctx.Command(), &cli)
}
