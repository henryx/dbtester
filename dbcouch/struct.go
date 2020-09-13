package dbcouch

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"gopkg.in/ini.v1"
)

type CouchDB struct {
	host     string
	port     int
	database string
	user     string
	password string
}

func (c *CouchDB) call(method, uri string, buffer io.Reader) *http.Response {
	client := &http.Client{}

	req, err := http.NewRequest(method, uri, buffer)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	return resp
}

func (c *CouchDB) clean() {
	c.call("DELETE", fmt.Sprintf("http://%s:%s@%s:%d/%s", c.user, c.password, c.host, c.port, c.database), bytes.NewBuffer([]byte(nil)))
}

func (c *CouchDB) create() {
	c.call("PUT", fmt.Sprintf("http://%s:%s@%s:%d/%s", c.user, c.password, c.host, c.port, c.database), bytes.NewBuffer([]byte(nil)))
}

func (c *CouchDB) New(cfg *ini.Section) {
	c.host = cfg.Key("host").MustString("localhost")
	c.port = cfg.Key("port").MustInt(5984)
	c.user = cfg.Key("user").MustString("admin")
	c.password = cfg.Key("password").MustString("password")
	c.database = cfg.Key("database").MustString("test")

	c.clean()
	c.create()
}

func (c *CouchDB) Close() {
	return
}

func (c *CouchDB) Name() string {
	return "CouchDB"
}

func (c *CouchDB) Url() string {
	return c.host
}
