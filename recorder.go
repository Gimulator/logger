package logger

import (
	"fmt"
	"os"
	"path/filepath"

	client "github.com/Gimulator/client-go"
)

type Recorder struct {
	ch     chan client.Object
	file   *os.File
	endKey string
}

func NewRecorder(ch chan client.Object, dir string, endKey string) (*Recorder, error) {
	path := filepath.Join(dir, "logger.log")

	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	return &Recorder{
		ch:     ch,
		file:   file,
		endKey: endKey,
	}, nil
}

func (r *Recorder) Record() {
	for {
		obj := <-r.ch
		r.file.WriteString(fmt.Sprintf("%v\n", obj))
		r.file.Sync()

		if obj.Key.Type == r.endKey {
			break
		}
	}
}
