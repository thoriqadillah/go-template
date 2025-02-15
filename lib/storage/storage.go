package storage

import (
	"io"
	"log"
	"os"
)

type option struct {
	// The directory to store temporary files for development purposes
	tmpDir string
}

type Option func(*option)

func WithTmpDir(tmpDir string) Option {
	return func(o *option) {
		o.tmpDir = tmpDir
	}
}

type Storage interface {
	Serve(filename string) (io.ReadCloser, error)
	Upload(filename string, src io.Reader) error
	Delete(filename string) error
}

type Factory func(option *option) Storage

var providers = map[string]Factory{}

func New(name string, opts ...Option) Storage {
	tmp, _ := os.Getwd()
	opt := &option{
		tmpDir: tmp + "/lib/storage/tmp",
	}

	for _, option := range opts {
		option(opt)
	}

	provider, ok := providers[name]
	if !ok {
		log.Fatalf("Storage provider %s not found", name)
		return nil
	}

	return provider(opt)
}

func register(name string, impl Factory) {
	providers[name] = impl
}
