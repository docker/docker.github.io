// +build linux freebsd solaris

package daemon

import "github.com/docker/docker/container"

// excludeByIsolation is a platform specific helper function to support PS
// filtering by Isolation. This is a Windows-only concept, so is a no-op on Unix.
func excludeByIsolation(container *container.Container, ctx *listContext) iterationAction {
	return includeContainer
}
