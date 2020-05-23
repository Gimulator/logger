package uploader

type Uploader interface {
	Upload(path, key string) error
}
