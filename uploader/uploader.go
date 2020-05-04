package uploader

type Uploader interface {
	Upload(path string) error
}
