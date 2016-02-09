package trustmanager

import (
	"encoding/pem"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/docker/notary"
	"github.com/docker/notary/passphrase"
	"github.com/docker/notary/tuf/data"
)

const (
	rootKeysSubdir    = "root_keys"
	nonRootKeysSubdir = "tuf_keys"
	privDir           = "private"
)

type keyInfoMap map[string]KeyInfo

// KeyFileStore persists and manages private keys on disk
type KeyFileStore struct {
	sync.Mutex
	SimpleFileStore
	passphrase.Retriever
	cachedKeys map[string]*cachedKey
	keyInfoMap
}

// KeyMemoryStore manages private keys in memory
type KeyMemoryStore struct {
	sync.Mutex
	MemoryFileStore
	passphrase.Retriever
	cachedKeys map[string]*cachedKey
	keyInfoMap
}

// KeyInfo stores the role, path, and gun for a corresponding private key ID
// It is assumed that each private key ID is unique
type KeyInfo struct {
	Gun  string
	Role string
}

// NewKeyFileStore returns a new KeyFileStore creating a private directory to
// hold the keys.
func NewKeyFileStore(baseDir string, passphraseRetriever passphrase.Retriever) (*KeyFileStore, error) {
	baseDir = filepath.Join(baseDir, notary.PrivDir)
	fileStore, err := NewPrivateSimpleFileStore(baseDir, keyExtension)
	if err != nil {
		return nil, err
	}
	cachedKeys := make(map[string]*cachedKey)
	keyInfoMap := make(keyInfoMap)

	keyStore := &KeyFileStore{SimpleFileStore: *fileStore,
		Retriever:  passphraseRetriever,
		cachedKeys: cachedKeys,
		keyInfoMap: keyInfoMap,
	}

	// Load this keystore's ID --> gun/role map
	keyStore.loadKeyInfo()
	return keyStore, nil
}

func generateKeyInfoMap(s LimitedFileStore) map[string]KeyInfo {
	keyInfoMap := make(map[string]KeyInfo)
	for _, keyPath := range s.ListFiles() {
		// Remove the prefix of the directory from the filename for GUN/role/ID parsing
		var keyIDAndGun, keyRole string
		if strings.HasPrefix(keyPath, rootKeysSubdir+"/") {
			keyIDAndGun = strings.TrimPrefix(keyPath, rootKeysSubdir+"/")
		} else {
			keyIDAndGun = strings.TrimPrefix(keyPath, nonRootKeysSubdir+"/")
		}

		// Separate the ID and GUN (can be empty) from the filepath
		keyIDAndGun = strings.TrimSpace(keyIDAndGun)
		keyGun := getGunFromFullID(keyIDAndGun)
		keyID := filepath.Base(keyIDAndGun)

		// If the key does not have a _, we'll attempt to
		// read it as a PEM
		underscoreIndex := strings.LastIndex(keyID, "_")
		if underscoreIndex == -1 {
			d, err := s.Get(keyPath)
			if err != nil {
				logrus.Error(err)
				continue
			}
			block, _ := pem.Decode(d)
			if block == nil {
				continue
			}
			if role, ok := block.Headers["role"]; ok {
				keyRole = role
			}
		} else {
			// This is the legacy KEYID_ROLE filename
			// The keyID is the first part of the keyname
			// The keyRole is the second part of the keyname
			// in a key named abcde_root, abcde is the keyID and root is the KeyAlias
			legacyID := keyID
			keyID = legacyID[:underscoreIndex]
			keyRole = legacyID[underscoreIndex+1:]
		}
		keyInfoMap[keyID] = KeyInfo{Gun: keyGun, Role: keyRole}
	}
	return keyInfoMap
}

func getGunFromFullID(fullKeyID string) string {
	keyGun := filepath.Dir(fullKeyID)
	// If the gun is empty, Dir will return .
	if keyGun == "." {
		keyGun = ""
	}
	return keyGun
}

func (s *KeyFileStore) loadKeyInfo() {
	s.keyInfoMap = generateKeyInfoMap(s)
}

func (s *KeyMemoryStore) loadKeyInfo() {
	s.keyInfoMap = generateKeyInfoMap(s)
}

func (s *KeyFileStore) GetKeyInfo(keyID string) (KeyInfo, error) {
	if info, ok := s.keyInfoMap[keyID]; ok {
		return info, nil
	}
	return KeyInfo{}, fmt.Errorf("Could not find info for keyID %s", keyID)
}

func (s *KeyMemoryStore) GetKeyInfo(keyID string) (KeyInfo, error) {
	if info, ok := s.keyInfoMap[keyID]; ok {
		return info, nil
	}
	return KeyInfo{}, fmt.Errorf("Could not find info for keyID %s", keyID)
}

