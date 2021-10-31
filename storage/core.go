package storage

import (
	localStorage "goproxy/storage/local"
	s3Storage "goproxy/storage/s3"

	"github.com/spf13/viper"
)

const (
	local = "local"
	s3    = "s3"
)

// Operations ...
type Operations interface {
	Create(project, version, fileName string, data []byte) (err error)
	Get(project, version, fileName string) (data []byte, err error)
	Check(project, version, fileName string) (exist bool, err error)
}

var storage Operations

// Init ...
func Init() {
	switch viper.GetString("storage.provider") {
	case local:
		storage = &localStorage.Local{
			Path: viper.GetString("storage.local.path"),
		}
	case s3:
		storage = &s3Storage.S3{}
	default:
		panic("The storage provider is not supported")
	}
}

// Use ...
func Use() (op Operations) {
	return storage
}
