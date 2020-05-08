package path

import (
	"os"
	"path/filepath"
)

var MonitororBaseDir = ""

func init() {
	MonitororBaseDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
}

func ToAbsolute(basedir, path string) string {
	if filepath.IsAbs(path) {
		path, _ = filepath.Abs(path)
	} else {
		path = filepath.Join(basedir, path)
	}

	return path
}
