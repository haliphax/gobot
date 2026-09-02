// Package internal contains application internals
package internal

import "github.com/spf13/afero"

var fs afero.Fs

// SetFs sets the file system singleton for Gobot
// (useful for unit tests)
func SetFs(newFs afero.Fs) {
	fs = newFs
}

// GetFs gets the file system singleton for Gobot
func GetFs() afero.Fs {
	return fs
}
