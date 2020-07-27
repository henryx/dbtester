package dbmongo

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m *Mongo) Load(size int, filename string) {
	//var j map[string]interface{}
	var j interface{}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	session, err := m.client.StartSession()
	if err != nil {
		panic("Cannot open session " + filename)
	}
	defer session.EndSession(ctx)

	inFile, err := os.Open(filename)
	if err != nil {
		panic("Cannot open file " + filename)
	}
	defer inFile.Close()

	reader := bufio.NewReader(inFile)

	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		counter := 0
		commit := 0

		err = session.StartTransaction()
		if err != nil {
			panic("Error when start transaction: " + err.Error())
		}

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}

				if line == "" {
					log.Println("Line empty")
					continue
				}

				panic("Error when read line: " + err.Error())
			}
			/* err = json.Unmarshal([]byte(line), &j) */
			err = bson.UnmarshalExtJSON([]byte(line), true, &j)
			if err != nil {
				panic("Error when unmarshal data: " + err.Error())
			}

			ctx, _ := context.WithTimeout(context.Background(), 120*time.Second)
			_, err = m.collection.InsertOne(ctx, &j)
			if err != nil {
				panic(err.Error())
			}

			counter++
			if counter == size {
				counter = 0
				err = session.CommitTransaction(ctx)
				if err != nil {
					panic("Error when commit data: " + err.Error())
				}
				commit++
				log.Printf("Committed %d...\n", commit)

				err = session.StartTransaction()
				if err != nil {
					panic("Error when start transaction: " + err.Error())
				}
			}
		}
		return nil

	})
}
