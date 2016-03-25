package libcontainerd

import "github.com/docker/docker/pkg/locker"

type remote struct {
}

func (r *remote) Client(b Backend) (Client, error) {
	c := &client{
		clientCommon: clientCommon{
			backend:    b,
			containers: make(map[string]*container),
			locker:     locker.New(),
		},
	}
	return c, nil
}

func (r *remote) Cleanup() {
}

// New creates a fresh instance of libcontainerd remote.
func New(_ string, _ ...RemoteOption) (Remote, error) {
	return &remote{}, nil
}
