// +build !linux,!windows

package sandbox

import "errors"

var (
	// ErrNotImplemented is for platforms which don't implement sandbox
	ErrNotImplemented = errors.New("not implemented")
)

// NewSandbox provides a new sandbox instance created in an os specific way
// provided a key which uniquely identifies the sandbox
func NewSandbox(key string, osCreate bool) (Sandbox, error) {
	return nil, ErrNotImplemented
}

// GenerateKey generates a sandbox key based on the passed
// container id.
func GenerateKey(containerID string) string {
	return ""
}
