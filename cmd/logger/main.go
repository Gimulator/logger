package main

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	client "github.com/Gimulator/client-go"
	"github.com/Gimulator/logger/concluder"
	"github.com/Gimulator/logger/recorder"
	"github.com/Gimulator/logger/uploader"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)

	formatter := &logrus.TextFormatter{
		TimestampFormat:  "2006-01-02 15:04:05",
		FullTimestamp:    true,
		PadLevelText:     true,
		QuoteEmptyFields: true,
		ForceQuote:       false,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return "", fmt.Sprintf(" %s:%d\t", path.Base(f.File), f.Line)
		},
	}
	logrus.SetFormatter(formatter)
}

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
	roomID     string
	roomEndKey string
	ch         chan client.Object
	cli        *client.Client

	recorder  *recorder.Recorder
	uploader  uploader.Uploader
	concluder concluder.Concluder

	log *logrus.Entry
}

func newRunner(noConcluder, noUploader bool) (*runner, error) {
	r := &runner{
		log: logrus.WithField("entity", "runner"),
	}
	r.log.Info(fmt.Sprintf("starting with options --no-concluder=%v, --no-uploader=%v", noConcluder, noUploader))

	r.log.Info("starting to read environment variables")
	if err := r.env(); err != nil {
		r.log.WithError(err).Error("could not read environment variables")
		return nil, err
	}

	r.log.Info("starting to initiate client")
	if err := r.loadClient(); err != nil {
		r.log.WithError(err).Error("could not initiate client")
		return nil, err
	}

	r.log.Info("starting to initiate recorder")
	if err := r.loadRecorder(); err != nil {
		r.log.WithError(err).Error("could not initiate recorder")
		return nil, err
	}

	r.log.Info("starting to initiate uploader")
	if err := r.loadUploader(noUploader); err != nil {
		r.log.WithError(err).Error("could not initiate uploader")
		return nil, err
	}

	r.log.Info("starting to initiate concluder")
	if err := r.loadConcluder(noConcluder); err != nil {
		r.log.WithError(err).Error("could not initiate concluder")
		return nil, err
	}

	return r, nil
}

func (r *runner) run() error {
	r.log.Info("starting to run")

	r.log.Info("starting to record objects")
	obj, err := r.recorder.Record()
	if err != nil {
		r.log.WithError(err).Error("could not record objects")
		return err
	}

	r.log.Info("starting to upload log-file to s3")
	if r.uploader != nil {
		if err := r.uploader.Upload(r.recorder.LogFilePath(), r.roomID); err != nil {
			r.log.WithError(err).Error("could not upload log-file to s3")
			return err
		}
	} else {
		r.log.Debug("nil uploader")
	}

	r.log.Info("starting to send conclusion to rabbitMQ")
	if r.concluder != nil {
		if err := r.concluder.Send(obj); err != nil {
			r.log.WithError(err).Error("could not send conclusion to rabbitMQ")
			return err
		}
	} else {
		r.log.Debug("nil concluder")
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
	r.recorder, err = recorder.NewRecorder(r.ch, r.roomEndKey)
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
	r.roomEndKey = os.Getenv("ROOM_END_OF_GAME_KEY")
	if r.roomEndKey == "" {
		return fmt.Errorf("set the 'ROOM_END_OF_GAME_KEY' environment variable for detecting the end of the game")
	}

	r.roomID = os.Getenv("ROOM_ID")
	if r.roomID == "" {
		return fmt.Errorf("set the 'ROOM_ID' environment variable for recording and uploading logs")
	}
	return nil
}