// Name returns a user friendly name for the location this store
// keeps its data
func (s *KeyFileStore) Name() string {
	return fmt.Sprintf("file (%s)", s.SimpleFileStore.BaseDir())
}

// AddKey stores the contents of a PEM-encoded private key as a PEM block
func (s *KeyFileStore) AddKey(name, role string, privKey data.PrivateKey) error {
	s.Lock()
	defer s.Unlock()
	err := addKey(s, s.Retriever, s.cachedKeys, name, role, privKey)
	if err != nil {
		return err
	}
	s.keyInfoMap[privKey.ID()] = KeyInfo{Gun: getGunFromFullID(name), Role: role}
	return nil
}

// GetKey returns the PrivateKey given a KeyID
func (s *KeyFileStore) GetKey(name string) (data.PrivateKey, string, error) {
	s.Lock()
	defer s.Unlock()
	// If this is a bare key ID without the gun, prepend the gun so the filestore lookup succeeds
	if keyInfo, ok := s.keyInfoMap[name]; ok {
		name = filepath.Join(keyInfo.Gun, name)
	}
	return getKey(s, s.Retriever, s.cachedKeys, name)
}

// ListKeys returns a list of unique PublicKeys present on the KeyFileStore.
func (s *KeyFileStore) ListKeys() map[string]KeyInfo {
	return s.keyInfoMap
}

// RemoveKey removes the key from the keyfilestore
func (s *KeyFileStore) RemoveKey(name string) error {
	s.Lock()
	defer s.Unlock()
	// If this is a bare key ID without the gun, prepend the gun so the filestore lookup succeeds
	if keyInfo, ok := s.keyInfoMap[name]; ok {
		name = filepath.Join(keyInfo.Gun, name)
	}
	err := removeKey(s, s.cachedKeys, name)
	if err != nil {
		return err
	}
	// Remove this key from our keyInfo map if we removed from our filesystem
	if _, ok := s.keyInfoMap[name]; ok {
		delete(s.keyInfoMap, name)
	} else {
		// This might be of the form GUN/ID  - try to delete without the gun
		delete(s.keyInfoMap, filepath.Base(name))
	}
	return nil
}

// ExportKey exports the encrypted bytes from the keystore and writes it to
// dest.
func (s *KeyFileStore) ExportKey(name string) ([]byte, error) {
	keyBytes, _, err := getRawKey(s, name)
	if err != nil {
		return nil, err
	}
	return keyBytes, nil
}

// ImportKey imports the private key in the encrypted bytes into the keystore
// with the given key ID and alias.
// This is only used for root, so no need to touch the keyInfoMap
func (s *KeyFileStore) ImportKey(pemBytes []byte, alias string) error {
	err := importKey(s, s.Retriever, s.cachedKeys, alias, pemBytes)
	if err != nil {
		return err
	}
	keyID, keyInfo := getImportedKeyInfo(pemBytes, alias)
	s.keyInfoMap[keyID] = keyInfo
	return nil
}

// NewKeyMemoryStore returns a new KeyMemoryStore which holds keys in memory
func NewKeyMemoryStore(passphraseRetriever passphrase.Retriever) *KeyMemoryStore {
	memStore := NewMemoryFileStore()
	cachedKeys := make(map[string]*cachedKey)

	keyInfoMap := make(keyInfoMap)

	keyStore := &KeyMemoryStore{MemoryFileStore: *memStore,
		Retriever:  passphraseRetriever,
		cachedKeys: cachedKeys,
		keyInfoMap: keyInfoMap,
	}

	// Load this keystore's ID --> gun/role map
	keyStore.loadKeyInfo()
	return keyStore
}

// Name returns a user friendly name for the location this store
// keeps its data
func (s *KeyMemoryStore) Name() string {
	return "memory"
}

// AddKey stores the contents of a PEM-encoded private key as a PEM block
func (s *KeyMemoryStore) AddKey(name, alias string, privKey data.PrivateKey) error {
	s.Lock()
	defer s.Unlock()
	err := addKey(s, s.Retriever, s.cachedKeys, name, alias, privKey)
	if err != nil {
		return err
	}
	s.keyInfoMap[privKey.ID()] = KeyInfo{Gun: getGunFromFullID(name), Role: alias}
	return nil
}

// GetKey returns the PrivateKey given a KeyID
func (s *KeyMemoryStore) GetKey(name string) (data.PrivateKey, string, error) {
	s.Lock()
	defer s.Unlock()
	// If this is a bare key ID without the gun, prepend the gun so the filestore lookup succeeds
	if keyInfo, ok := s.keyInfoMap[name]; ok {
		name = filepath.Join(keyInfo.Gun, name)
	}
	return getKey(s, s.Retriever, s.cachedKeys, name)
}

