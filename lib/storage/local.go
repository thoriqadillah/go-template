package storage

import (
	"io"
	"os"
	"path/filepath"
)

type local struct {
	tmpDir string
}

func createLocal(option *option) Storage {
	return &local{
		tmpDir: option.tmpDir,
	}
}

func (l *local) Serve(filename string) (io.ReadCloser, error) {
	return os.Open(l.tmpDir + "/" + filename)
}

func (l *local) Upload(filename string, src io.Reader) error {
	file := filepath.Join(l.tmpDir, filename)

	dst, err := os.Create(file)
	if err != nil {
		return err
	}

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	return nil
}

func (l *local) Delete(filename string) error {
	return os.Remove(l.tmpDir + "/" + filename)
}

func init() {
	register("local", createLocal)
}
