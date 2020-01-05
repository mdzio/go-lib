package releng

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mdzio/go-lib/logging"
)

var log = logging.Get("releng")

// CopySpec specifies files to copy.
type CopySpec struct {
	Inc    string
	DstDir string
	Exe    bool
}

// Must exits on error.
func Must(err error) {
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

// RequireFiles checks for existence and readability of files.
func RequireFiles(files []string) {
	for _, filePath := range files {
		file, err := os.Open(filePath)
		if err != nil {
			Must(fmt.Errorf("Required file not found: %s", filePath))
		}
		file.Close()
	}
}

// Mkdir creates the specified directory.
func Mkdir(path string) {
	log.Info("Creating directory: ", path)
	err := os.MkdirAll(path, 0777)
	Must(err)
}

// Getwd returns the current working directory.
func Getwd() string {
	wd, err := os.Getwd()
	Must(err)
	return wd
}

// WriteFile writes the data to the specified file.
func WriteFile(path string, data []byte) {
	log.Info("Writing file: ", path)
	Must(ioutil.WriteFile(path, data, 0666))
}

// expand expands wildcards in CopySpec's.
func expand(specs []CopySpec) []CopySpec {
	var ret []CopySpec
	for _, s := range specs {
		fs, err := filepath.Glob(s.Inc)
		Must(err)
		cnt := 0
		for _, f := range fs {
			// is directory?
			inf, err := os.Lstat(f)
			Must(err)
			if inf.IsDir() {
				continue
			}
			// normalize path separators
			if os.PathSeparator != '/' {
				f = strings.ReplaceAll(f, string(os.PathSeparator), "/")
			}
			ret = append(ret, CopySpec{Inc: f, DstDir: s.DstDir, Exe: s.Exe})
			cnt++
		}
		if cnt == 0 {
			log.Warning("No files matches pattern: ", s.Inc)
		}
	}
	return ret
}