// ListKeys returns a list of unique PublicKeys present on the KeyFileStore.
func (s *KeyMemoryStore) ListKeys() map[string]KeyInfo {
	return s.keyInfoMap
}

// RemoveKey removes the key from the keystore
func (s *KeyMemoryStore) RemoveKey(name string) error {
	s.Lock()
	defer s.Unlock()
	// If this is a bare key ID without the gun, prepend the gun so the filestore lookup succeeds
	if keyInfo, ok := s.keyInfoMap[name]; ok {
		name = filepath.Join(keyInfo.Gun, name)
	}
	err := removeKey(s, s.cachedKeys, name)
	if err != nil {
		return err
	}
	// Remove this key from our keyInfo map if we removed from our filesystem
	if _, ok := s.keyInfoMap[name]; ok {
		delete(s.keyInfoMap, name)
	} else {
		// This might be of the form GUN/ID  - try to delete without the gun
		delete(s.keyInfoMap, filepath.Base(name))
	}
	return nil
}

// ExportKey exports the encrypted bytes from the keystore and writes it to
// dest.
func (s *KeyMemoryStore) ExportKey(name string) ([]byte, error) {
	keyBytes, _, err := getRawKey(s, name)
	if err != nil {
		return nil, err
	}
	return keyBytes, nil
}

// ImportKey imports the private key in the encrypted bytes into the keystore
// with the given key ID and alias.
func (s *KeyMemoryStore) ImportKey(pemBytes []byte, alias string) error {
	err := importKey(s, s.Retriever, s.cachedKeys, alias, pemBytes)
	if err != nil {
		return err
	}
	keyID, keyInfo := getImportedKeyInfo(pemBytes, alias)
	s.keyInfoMap[keyID] = keyInfo
	return nil
}

func getImportedKeyInfo(pemBytes []byte, filename string) (string, KeyInfo) {
	var keyID, gun, role string
	// Try to read the role and gun from the filepath and PEM headers
	if filepath.HasPrefix(filename, "root_keys") {
		role = data.CanonicalRootRole
		gun = ""
		filename = strings.TrimPrefix(filename, "root_keys")
	} else {
		filename = strings.TrimPrefix(filename, "tuf_keys")
		gun = getGunFromFullID(filename)
	}
	keyIDFull := filepath.Base(filename)
	underscoreIndex := strings.LastIndex(keyIDFull, "_")
	if underscoreIndex == -1 {
		block, _ := pem.Decode(pemBytes)
		if keyRole, ok := block.Headers["role"]; ok {
			role = keyRole
		}
		keyID = filepath.Base(keyIDFull)
	} else {
		// This is the legacy KEYID_ROLE filename
		// The keyID is the first part of the keyname
		// The keyRole is the second part of the keyname
		// in a key named abcde_root, abcde is the keyID and root is the KeyAlias
		legacyID := keyIDFull
		keyID = legacyID[:underscoreIndex]
		role = legacyID[underscoreIndex+1:]
	}
	return keyID, KeyInfo{Gun: gun, Role: role}
}

func addKey(s LimitedFileStore, passphraseRetriever passphrase.Retriever, cachedKeys map[string]*cachedKey, name, role string, privKey data.PrivateKey) error {

	var (
		chosenPassphrase string
		giveup           bool
		err              error
	)

	for attempts := 0; ; attempts++ {
		chosenPassphrase, giveup, err = passphraseRetriever(name, role, true, attempts)
		if err != nil {
			continue
		}
		if giveup {
			return ErrAttemptsExceeded{}
		}
		if attempts > 10 {
			return ErrAttemptsExceeded{}
		}
		break
	}

	return encryptAndAddKey(s, chosenPassphrase, cachedKeys, name, role, privKey)
}

// getKeyRole finds the role for the given keyID. It attempts to look
// both in the newer format PEM headers, and also in the legacy filename
// format. It returns: the role, whether it was found in the legacy format
// (true == legacy), and an error
func getKeyRole(s LimitedFileStore, keyID string) (string, bool, error) {
	name := strings.TrimSpace(strings.TrimSuffix(filepath.Base(keyID), filepath.Ext(keyID)))

	for _, file := range s.ListFiles() {
		filename := filepath.Base(file)

		if strings.HasPrefix(filename, name) {
			d, err := s.Get(file)
			if err != nil {
				return "", false, err
			}
			block, _ := pem.Decode(d)
			if block != nil {
				if role, ok := block.Headers["role"]; ok {
					return role, false, nil
				}
			}

			role := strings.TrimPrefix(filename, name+"_")
			return role, true, nil
		}
	}

	return "", false, &ErrKeyNotFound{KeyID: keyID}
}

