package cryptoservice

import (
	"crypto/rand"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/signed"
)

var algoToSigType = map[string]data.SigAlgorithm{
	data.ECDSAKey:   data.ECDSASignature,
	data.ED25519Key: data.EDDSASignature,
	data.RSAKey:     data.RSAPSSSignature,
}

var passphraseRetriever = func(string, string, bool, int) (string, bool, error) { return "", false, nil }

type CryptoServiceTester struct {
	cryptoServiceFactory func() *CryptoService
	role                 string
	keyAlgo              string
}

// asserts that created key exists
func (c CryptoServiceTester) TestCreateAndGetKey(t *testing.T) {
	cryptoService := c.cryptoServiceFactory()

	// Test Create
	tufKey, err := cryptoService.Create(c.role, c.keyAlgo)
	assert.NoError(t, err, c.errorMsg("error creating key"))

	// Test GetKey
	retrievedKey := cryptoService.GetKey(tufKey.ID())
	assert.NotNil(t, retrievedKey,
		c.errorMsg("Could not find key ID %s", tufKey.ID()))
	assert.Equal(t, tufKey.Public(), retrievedKey.Public(),
		c.errorMsg("retrieved public key didn't match"))

	// Test GetPrivateKey
	retrievedKey, alias, err := cryptoService.GetPrivateKey(tufKey.ID())
	assert.NoError(t, err)
	assert.Equal(t, tufKey.ID(), retrievedKey.ID(),
		c.errorMsg("retrieved private key didn't have the right ID"))
	assert.Equal(t, c.role, alias)
}

// If there are multiple keystores, ensure that a key is only added to one -
// the first in the list of keyStores (which is in order of preference)
func (c CryptoServiceTester) TestCreateAndGetWhenMultipleKeystores(t *testing.T) {
	cryptoService := c.cryptoServiceFactory()
	cryptoService.keyStores = append(cryptoService.keyStores,
		trustmanager.NewKeyMemoryStore(passphraseRetriever))

	// Test Create
	tufKey, err := cryptoService.Create(c.role, c.keyAlgo)
	assert.NoError(t, err, c.errorMsg("error creating key"))

	// Only the first keystore should have the key
	keyPath := tufKey.ID()
	if c.role != data.CanonicalRootRole && cryptoService.gun != "" {
		keyPath = filepath.Join(cryptoService.gun, keyPath)
	}
	_, _, err = cryptoService.keyStores[0].GetKey(keyPath)
	assert.NoError(t, err, c.errorMsg(
		"First keystore does not have the key %s", keyPath))
	_, _, err = cryptoService.keyStores[1].GetKey(keyPath)
	assert.Error(t, err, c.errorMsg(
		"Second keystore has the key %s", keyPath))

	// GetKey works across multiple keystores
	retrievedKey := cryptoService.GetKey(tufKey.ID())
	assert.NotNil(t, retrievedKey,
		c.errorMsg("Could not find key ID %s", tufKey.ID()))
}

// asserts that getting key fails for a non-existent key
func (c CryptoServiceTester) TestGetNonexistentKey(t *testing.T) {
	cryptoService := c.cryptoServiceFactory()

	assert.Nil(t, cryptoService.GetKey("boguskeyid"),
		c.errorMsg("non-nil result for bogus keyid"))

	_, _, err := cryptoService.GetPrivateKey("boguskeyid")
	assert.NotNil(t, err)
}

// asserts that signing with a created key creates a valid signature
func (c CryptoServiceTester) TestSignWithKey(t *testing.T) {
	cryptoService := c.cryptoServiceFactory()
	content := []byte("this is a secret")

	tufKey, err := cryptoService.Create(c.role, c.keyAlgo)
	assert.NoError(t, err, c.errorMsg("error creating key"))

	// Test Sign
	privKey, role, err := cryptoService.GetPrivateKey(tufKey.ID())
	assert.NoError(t, err, c.errorMsg("failed to get private key"))
	assert.Equal(t, c.role, role)

	signature, err := privKey.Sign(rand.Reader, content, nil)
	assert.NoError(t, err, c.errorMsg("signing failed"))

	verifier, ok := signed.Verifiers[algoToSigType[c.keyAlgo]]
	assert.True(t, ok, c.errorMsg("Unknown verifier for algorithm"))

	err = verifier.Verify(tufKey, signature, content)
	assert.NoError(t, err,
		c.errorMsg("verification failed for %s key type", c.keyAlgo))
}

// asserts that signing, if there are no matching keys, produces no signatures
func (c CryptoServiceTester) TestSignNoMatchingKeys(t *testing.T) {
	cryptoService := c.cryptoServiceFactory()

	privKey, err := trustmanager.GenerateECDSAKey(rand.Reader)
	assert.NoError(t, err, c.errorMsg("error creating key"))

	// Test Sign
	_, _, err = cryptoService.GetPrivateKey(privKey.ID())
	assert.Error(t, err, c.errorMsg("Should not have found private key"))
}

// Test GetPrivateKey succeeds when multiple keystores have the same key
func (c CryptoServiceTester) TestGetPrivateKeyMultipleKeystores(t *testing.T) {
	cryptoService := c.cryptoServiceFactory()
	cryptoService.keyStores = append(cryptoService.keyStores,
		trustmanager.NewKeyMemoryStore(passphraseRetriever))

	privKey, err := trustmanager.GenerateECDSAKey(rand.Reader)
	assert.NoError(t, err, c.errorMsg("error creating key"))

	for _, store := range cryptoService.keyStores {
		err := store.AddKey(privKey.ID(), c.role, privKey)
		assert.NoError(t, err)
	}

	foundKey, role, err := cryptoService.GetPrivateKey(privKey.ID())
	assert.NoError(t, err, c.errorMsg("failed to get private key"))
	assert.Equal(t, c.role, role)
	assert.Equal(t, privKey.ID(), foundKey.ID())
}

