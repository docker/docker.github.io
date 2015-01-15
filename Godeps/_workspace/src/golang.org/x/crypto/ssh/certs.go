// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"sort"
	"time"
)

// These constants from [PROTOCOL.certkeys] represent the algorithm names
// for certificate types supported by this package.
const (
	CertAlgoRSAv01      = "ssh-rsa-cert-v01@openssh.com"
	CertAlgoDSAv01      = "ssh-dss-cert-v01@openssh.com"
	CertAlgoECDSA256v01 = "ecdsa-sha2-nistp256-cert-v01@openssh.com"
	CertAlgoECDSA384v01 = "ecdsa-sha2-nistp384-cert-v01@openssh.com"
	CertAlgoECDSA521v01 = "ecdsa-sha2-nistp521-cert-v01@openssh.com"
)

// Certificate types distinguish between host and user
// certificates. The values can be set in the CertType field of
// Certificate.
const (
	UserCert = 1
	HostCert = 2
)

// Signature represents a cryptographic signature.
type Signature struct {
	Format string
	Blob   []byte
}

// CertTimeInfinity can be used for OpenSSHCertV01.ValidBefore to indicate that
// a certificate does not expire.
const CertTimeInfinity = 1<<64 - 1

// An Certificate represents an OpenSSH certificate as defined in
// [PROTOCOL.certkeys]?rev=1.8.
type Certificate struct {
	Nonce           []byte
	Key             PublicKey
	Serial          uint64
	CertType        uint32
	KeyId           string
	ValidPrincipals []string
	ValidAfter      uint64
	ValidBefore     uint64
	Permissions
	Reserved     []byte
	SignatureKey PublicKey
	Signature    *Signature
}

// genericCertData holds the key-independent part of the certificate data.
// Overall, certificates contain an nonce, public key fields and
// key-independent fields.
type genericCertData struct {
	Serial          uint64
	CertType        uint32
	KeyId           string
	ValidPrincipals []byte
	ValidAfter      uint64
	ValidBefore     uint64
	CriticalOptions []byte
	Extensions      []byte
	Reserved        []byte
	SignatureKey    []byte
	Signature       []byte
}

func marshalStringList(namelist []string) []byte {
	var to []byte
	for _, name := range namelist {
		s := struct{ N string }{name}
		to = append(to, Marshal(&s)...)
	}
	return to
}

