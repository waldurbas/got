//go:build linux
// +build linux

package runonce

// ----------------------------------------------------------------------------------
// runonce_lin.go (https://github.com/waldurbas/got)
// Copyright 2019,2022 by Waldemar Urbas
//-----------------------------------------------------------------------------------

import (
	"os"
	"syscall"
)

// TryLock #
func (r *RunOnce) Trylock() error {
	f, err := os.OpenFile(r.LockFname(), os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	r.f = f

	flock := syscall.Flock_t{
		Type: syscall.F_WRLCK,
		Pid:  int32(os.Getpid()),
	}

	if err := syscall.FcntlFlock(r.f.Fd(), syscall.F_SETLK, &flock); err != nil {
		return ErrAlreadyRunning
	}

	return nil
}

// Unlock #
func (r *RunOnce) Unlock() error {

	flock := syscall.Flock_t{
		Type: syscall.F_UNLCK,
		Pid:  int32(os.Getpid()),
	}

	if err := syscall.FcntlFlock(r.f.Fd(), syscall.F_SETLK, &flock); err != nil {
		return err
	}

	if err := r.f.Close(); err != nil {
		return err
	}

	if err := os.Remove(r.LockFname()); err != nil {
		return err
	}

	return nil
}
