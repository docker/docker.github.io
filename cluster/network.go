package cluster

import (
	"strings"

	"github.com/docker/docker/pkg/stringid"
	"github.com/samalba/dockerclient"
)

// Network is exported
type Network struct {
	dockerclient.NetworkResource

	Engine *Engine
}

// Networks represents a map of networks
type Networks []*Network

// Uniq returns all uniq networks
func (networks Networks) Uniq() Networks {
	tmp := make(map[string]*Network)
	for _, network := range networks {
		tmp[network.ID] = network
	}
	uniq := Networks{}
	for _, network := range tmp {
		uniq = append(uniq, network)
	}
	return uniq
}

// Get returns a network using it's ID or Name
func (networks Networks) Get(IDOrName string) *Network {
	// Abort immediately if the name is empty.
	if len(IDOrName) == 0 {
		return nil
	}

	// Match exact or short Network ID.
	for _, network := range networks {
		if network.ID == IDOrName || stringid.TruncateID(network.ID) == IDOrName {
			return network
		}
	}

	candidates := []*Network{}

	// Match name, /name or engine/name.
	for _, network := range networks {
		if network.Name == IDOrName || network.Engine.ID+"/"+network.Name == IDOrName || network.Engine.Name+"/"+network.Name == IDOrName {
			candidates = append(candidates, network)
		}
	}

	if size := len(candidates); size == 1 {
		return candidates[0]
	} else if size > 1 {
		return nil
	}

	// Match name, /name or engine/name.
	for _, network := range networks {
		if network.Name == "/"+IDOrName {
			return network
		}
	}

	// Match Network ID prefix.
	for _, network := range networks {
		if strings.HasPrefix(network.ID, IDOrName) {
			candidates = append(candidates, network)
		}
	}

	if len(candidates) == 1 {
		return candidates[0]
	}

	return nil

}
