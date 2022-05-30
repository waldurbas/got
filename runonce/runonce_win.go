//go:build windows
// +build windows

package runonce

// ----------------------------------------------------------------------------------
// runonce_win.go (https://github.com/waldurbas/got)
// Copyright 2019,2022 by Waldemar Urbas
//-----------------------------------------------------------------------------------

import (
	"os"
)

// Trylock #
func (r *RunOnce) TryLock() error {
	if err := os.Remove(r.LockFname()); err != nil && !os.IsNotExist(err) {
		return ErrAlreadyRunning
	}

	f, err := os.OpenFile(r.LockFname(), os.O_EXCL|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	r.f = f

	return nil
}

// Unlock #
func (r *RunOnce) Unlock() error {
	if err := r.f.Close(); err != nil {
		return err
	}

	if err := os.Remove(r.LockFname()); err != nil {
		return err
	}

	return nil
}