// asserts that removing key that exists succeeds
func (c CryptoServiceTester) TestRemoveCreatedKey(t *testing.T) {
	cryptoService := c.cryptoServiceFactory()

	tufKey, err := cryptoService.Create(c.role, c.keyAlgo)
	assert.NoError(t, err, c.errorMsg("error creating key"))
	assert.NotNil(t, cryptoService.GetKey(tufKey.ID()))

	// Test RemoveKey
	err = cryptoService.RemoveKey(tufKey.ID())
	assert.NoError(t, err, c.errorMsg("could not remove key"))
	retrievedKey := cryptoService.GetKey(tufKey.ID())
	assert.Nil(t, retrievedKey, c.errorMsg("remove didn't work"))
}

// asserts that removing key will remove it from all keystores
func (c CryptoServiceTester) TestRemoveFromMultipleKeystores(t *testing.T) {
	cryptoService := c.cryptoServiceFactory()
	cryptoService.keyStores = append(cryptoService.keyStores,
		trustmanager.NewKeyMemoryStore(passphraseRetriever))

	privKey, err := trustmanager.GenerateECDSAKey(rand.Reader)
	assert.NoError(t, err, c.errorMsg("error creating key"))

	for _, store := range cryptoService.keyStores {
		err := store.AddKey(privKey.ID(), "root", privKey)
		assert.NoError(t, err)
	}

	assert.NotNil(t, cryptoService.GetKey(privKey.ID()))

	// Remove removes it from all key stores
	err = cryptoService.RemoveKey(privKey.ID())
	assert.NoError(t, err, c.errorMsg("could not remove key"))

	for _, store := range cryptoService.keyStores {
		_, _, err := store.GetKey(privKey.ID())
		assert.Error(t, err)
	}
}

// asserts that listing keys works with multiple keystores, and that the
// same keys are deduplicated
func (c CryptoServiceTester) TestListFromMultipleKeystores(t *testing.T) {
	cryptoService := c.cryptoServiceFactory()
	cryptoService.keyStores = append(cryptoService.keyStores,
		trustmanager.NewKeyMemoryStore(passphraseRetriever))

	expectedKeysIDs := make(map[string]bool) // just want to be able to index by key

	for i := 0; i < 3; i++ {
		privKey, err := trustmanager.GenerateECDSAKey(rand.Reader)
		assert.NoError(t, err, c.errorMsg("error creating key"))
		expectedKeysIDs[privKey.ID()] = true

		// adds one different key to each keystore, and then one key to
		// both keystores
		for j, store := range cryptoService.keyStores {
			if i == j || i == 2 {
				store.AddKey(privKey.ID(), "root", privKey)
			}
		}
	}
	// sanity check - each should have 2
	for _, store := range cryptoService.keyStores {
		assert.Len(t, store.ListKeys(), 2, c.errorMsg("added keys wrong"))
	}

	keyList := cryptoService.ListKeys("root")
	assert.Len(t, keyList, 4,
		c.errorMsg(
			"ListKeys should have 4 keys (not necesarily unique) but does not: %v", keyList))
	for _, k := range keyList {
		_, ok := expectedKeysIDs[k]
		assert.True(t, ok, c.errorMsg("Unexpected key %s", k))
	}

	keyMap := cryptoService.ListAllKeys()
	assert.Len(t, keyMap, 3,
		c.errorMsg("ListAllKeys should have 3 unique keys but does not: %v", keyMap))

	for k, role := range keyMap {
		_, ok := expectedKeysIDs[k]
		assert.True(t, ok)
		assert.Equal(t, "root", role)
	}
}

// Prints out an error message with information about the key algorithm,
// role, and test name. Ideally we could generate different tests given
// data, without having to put for loops in one giant test function, but
// that involves a lot of boilerplate.  So as a compromise, everything will
// still be run in for loops in one giant test function, but we can at
// least provide an error message stating what data/helper test function
// failed.
func (c CryptoServiceTester) errorMsg(message string, args ...interface{}) string {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)    // the caller of errorMsg
	f := runtime.FuncForPC(pc[0])
	return fmt.Sprintf("%s (role: %s, keyAlgo: %s): %s", f.Name(), c.role,
		c.keyAlgo, fmt.Sprintf(message, args...))
}

func testCryptoService(t *testing.T, gun string) {
	getTestingCryptoService := func() *CryptoService {
		return NewCryptoService(
			gun, trustmanager.NewKeyMemoryStore(passphraseRetriever))
	}
	roles := []string{
		data.CanonicalRootRole,
		data.CanonicalTargetsRole,
		data.CanonicalSnapshotRole,
		data.CanonicalTimestampRole,
	}

	for _, role := range roles {
		for algo := range algoToSigType {
			cst := CryptoServiceTester{
				cryptoServiceFactory: getTestingCryptoService,
				role:                 role,
				keyAlgo:              algo,
			}
			cst.TestCreateAndGetKey(t)
			cst.TestCreateAndGetWhenMultipleKeystores(t)
			cst.TestGetNonexistentKey(t)
			cst.TestSignWithKey(t)
			cst.TestSignNoMatchingKeys(t)
			cst.TestGetPrivateKeyMultipleKeystores(t)
			cst.TestRemoveCreatedKey(t)
			cst.TestRemoveFromMultipleKeystores(t)
			cst.TestListFromMultipleKeystores(t)
		}
	}
}

func TestCryptoServiceWithNonEmptyGUN(t *testing.T) {
	testCryptoService(t, "org/repo")
}

func TestCryptoServiceWithEmptyGUN(t *testing.T) {
	testCryptoService(t, "")
}
