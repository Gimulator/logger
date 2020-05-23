package main

import (
	"fmt"
	"os"
	"strings"

	client "github.com/Gimulator/client-go"
	"github.com/Gimulator/logger/concluder"
	"github.com/Gimulator/logger/recorder"
	"github.com/Gimulator/logger/uploader"
)

func readArgs() (noConcluder bool, noUploader bool) {
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		if strings.Contains(arg, "no-concluder") {
			noConcluder = true
		}
		if strings.Contains(arg, "no-uploader") {
			noUploader = true
		}
	}
	return
}

func main() {
	noConcluder, noUploader := readArgs()

	runner, err := newRunner(noConcluder, noUploader)
	if err != nil {
		panic(err)
	}

	if err := runner.run(); err != nil {
		panic(err)
	}
}

type runner struct {
	id     string
	endKey string
	ch     chan client.Object
	cli    *client.Client

	recorder  *recorder.Recorder
	uploader  uploader.Uploader
	concluder concluder.Concluder
}

func newRunner(noConcluder, noUploader bool) (*runner, error) {
	r := &runner{}

	if err := r.env(); err != nil {
		return nil, err
	}

	if err := r.loadClient(); err != nil {
		return nil, err
	}

	if err := r.loadRecorder(); err != nil {
		return nil, err
	}

	if !noUploader {
		if err := r.loadUploader(noUploader); err != nil {
			return nil, err
		}
	}

	if !noConcluder {
		if err := r.loadConcluder(noConcluder); err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (r *runner) run() error {
	obj, err := r.recorder.Record()
	if err != nil {
		return err
	}

	if r.uploader != nil {
		if err := r.uploader.Upload(r.recorder.LogFilePath(), r.id); err != nil {
			return err
		}
	}

	if r.concluder != nil {
		if err := r.concluder.Send(obj); err != nil {
			return err
		}
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
	r.recorder, err = recorder.NewRecorder(r.ch, r.endKey)
	return err
}

func (r *runner) loadUploader(noUploader bool) (err error) {
	r.uploader = nil
	if !noUploader {
		r.uploader, err = uploader.NewS3()
	}
	return err
}

func (r *runner) loadConcluder(noConcluder bool) (err error) {
	r.concluder = nil
	if !noConcluder {
		r.concluder, err = concluder.NewRabbit()
	}
	return err
}

func (r *runner) env() error {
	r.endKey = os.Getenv("ROOM_END_OF_GAME_KEY")
	if r.endKey == "" {
		return fmt.Errorf("set the 'ROOM_END_OF_GAME_KEY' environment variable to detect when game is over")
	}

	r.id = os.Getenv("ROOM_ID")
	if r.id == "" {
		return fmt.Errorf("set the 'ROOM_ID' environment variable to record")
	}
	return nil
}
