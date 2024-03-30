package util

import (
	"github.com/gofrs/flock"
	"github.com/gruntwork-io/go-commons/errors"
)

type Lockfile struct {
	*flock.Flock
}

func NewLockfile(filename string) *Lockfile {
	return &Lockfile{
		flock.New(filename),
	}
}

func (lockfile *Lockfile) Unlock() error {
	if err := lockfile.Flock.Unlock(); err != nil {
		return errors.WithStackTrace(err)
	}
	return nil
}

func (lockfile *Lockfile) Lock() error {
	if locked, err := lockfile.Flock.TryLock(); err != nil {
		return errors.WithStackTrace(err)
	} else if !locked {
		return errors.Errorf("unable to lock file %s", lockfile.Path())
	}
	return nil
}

func AcquireLockfile(filename string) (*Lockfile, error) {
	lockfile := NewLockfile(filename)
	err := lockfile.Lock()
	return lockfile, err
}
