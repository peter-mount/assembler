package common

import "github.com/peter-mount/go-kernel/v2/log"

var Version string

// VersionService is a service which ensures that logging is initialised and
// if verbose then log the application header.
// Reference this as the first service in the kernel.Launch() function
type VersionService struct{}

func (r *VersionService) Start() error {
	log.Println(Version)
	return nil
}
