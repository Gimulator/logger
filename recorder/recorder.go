package recorder

import (
	"fmt"
	"os"
	"path/filepath"

	client "github.com/Gimulator/client-go"
)

type Recorder struct {
	ch     chan client.Object
	file   *os.File
	path   string
	endKey string
}

func NewRecorder(ch chan client.Object, endKey string) (*Recorder, error) {
	if ch == nil {
		return nil, fmt.Errorf("invalid channel to read objects")
	}
	r := &Recorder{
		ch:     ch,
		endKey: endKey,
	}

	err := r.env()
	if err != nil {
		return nil, err
	}

	r.file, err = os.Create(r.path)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Recorder) env() error {
	dir := os.Getenv("LOGGER_RECORDER_DIR")
	if dir == "" {
		return fmt.Errorf("set the 'LOGGER_RECORDER_DIR' environment variable for storing logs in it")
	}
	r.path = filepath.Join(dir, "logger.log")

	return nil
}

func (r *Recorder) Record() (client.Object, error) {
	defer r.file.Close()

	for {
		obj := <-r.ch
		_, err := r.file.WriteString(fmt.Sprintf("%v\n", obj))
		if err != nil {
			return client.Object{}, err
		}

		if err := r.file.Sync(); err != nil {
			return client.Object{}, err
		}

		if obj.Key.Type == r.endKey {
			return obj, nil
		}
	}
}

func (r *Recorder) LogFilePath() string {
	return r.path
}
