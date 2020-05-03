package main

import (
	"os"

	client "github.com/Gimulator/client-go"
	"github.com/Gimulator/logger"
)

var (
	dir    string
	endKey string
	addr   string
)

func init() {
	dir = os.Getenv("LOGGER_DIR")
	if dir == "" {
		panic("set the 'LOGGER_DIR' environment variable for storing logs in it")
	}

	endKey = os.Getenv("LOGGER_END_KEY")
	if endKey == "" {
		panic("set the 'LOGGER_END_KEY' environment variable to exit the program")
	}

	addr = os.Getenv("LOGGER_ADDRESS")
	if addr == "" {
		panic("set the 'LOGGER_ADDRESS' environment variable to connect to the Gimulator, like localhost:3030")
	}
}

func main() {

	cli := client.NewClient(addr)
	err := cli.Register()
	if err != nil {
		panic(err)
	}

	ch := make(chan client.Object, 1024)
	err = cli.Socket(ch)
	if err != nil {
		panic(err)
	}

	err = cli.Watch(client.Key{
		Type:      "",
		Name:      "",
		Namespace: "",
	})

	recorder, err := logger.NewRecorder(ch, dir, endKey)
	if err != nil {
		panic(err)
	}
	recorder.Record()
}