func marshalTuples(tups map[string]string) []byte {
	keys := make([]string, 0, len(tups))
	for k := range tups {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var r []byte
	for _, k := range keys {
		s := struct{ K, V string }{k, tups[k]}
		r = append(r, Marshal(&s)...)
	}
	return r
}

func parseTuples(in []byte) (map[string]string, error) {
	tups := map[string]string{}
	var lastKey string
	var haveLastKey bool

	for len(in) > 0 {
		nameBytes, rest, ok := parseString(in)
		if !ok {
			return nil, errShortRead
		}
		data, rest, ok := parseString(rest)
		if !ok {
			return nil, errShortRead
		}
		name := string(nameBytes)

		// according to [PROTOCOL.certkeys], the names must be in
		// lexical order.
		if haveLastKey && name <= lastKey {
			return nil, fmt.Errorf("ssh: certificate options are not in lexical order")
		}
		lastKey, haveLastKey = name, true

		tups[name] = string(data)
		in = rest
	}
	return tups, nil
}

func parseCert(in []byte, privAlgo string) (*Certificate, error) {
	nonce, rest, ok := parseString(in)
	if !ok {
		return nil, errShortRead
	}

	key, rest, err := parsePubKey(rest, privAlgo)
	if err != nil {
		return nil, err
	}

	var g genericCertData
	if err := Unmarshal(rest, &g); err != nil {
		return nil, err
	}

	c := &Certificate{
		Nonce:       nonce,
		Key:         key,
		Serial:      g.Serial,
		CertType:    g.CertType,
		KeyId:       g.KeyId,
		ValidAfter:  g.ValidAfter,
		ValidBefore: g.ValidBefore,
	}

	for principals := g.ValidPrincipals; len(principals) > 0; {
		principal, rest, ok := parseString(principals)
		if !ok {
			return nil, errShortRead
		}
		c.ValidPrincipals = append(c.ValidPrincipals, string(principal))
		principals = rest
	}

	c.CriticalOptions, err = parseTuples(g.CriticalOptions)
	if err != nil {
		return nil, err
	}
	c.Extensions, err = parseTuples(g.Extensions)
	if err != nil {
		return nil, err
	}
	c.Reserved = g.Reserved
	k, err := ParsePublicKey(g.SignatureKey)
	if err != nil {
		return nil, err
	}

	c.SignatureKey = k
	c.Signature, rest, ok = parseSignatureBody(g.Signature)
	if !ok || len(rest) > 0 {
		return nil, errors.New("ssh: signature parse error")
	}

	return c, nil
}

type openSSHCertSigner struct {
	pub    *Certificate
	signer Signer
}

// NewCertSigner returns a Signer that signs with the given Certificate, whose
// private key is held by signer. It returns an error if the public key in cert
// doesn't match the key used by signer.
func NewCertSigner(cert *Certificate, signer Signer) (Signer, error) {
	if bytes.Compare(cert.Key.Marshal(), signer.PublicKey().Marshal()) != 0 {
		return nil, errors.New("ssh: signer and cert have different public key")
	}

	return &openSSHCertSigner{cert, signer}, nil
}

func (s *openSSHCertSigner) Sign(rand io.Reader, data []byte) (*Signature, error) {
	return s.signer.Sign(rand, data)
}

func (s *openSSHCertSigner) PublicKey() PublicKey {
	return s.pub
}

const sourceAddressCriticalOption = "source-address"

// CertChecker does the work of verifying a certificate. Its methods
// can be plugged into ClientConfig.HostKeyCallback and
// ServerConfig.PublicKeyCallback. For the CertChecker to work,
// minimally, the IsAuthority callback should be set.
type CertChecker struct {
	// SupportedCriticalOptions lists the CriticalOptions that the
	// server application layer understands. These are only used
	// for user certificates.
	SupportedCriticalOptions []string

	// IsAuthority should return true if the key is recognized as
	// an authority. This allows for certificates to be signed by other
	// certificates.
	IsAuthority func(auth PublicKey) bool

	// Clock is used for verifying time stamps. If nil, time.Now
	// is used.
	Clock func() time.Time

	// UserKeyFallback is called when CertChecker.Authenticate encounters a
	// public key that is not a certificate. It must implement validation
	// of user keys or else, if nil, all such keys are rejected.
	UserKeyFallback func(conn ConnMetadata, key PublicKey) (*Permissions, error)

	// HostKeyFallback is called when CertChecker.CheckHostKey encounters a
	// public key that is not a certificate. It must implement host key
	// validation or else, if nil, all such keys are rejected.
	HostKeyFallback func(addr string, remote net.Addr, key PublicKey) error

	// IsRevoked is called for each certificate so that revocation checking
	// can be implemented. It should return true if the given certificate
	// is revoked and false otherwise. If nil, no certificates are
	// considered to have been revoked.
	IsRevoked func(cert *Certificate) bool
}

// CheckHostKey checks a host key certificate. This method can be
// plugged into ClientConfig.HostKeyCallback.
func (c *CertChecker) CheckHostKey(addr string, remote net.Addr, key PublicKey) error {
	cert, ok := key.(*Certificate)
	if !ok {
		if c.HostKeyFallback != nil {
			return c.HostKeyFallback(addr, remote, key)
		}
		return errors.New("ssh: non-certificate host key")
	}
	if cert.CertType != HostCert {
		return fmt.Errorf("ssh: certificate presented as a host key has type %d", cert.CertType)
	}

	return c.CheckCert(addr, cert)
}

// Authenticate checks a user certificate. Authenticate can be used as
// a value for ServerConfig.PublicKeyCallback.
func (c *CertChecker) Authenticate(conn ConnMetadata, pubKey PublicKey) (*Permissions, error) {
	cert, ok := pubKey.(*Certificate)
	if !ok {
		if c.UserKeyFallback != nil {
			return c.UserKeyFallback(conn, pubKey)
		}
		return nil, errors.New("ssh: normal key pairs not accepted")
	}

	if cert.CertType != UserCert {
		return nil, fmt.Errorf("ssh: cert has type %d", cert.CertType)
	}

	if err := c.CheckCert(conn.User(), cert); err != nil {
		return nil, err
	}

	return &cert.Permissions, nil
}

// CheckCert checks CriticalOptions, ValidPrincipals, revocation, timestamp and
// the signature of the certificate.
func (c *CertChecker) CheckCert(principal string, cert *Certificate) error {
	if c.IsRevoked != nil && c.IsRevoked(cert) {
		return fmt.Errorf("ssh: certicate serial %d revoked", cert.Serial)
	}

	for opt, _ := range cert.CriticalOptions {
		// sourceAddressCriticalOption will be enforced by
		// serverAuthenticate
		if opt == sourceAddressCriticalOption {
			continue
		}

		found := false
		for _, supp := range c.SupportedCriticalOptions {
			if supp == opt {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("ssh: unsupported critical option %q in certificate", opt)
		}
	}

	if len(cert.ValidPrincipals) > 0 {
		// By default, certs are valid for all users/hosts.
		found := false
		for _, p := range cert.ValidPrincipals {
			if p == principal {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("ssh: principal %q not in the set of valid principals for given certificate: %q", principal, cert.ValidPrincipals)
		}
	}

	if !c.IsAuthority(cert.SignatureKey) {
		return fmt.Errorf("ssh: certificate signed by unrecognized authority")
	}

	clock := c.Clock
	if clock == nil {
		clock = time.Now
	}

	unixNow := clock().Unix()
	if after := int64(cert.ValidAfter); after < 0 || unixNow < int64(cert.ValidAfter) {
		return fmt.Errorf("ssh: cert is not yet valid")
	}
	if before := int64(cert.ValidBefore); cert.ValidBefore != CertTimeInfinity && (unixNow >= before || before < 0) {
		return fmt.Errorf("ssh: cert has expired")
	}
	if err := cert.SignatureKey.Verify(cert.bytesForSigning(), cert.Signature); err != nil {
		return fmt.Errorf("ssh: certificate signature does not verify")
	}

	return nil
}

// SignCert sets c.SignatureKey to the authority's public key and stores a
// Signature, by authority, in the certificate.
func (c *Certificate) SignCert(rand io.Reader, authority Signer) error {
	c.Nonce = make([]byte, 32)
	if _, err := io.ReadFull(rand, c.Nonce); err != nil {
		return err
	}
	c.SignatureKey = authority.PublicKey()

	sig, err := authority.Sign(rand, c.bytesForSigning())
	if err != nil {
		return err
	}
	c.Signature = sig
	return nil
}

var certAlgoNames = map[string]string{
	KeyAlgoRSA:      CertAlgoRSAv01,
	KeyAlgoDSA:      CertAlgoDSAv01,
	KeyAlgoECDSA256: CertAlgoECDSA256v01,
	KeyAlgoECDSA384: CertAlgoECDSA384v01,
	KeyAlgoECDSA521: CertAlgoECDSA521v01,
}

// certToPrivAlgo returns the underlying algorithm for a certificate algorithm.
// Panics if a non-certificate algorithm is passed.
func certToPrivAlgo(algo string) string {
	for privAlgo, pubAlgo := range certAlgoNames {
		if pubAlgo == algo {
			return privAlgo
		}
	}
	panic("unknown cert algorithm")
}

func (cert *Certificate) bytesForSigning() []byte {
	c2 := *cert
	c2.Signature = nil
	out := c2.Marshal()
	// Drop trailing signature length.
	return out[:len(out)-4]
}

// Marshal serializes c into OpenSSH's wire format. It is part of the
// PublicKey interface.
func (c *Certificate) Marshal() []byte {
	generic := genericCertData{
		Serial:          c.Serial,
		CertType:        c.CertType,
		KeyId:           c.KeyId,
		ValidPrincipals: marshalStringList(c.ValidPrincipals),
		ValidAfter:      uint64(c.ValidAfter),
		ValidBefore:     uint64(c.ValidBefore),
		CriticalOptions: marshalTuples(c.CriticalOptions),
		Extensions:      marshalTuples(c.Extensions),
		Reserved:        c.Reserved,
		SignatureKey:    c.SignatureKey.Marshal(),
	}
	if c.Signature != nil {
		generic.Signature = Marshal(c.Signature)
	}
	genericBytes := Marshal(&generic)
	keyBytes := c.Key.Marshal()
	_, keyBytes, _ = parseString(keyBytes)
	prefix := Marshal(&struct {
		Name  string
		Nonce []byte
		Key   []byte `ssh:"rest"`
	}{c.Type(), c.Nonce, keyBytes})

	result := make([]byte, 0, len(prefix)+len(genericBytes))
	result = append(result, prefix...)
	result = append(result, genericBytes...)
	return result
}

// Type returns the key name. It is part of the PublicKey interface.
func (c *Certificate) Type() string {
	algo, ok := certAlgoNames[c.Key.Type()]
	if !ok {
		panic("unknown cert key type")
	}
	return algo
}

// Verify verifies a signature against the certificate's public
// key. It is part of the PublicKey interface.
func (c *Certificate) Verify(data []byte, sig *Signature) error {
	return c.Key.Verify(data, sig)
}

func parseSignatureBody(in []byte) (out *Signature, rest []byte, ok bool) {
	format, in, ok := parseString(in)
	if !ok {
		return
	}

	out = &Signature{
		Format: string(format),
	}

	if out.Blob, in, ok = parseString(in); !ok {
		return
	}

	return out, in, ok
}

func parseSignature(in []byte) (out *Signature, rest []byte, ok bool) {
	sigBytes, rest, ok := parseString(in)
	if !ok {
		return
	}

	out, trailing, ok := parseSignatureBody(sigBytes)
	if !ok || len(trailing) > 0 {
		return nil, nil, false
	}
	return
}
