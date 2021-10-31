package local

import (
	"os"
	"path/filepath"
)

// Get ...
func (l *Local) Get(project, version, fileName string) (data []byte, err error) {
	directory := filepath.Join(l.Path, project, version)
	return os.ReadFile(filepath.Join(directory, fileName))
}
