package releng

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/mdzio/go-logging"
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
	Must(os.WriteFile(path, data, 0666))
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

// WriteTemplate reads a template file, executes it with the provided data,
// and writes the result to the specified output file.
func WriteTemplate(templatePath, outputPath string, data interface{}) {
	log.Debug("Processing template: ", templatePath)

	// Read template file
	templateContent, err := os.ReadFile(templatePath)
	Must(err)

	// Parse and execute template
	tmpl, err := template.New(filepath.Base(templatePath)).Parse(string(templateContent))
	Must(err)

	// Create output file
	outputFile, err := os.Create(outputPath)
	Must(err)
	defer outputFile.Close()

	// Execute template with data and write to output file
	err = tmpl.Execute(outputFile, data)
	Must(err)

	log.Info("Written file from template: ", outputPath)
}
