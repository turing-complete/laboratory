package file

import (
	"errors"
	"fmt"
	"os"

	"github.com/ready-steady/hdf5"
)

func Create(path string) (*hdf5.File, error) {
	if len(path) == 0 {
		return nil, errors.New("expected a filename")
	}

	return hdf5.Create(path)
}

func Open(path string) (*hdf5.File, error) {
	if len(path) == 0 {
		return nil, errors.New("expected a filename")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.New(fmt.Sprintf("the file “%s” does not exist", path))
	}

	return hdf5.Open(path)
}
