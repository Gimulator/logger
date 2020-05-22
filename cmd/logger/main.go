package main

import (
	"fmt"
	"os"

	client "github.com/Gimulator/client-go"
	"github.com/Gimulator/logger/concluder"
	"github.com/Gimulator/logger/recorder"
	"github.com/Gimulator/logger/uploader"
)

func main() {
	runner, err := newRunner()
	if err != nil {
		panic(err)
	}

	if err := runner.run(); err != nil {
		panic(err)
	}
}

type runner struct {
	addr     string
	ch       chan client.Object
	cli      *client.Client
	recorder *recorder.Recorder

	uploader  uploader.Uploader
	concluder concluder.Concluder
}

func newRunner() (*runner, error) {
	r := &runner{}

	if err := r.env(); err != nil {
		return nil, err
	}

	if err := r.loadClient(); err != nil {
		return nil, err
	}

	if err := r.loadUploader(); err != nil {
		return nil, err
	}

	if err := r.loadConcluder(); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *runner) run() error {
	obj, path := r.recorder.Record()

	if err := r.uploader.Upload(path); err != nil {
		return err
	}

	if err := r.concluder.Send(obj); err != nil {
		return err
	}

	return nil
}

func (r *runner) loadClient() error {
	r.ch = make(chan client.Object, 1024)

	var err error
	r.cli, err = client.NewClient(r.ch)
	if err != nil {
		return err
	}

	if err := r.cli.Watch(client.Key{
		Type:      "",
		Name:      "",
		Namespace: "",
	}); err != nil {
		return err
	}
	return nil
}

func (r *runner) loadRecorder() (err error) {
	r.recorder, err = recorder.NewRecorder(r.ch)
	return err
}

func (r *runner) loadUploader() (err error) {
	r.uploader, err = uploader.NewS3()
	return err
}

func (r *runner) loadConcluder() (err error) {
	r.concluder, err = concluder.NewRabbit()
	return err
}

func (r *runner) env() error {
	r.addr = os.Getenv("LOGGER_ADDRESS")
	if r.addr == "" {
		return fmt.Errorf("set the 'LOGGER_ADDRESS' environment variable to connect to Gimulator, like localhost:3030")
	}
	return nil
}
