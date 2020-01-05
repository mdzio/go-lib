package releng

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"
)

// Archive bundles files into an archive.
func Archive(dest string, specs []CopySpec) {
	log.Info("Building archive: ", dest)

	// select archive format
	var b arcBuilder
	dl := strings.ToLower(dest)
	if strings.HasSuffix(dl, ".zip") {
		b = newZipBuilder(dest)
	} else if strings.HasSuffix(dl, ".tgz") || strings.HasSuffix(dl, ".tar.gz") {
		b = newTgzBuilder(dest)
	} else {
		Must(fmt.Errorf("Unsupported archive format: %s", dest))
	}
	defer func() {
		b.close()
	}()

	// expand wildcards
	specs = expand(specs)

	// add files
	for _, s := range specs {
		log.Debug("Adding: ", s.Inc)
		s.DstDir = strings.TrimSuffix(s.DstDir, "/")
		b.add(s)
	}
}

// generic archive builder
type arcBuilder interface {
	add(e CopySpec)
	close()
}

// helper for creating directory entries
type dirAdder interface {
	addDir(path string)
}

type dirCreator struct {
	dirAdder
	dirs map[string]struct{}
}

func (c *dirCreator) addDirAll(dir string) {
	// no directory?
	if dir == "" {
		return
	}
	// directory already created?
	_, ok := c.dirs[dir]
	if ok {
		return
	}
	// create parent directory
	parentDir := path.Dir(dir)
	if parentDir != "." {
		c.addDirAll(parentDir)
	}
	// create directory
	c.dirAdder.addDir(dir)
	// remember created directories
	if c.dirs == nil {
		c.dirs = make(map[string]struct{})
	}
	c.dirs[dir] = struct{}{}
}

// zipBuilder implements arcBuilder.
type zipBuilder struct {
	f  *os.File
	zw *zip.Writer
	dirCreator
}

func newZipBuilder(path string) *zipBuilder {
	// start zip file
	f, err := os.Create(path)
	Must(err)
	zb := &zipBuilder{f: f, zw: zip.NewWriter(f)}
	zb.dirCreator.dirAdder = zb
	return zb
}

func (b *zipBuilder) addDir(dir string) {
	// create directory
	h := &zip.FileHeader{
		Name:     dir + "/",
		Modified: time.Now(), // set timestamp to build time
	}
	_, err := b.zw.CreateHeader(h)
	Must(err)
}

func (b *zipBuilder) add(e CopySpec) {
	// create directory entries, if needed
	b.addDirAll(e.DstDir)

	// open source file
	f, err := os.Open(e.Inc)
	Must(err)
	defer f.Close()

	// fill info header
	h := &zip.FileHeader{
		Name:     path.Join(e.DstDir, path.Base(e.Inc)),
		Modified: time.Now(), // set timestamp to build time
		Method:   zip.Deflate,
	}

	// create header
	w, err := b.zw.CreateHeader(h)
	Must(err)

	// write content
	_, err = io.Copy(w, f)
	Must(err)
}

func (b *zipBuilder) close() {
	// end zip file
	Must(b.zw.Close())
	Must(b.f.Close())
}

// tgzBuilder implements arcBuilder.
type tgzBuilder struct {
	f  *os.File
	gw *gzip.Writer
	tw *tar.Writer
	dirCreator
}

func newTgzBuilder(path string) *tgzBuilder {
	// start tgz file
	f, err := os.Create(path)
	Must(err)
	gw := gzip.NewWriter(f)
	tw := tar.NewWriter(gw)
	b := &tgzBuilder{f: f, gw: gw, tw: tw}
	b.dirCreator.dirAdder = b
	return b
}

func (b *tgzBuilder) addDir(dir string) {
	// create directory
	h := &tar.Header{
		Name:    dir + "/",
		ModTime: time.Now(), // set timestamp to build time
		Uname:   "root",
		Gname:   "root",
		Mode:    0755,
	}
	// create header
	Must(b.tw.WriteHeader(h))
}

func (b *tgzBuilder) add(e CopySpec) {
	// create directory entries, if needed
	b.addDirAll(e.DstDir)

	// open source file
	f, err := os.Open(e.Inc)
	Must(err)
	defer f.Close()

	// fill info header
	i, err := f.Stat()
	Must(err)
	var mode int64
	if e.Exe {
		mode = 0755
	} else {
		mode = 0644
	}
	h := &tar.Header{
		Name:    path.Join(e.DstDir, path.Base(e.Inc)),
		ModTime: time.Now(), // set timestamp to build time
		Uname:   "root",
		Gname:   "root",
		Mode:    mode,
		Size:    i.Size(),
	}

	// create header
	Must(b.tw.WriteHeader(h))

	// write content
	_, err = io.Copy(b.tw, f)
	Must(err)
}

func (b *tgzBuilder) close() {
	// end tgz file
	Must(b.tw.Close())
	Must(b.gw.Close())
	Must(b.f.Close())
}
