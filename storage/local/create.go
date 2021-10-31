package local

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

// Create ...
func (l *Local) Create(project, version, fileName string, data []byte) (err error) {
	directory := filepath.Join(l.Path, project, version)
	// if _, err := os.Stat(directory); os.IsNotExist(err) {
	// 	os.MkdirAll(directory, 0775)
	// }
	exist, err := check(directory)
	if err != nil {
		return err
	}
	if !exist {
		os.MkdirAll(directory, 0775)
	}
	file, err := os.Create(filepath.Join(directory, fileName))
	defer file.Close()
	if err != nil {
		log.Error(err)
		return err
	}
	if _, err := file.Write(data); err != nil {
		log.Error(err)
		return err
	}
	return nil
}
