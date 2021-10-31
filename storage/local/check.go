package local

import (
	"os"
	"path/filepath"
)

// Check ...
func (l *Local) Check(project, version, fileName string) (exist bool, err error) {
	directory := filepath.Join(l.Path, project, version, fileName)
	return check(directory)
}

func check(path string) (exist bool, err error) {
	_, err = os.Stat(path)
	switch {
	case os.IsNotExist(err):
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, nil
	}
}
