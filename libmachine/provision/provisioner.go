package provision

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/docker/machine/drivers"
	"github.com/docker/machine/libmachine/auth"
	"github.com/docker/machine/libmachine/provision/pkgaction"
	"github.com/docker/machine/libmachine/swarm"
)

var provisioners = make(map[string]*RegisteredProvisioner)

// Distribution specific actions
type Provisioner interface {
	// Create the files for the daemon to consume configuration settings (return struct of content and path)
	GenerateDockerConfig(dockerPort int, authConfig auth.AuthOptions) (*DockerConfig, error)

	// Get the directory where the settings files for docker are to be found
	GetDockerConfigDir() string

	// Run a package action e.g. install
	Package(name string, action pkgaction.PackageAction) error

	// Get Hostname
	Hostname() (string, error)

	// Set hostname
	SetHostname(hostname string) error

	// Figure out if this is the right provisioner to use based on /etc/os-release info
	CompatibleWithHost() bool

	// Do the actual provisioning piece:
	//     1. Set the hostname on the instance.
	//     2. Install Docker if it is not present.
	//     3. Configure the daemon to accept connections over TLS.
	//     4. Copy the needed certificates to the server and local config dir.
	//     5. Configure / activate swarm if applicable.
	Provision(swarmConfig swarm.SwarmOptions, authConfig auth.AuthOptions) error

	// Perform action on a named service e.g. stop
	Service(name string, action pkgaction.ServiceAction) error

	// Get the driver which is contained in the provisioner.
	GetDriver() drivers.Driver

	// Short-hand for accessing an SSH command from the driver.
	SSHCommand(args ...string) (*exec.Cmd, error)

	// Set the OS Release info depending on how it's represented
	// internally
	SetOsReleaseInfo(info *OsRelease)
}

// Detection
type RegisteredProvisioner struct {
	New func(d drivers.Driver) Provisioner
}

func Register(name string, p *RegisteredProvisioner) {
	provisioners[name] = p
}

func DetectProvisioner(d drivers.Driver) (Provisioner, error) {
	var (
		osReleaseOut bytes.Buffer
	)
	catOsReleaseCmd, err := drivers.GetSSHCommandFromDriver(d, "cat /etc/os-release")
	if err != nil {
		return nil, fmt.Errorf("Error getting SSH command: %s", err)
	}

	// Normally I would just use Output() for this, but d.GetSSHCommand
	// defaults to sending the output of the command to stdout in debug
	// mode, so that will be broken if we don't set it ourselves.
	catOsReleaseCmd.Stdout = &osReleaseOut

	if err := catOsReleaseCmd.Run(); err != nil {
		return nil, fmt.Errorf("Error running SSH command to get /etc/os-release: %s", err)
	}

	osReleaseInfo, err := NewOsRelease(osReleaseOut.Bytes())
	if err != nil {
		return nil, fmt.Errorf("Error parsing /etc/os-release file: %s", err)
	}

	for _, p := range provisioners {
		provisioner := p.New(d)
		provisioner.SetOsReleaseInfo(osReleaseInfo)

		if provisioner.CompatibleWithHost() {
			return provisioner, nil
		}
	}

	return nil, ErrDetectionFailed
}
