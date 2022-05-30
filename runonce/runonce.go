package runonce

import (
	"errors"
	"os"
	"path/filepath"
)

var (
	ErrAlreadyRunning = errors.New("program already running")
)

type RunOnce struct {
	name string
	f    *os.File
}

func New(name string) *RunOnce {
	s := &RunOnce{
		name: name,
	}

	return s
}

func (s *RunOnce) LockFname() string {
	return filepath.Join(os.TempDir(), s.name) + ".lck"
}
