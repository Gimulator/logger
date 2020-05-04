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

func NewRecorder(ch chan client.Object) (*Recorder, error) {
	r := &Recorder{
		ch: ch,
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
	dir := os.Getenv("LOGGER_RECORD_DIR")
	if dir == "" {
		return fmt.Errorf("set the 'LOGGER_RECORD_DIR' environment variable for storing logs in it")
	}
	r.path = filepath.Join(dir, "logger.log")

	endKey := os.Getenv("LOGGER_RECORD_END_KEY")
	if endKey == "" {
		return fmt.Errorf("set the 'LOGGER_RECORD_END_KEY' environment variable to exit the program")
	}
	r.endKey = endKey

	return nil
}

func (r *Recorder) Record() (client.Object, string) {
	for {
		obj := <-r.ch
		r.file.WriteString(fmt.Sprintf("%v\n", obj))
		r.file.Sync()

		if obj.Key.Type == r.endKey {
			return obj, r.path
		}
	}
}
