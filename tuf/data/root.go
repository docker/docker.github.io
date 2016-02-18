package data

import (
	"fmt"
	"time"

	"github.com/docker/go/canonical/json"
)

// SignedRoot is a fully unpacked root.json
type SignedRoot struct {
	Signatures []Signature
	Signed     Root
	Dirty      bool
}

// Root is the Signed component of a root.json
type Root struct {
	Type               string               `json:"_type"`
	Version            int                  `json:"version"`
	Expires            time.Time            `json:"expires"`
	Keys               Keys                 `json:"keys"`
	Roles              map[string]*RootRole `json:"roles"`
	ConsistentSnapshot bool                 `json:"consistent_snapshot"`
}

// isValidRootStructure returns an error, or nil, depending on whether the content of the struct
// is valid for root metadata.  This does not check signatures or expiry, just that
// the metadata content is valid.
func isValidRootStructure(r Root) error {
	expectedType := TUFTypes[CanonicalRootRole]
	if r.Type != expectedType {
		return ErrInvalidMetadata{
			role: CanonicalRootRole, msg: fmt.Sprintf("expected type %s, not %s", expectedType, r.Type)}
	}
	if len(r.Roles) < 4 {
		return ErrInvalidMetadata{role: CanonicalRootRole, msg: "does not have all required roles"}
	} else if len(r.Roles) > 4 {
		return ErrInvalidMetadata{role: CanonicalRootRole, msg: "specifies too many roles"}
	}

	for _, roleName := range BaseRoles {
		roleObj, ok := r.Roles[roleName]
		if !ok || roleObj == nil {
			return ErrInvalidMetadata{
				role: CanonicalRootRole, msg: fmt.Sprintf("missing %s role specification", roleName)}
		}
		if err := isValidRootRoleStructure(*roleObj, r.Keys); err != nil {
			return ErrInvalidMetadata{
				role: CanonicalRootRole,
				msg:  fmt.Sprintf("role %s: %s", roleName, err.Error()),
			}
		}
	}
	return nil
}

// NewRoot initializes a new SignedRoot with a set of keys, roles, and the consistent flag
func NewRoot(keys map[string]PublicKey, roles map[string]*RootRole, consistent bool) (*SignedRoot, error) {
	signedRoot := &SignedRoot{
		Signatures: make([]Signature, 0),
		Signed: Root{
			Type:               TUFTypes[CanonicalRootRole],
			Version:            0,
			Expires:            DefaultExpires(CanonicalRootRole),
			Keys:               keys,
			Roles:              roles,
			ConsistentSnapshot: consistent,
		},
		Dirty: true,
	}

	return signedRoot, nil
}

// ToSigned partially serializes a SignedRoot for further signing
func (r *SignedRoot) ToSigned() (*Signed, error) {
	s, err := defaultSerializer.MarshalCanonical(r.Signed)
	if err != nil {
		return nil, err
	}
	// cast into a json.RawMessage
	signed := json.RawMessage{}
	err = signed.UnmarshalJSON(s)
	if err != nil {
		return nil, err
	}
	sigs := make([]Signature, len(r.Signatures))
	copy(sigs, r.Signatures)
	return &Signed{
		Signatures: sigs,
		Signed:     signed,
	}, nil
}

// MarshalJSON returns the serialized form of SignedRoot as bytes
func (r *SignedRoot) MarshalJSON() ([]byte, error) {
	signed, err := r.ToSigned()
	if err != nil {
		return nil, err
	}
	return defaultSerializer.Marshal(signed)
}

// RootFromSigned fully unpacks a Signed object into a SignedRoot and ensures
// that it is a valid SignedRoot
func RootFromSigned(s *Signed) (*SignedRoot, error) {
	r := Root{}
	if err := defaultSerializer.Unmarshal(s.Signed, &r); err != nil {
		return nil, err
	}
	if err := isValidRootStructure(r); err != nil {
		return nil, err
	}
	sigs := make([]Signature, len(s.Signatures))
	copy(sigs, s.Signatures)
	return &SignedRoot{
		Signatures: sigs,
		Signed:     r,
	}, nil
}
