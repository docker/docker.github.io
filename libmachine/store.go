package libmachine

type Store interface {
	// Exists returns whether a machine exists or not
	Exists(name string) (bool, error)
	// GetActive returns the active host
	GetActive() (*Host, error)
	// GetPath returns the path to the store
	GetPath() string
	// GetCACertPath returns the CA certificate
	GetCACertificatePath() (string, error)
	// GetPrivateKeyPath returns the private key
	GetPrivateKeyPath() (string, error)
	// IsActive returns whether the host is active or not
	IsActive(host *Host) (bool, error)
	// List returns a list of hosts
	List() ([]*Host, error)
	// Load loads a host by name
	Get(name string) (*Host, error)
	// Remove removes a machine from the store
	Remove(name string, force bool) error
	// RemoveActive removes the active machine from the store
	RemoveActive() error
	// Save persists a machine in the store
	Save(host *Host) error
	// SetActive sets the specified host as the active host
	SetActive(host *Host) error
}
