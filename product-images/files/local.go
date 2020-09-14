package files

import (
	"io"
	"os"
	"path/filepath"

	"golang.org/x/xerrors"
)

type Local struct {
	maxFileSize int // maximum no. of bytes
	basePath    string
}

// Create New Local Instance
func NewLocal(basePath string, maxSize int) (*Local, error) {
	p, err := filepath.Abs(basePath)

	if err != nil {
		return nil, err
	}

	return &Local{maxSize, p}, nil
}

func (l *Local) Save(path string, contents io.Reader) error {
	// get full path
	fp := l.fullPath(path)

	// get directory
	d := filepath.Dir(fp)
	err := os.MkdirAll(d, os.ModePerm)
	if err != nil {
		return xerrors.Errorf("Unable to create directory: %w", err)
	}

	_, err = os.Stat(fp)
	if err == nil {
		err = os.Remove(fp)
		if err != nil {
			return xerrors.Errorf("Unable to delete file %w", err)
		}
	} else if !os.IsNotExist(err) {
		return xerrors.Errorf("Unable to get file info: %w", err)
	}

	// create a new file at the path
	f, err := os.Create(fp)
	if err != nil {
		return xerrors.Errorf("Unable to create file: %w", err)
	}

	defer f.Close()

	// write contents to the new file
	// ensure were not exceeding max bytes
	_, err = io.Copy(f, contents)
	if err != nil {
		return xerrors.Errorf("Unable to write to file : %w", err)
	}

	return nil
}

// Get the file
func (l *Local) Get(path string) (*os.File, error) {
	fp := l.fullPath(path)

	f, err := os.Open(fp)
	if err != nil {
		return nil, xerrors.Errorf("Unable to open file: %w", err)
	}
	return f, nil
}

func (l *Local) fullPath(path string) string {
	return filepath.Join(l.basePath, path)
}