// GetKey returns the PrivateKey given a KeyID
func getKey(s LimitedFileStore, passphraseRetriever passphrase.Retriever, cachedKeys map[string]*cachedKey, name string) (data.PrivateKey, string, error) {
	cachedKeyEntry, ok := cachedKeys[name]
	if ok {
		return cachedKeyEntry.key, cachedKeyEntry.alias, nil
	}

	keyBytes, keyAlias, err := getRawKey(s, name)
	if err != nil {
		return nil, "", err
	}

	// See if the key is encrypted. If its encrypted we'll fail to parse the private key
	privKey, err := ParsePEMPrivateKey(keyBytes, "")
	if err != nil {
		privKey, _, err = GetPasswdDecryptBytes(passphraseRetriever, keyBytes, name, string(keyAlias))
		if err != nil {
			return nil, "", err
		}
	}
	cachedKeys[name] = &cachedKey{alias: keyAlias, key: privKey}
	return privKey, keyAlias, nil
}

// RemoveKey removes the key from the keyfilestore
func removeKey(s LimitedFileStore, cachedKeys map[string]*cachedKey, name string) error {
	role, legacy, err := getKeyRole(s, name)
	if err != nil {
		return err
	}

	delete(cachedKeys, name)

	if legacy {
		name = name + "_" + role
	}

	// being in a subdirectory is for backwards compatibliity
	err = s.Remove(filepath.Join(getSubdir(role), name))
	if err != nil {
		return err
	}
	return nil
}

// Assumes 2 subdirectories, 1 containing root keys and 1 containing tuf keys
func getSubdir(alias string) string {
	if alias == "root" {
		return notary.RootKeysSubdir
	}
	return notary.NonRootKeysSubdir
}

// Given a key ID, gets the bytes and alias belonging to that key if the key
// exists
func getRawKey(s LimitedFileStore, name string) ([]byte, string, error) {
	role, legacy, err := getKeyRole(s, name)
	if err != nil {
		return nil, "", err
	}

	if legacy {
		name = name + "_" + role
	}

	var keyBytes []byte
	keyBytes, err = s.Get(filepath.Join(getSubdir(role), name))
	if err != nil {
		return nil, "", err
	}
	return keyBytes, role, nil
}

// GetPasswdDecryptBytes gets the password to decrypt the given pem bytes.
// Returns the password and private key
func GetPasswdDecryptBytes(passphraseRetriever passphrase.Retriever, pemBytes []byte, name, alias string) (data.PrivateKey, string, error) {
	var (
		passwd  string
		retErr  error
		privKey data.PrivateKey
	)
	for attempts := 0; ; attempts++ {
		var (
			giveup bool
			err    error
		)
		passwd, giveup, err = passphraseRetriever(name, alias, false, attempts)
		// Check if the passphrase retriever got an error or if it is telling us to give up
		if giveup || err != nil {
			return nil, "", ErrPasswordInvalid{}
		}
		if attempts > 10 {
			return nil, "", ErrAttemptsExceeded{}
		}

		// Try to convert PEM encoded bytes back to a PrivateKey using the passphrase
		privKey, err = ParsePEMPrivateKey(pemBytes, passwd)
		if err != nil {
			retErr = ErrPasswordInvalid{}
		} else {
			// We managed to parse the PrivateKey. We've succeeded!
			retErr = nil
			break
		}
	}
	if retErr != nil {
		return nil, "", retErr
	}
	return privKey, passwd, nil
}

func encryptAndAddKey(s LimitedFileStore, passwd string, cachedKeys map[string]*cachedKey, name, role string, privKey data.PrivateKey) error {

	var (
		pemPrivKey []byte
		err        error
	)

	if passwd != "" {
		pemPrivKey, err = EncryptPrivateKey(privKey, role, passwd)
	} else {
		pemPrivKey, err = KeyToPEM(privKey, role)
	}

	if err != nil {
		return err
	}

	cachedKeys[name] = &cachedKey{alias: role, key: privKey}
	return s.Add(filepath.Join(getSubdir(role), name), pemPrivKey)
}

func importKey(s LimitedFileStore, passphraseRetriever passphrase.Retriever, cachedKeys map[string]*cachedKey, alias string, pemBytes []byte) error {

	if alias != data.CanonicalRootRole {
		return s.Add(alias, pemBytes)
	}

	privKey, passphrase, err := GetPasswdDecryptBytes(
		passphraseRetriever, pemBytes, "", "imported "+alias)

	if err != nil {
		return err
	}

	var name string
	name = privKey.ID()
	return encryptAndAddKey(s, passphrase, cachedKeys, name, alias, privKey)
}
