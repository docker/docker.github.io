// +build pkcs11

package yubikey

import (
	"crypto/rand"
	"fmt"
	"reflect"
	"testing"

	"github.com/docker/notary/passphrase"
	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/tuf/data"
	"github.com/miekg/pkcs11"
	"github.com/stretchr/testify/assert"
)

var ret = passphrase.ConstantRetriever("passphrase")

// create a new store for clearing out keys, because we don't want to pollute
// any cache
func clearAllKeys(t *testing.T) {
	store, err := NewYubiKeyStore(trustmanager.NewKeyMemoryStore(ret), ret)
	assert.NoError(t, err)

	for k := range store.ListKeys() {
		err := store.RemoveKey(k)
		assert.NoError(t, err)
	}
}

func TestEnsurePrivateKeySizePassesThroughRightSizeArrays(t *testing.T) {
	fullByteArray := make([]byte, ecdsaPrivateKeySize)
	for i := range fullByteArray {
		fullByteArray[i] = byte(1)
	}

	result := ensurePrivateKeySize(fullByteArray)
	assert.True(t, reflect.DeepEqual(fullByteArray, result))
}

// The pad32Byte helper function left zero-pads byte arrays that are less than
// ecdsaPrivateKeySize bytes
func TestEnsurePrivateKeySizePadsLessThanRequiredSizeArrays(t *testing.T) {
	shortByteArray := make([]byte, ecdsaPrivateKeySize/2)
	for i := range shortByteArray {
		shortByteArray[i] = byte(1)
	}

	expected := append(
		make([]byte, ecdsaPrivateKeySize-ecdsaPrivateKeySize/2),
		shortByteArray...)

	result := ensurePrivateKeySize(shortByteArray)
	assert.True(t, reflect.DeepEqual(expected, result))
}

func testAddKey(t *testing.T, store trustmanager.KeyStore) (data.PrivateKey, error) {
	privKey, err := trustmanager.GenerateECDSAKey(rand.Reader)
	assert.NoError(t, err)

	err = store.AddKey(privKey.ID(), data.CanonicalRootRole, privKey)
	return privKey, err
}

func addMaxKeys(t *testing.T, store trustmanager.KeyStore) []string {
	keys := make([]string, 0, numSlots)
	// create the maximum number of keys
	for i := 0; i < numSlots; i++ {
		privKey, err := testAddKey(t, store)
		assert.NoError(t, err)
		keys = append(keys, privKey.ID())
	}
	return keys
}

// We can add keys enough times to fill up all the slots in the Yubikey.
// They are backed up, and we can then list them and get the keys.
func TestYubiAddKeysAndRetrieve(t *testing.T) {
	if !YubikeyAccessible() {
		t.Skip("Must have Yubikey access.")
	}
	clearAllKeys(t)

	SetYubikeyKeyMode(KeymodeNone)
	defer func() {
		SetYubikeyKeyMode(KeymodeTouch | KeymodePinOnce)
	}()

	// create 4 keys on the original store
	backup := trustmanager.NewKeyMemoryStore(ret)
	origStore, err := NewYubiKeyStore(backup, ret)
	assert.NoError(t, err)
	keys := addMaxKeys(t, origStore)

	// create a new store, since we want to be sure the original store's cache
	// is not masking any issues
	cleanStore, err := NewYubiKeyStore(trustmanager.NewKeyMemoryStore(ret), ret)
	assert.NoError(t, err)

	// All 4 keys should be in the original store, in the clean store (which
	// makes sure the keys are actually on the Yubikey and not on the original
	// store's cache, and on the backup store)
	for _, store := range []trustmanager.KeyStore{origStore, cleanStore, backup} {
		listedKeys := store.ListKeys()
		assert.Len(t, listedKeys, numSlots)
		for _, k := range keys {
			r, ok := listedKeys[k]
			assert.True(t, ok)
			assert.Equal(t, data.CanonicalRootRole, r)

			_, _, err := store.GetKey(k)
			assert.NoError(t, err)
		}
	}
}

// We can't add a key if there are no more slots
func TestYubiAddKeyFailureIfNoMoreSlots(t *testing.T) {
	if !YubikeyAccessible() {
		t.Skip("Must have Yubikey access.")
	}
	clearAllKeys(t)

	SetYubikeyKeyMode(KeymodeNone)
	defer func() {
		SetYubikeyKeyMode(KeymodeTouch | KeymodePinOnce)
	}()

	// create 4 keys on the original store
	backup := trustmanager.NewKeyMemoryStore(ret)
	origStore, err := NewYubiKeyStore(backup, ret)
	assert.NoError(t, err)
	addMaxKeys(t, origStore)

	// add another key - should fail because there are no more slots
	badKey, err := testAddKey(t, origStore)
	assert.Error(t, err)

	// create a new store, since we want to be sure the original store's cache
	// is not masking any issues
	cleanStore, err := NewYubiKeyStore(trustmanager.NewKeyMemoryStore(ret), ret)
	assert.NoError(t, err)

	// The key should not be in the original store, in the new clean store, or
	// in teh backup store.
	for _, store := range []trustmanager.KeyStore{origStore, cleanStore, backup} {
		// the key that wasn't created should not appear in ListKeys or GetKey
		_, _, err := store.GetKey(badKey.ID())
		assert.Error(t, err)
		for k := range store.ListKeys() {
			assert.NotEqual(t, badKey, k)
		}
	}
}

// If some random key in the middle was removed, adding a key will work (keys
// do not have to be deleted/added in order)
func TestYubiAddKeyCanAddToMiddleSlot(t *testing.T) {
	if !YubikeyAccessible() {
		t.Skip("Must have Yubikey access.")
	}
	clearAllKeys(t)

	SetYubikeyKeyMode(KeymodeNone)
	defer func() {
		SetYubikeyKeyMode(KeymodeTouch | KeymodePinOnce)
	}()

	// create 4 keys on the original store
	backup := trustmanager.NewKeyMemoryStore(ret)
	origStore, err := NewYubiKeyStore(backup, ret)
	assert.NoError(t, err)
	keys := addMaxKeys(t, origStore)

	// delete one of the middle keys, and assert we can still create a new key
	keyIDToDelete := keys[numSlots/2]
	err = origStore.RemoveKey(keyIDToDelete)
	assert.NoError(t, err)

	newKey, err := testAddKey(t, origStore)
	assert.NoError(t, err)

	// create a new store, since we want to be sure the original store's cache
	// is not masking any issues
	cleanStore, err := NewYubiKeyStore(trustmanager.NewKeyMemoryStore(ret), ret)
	assert.NoError(t, err)

	// The new key should be in the original store, in the new clean store, and
	// in the backup store.  The old key should not be in the original store,
	// or the new clean store.
	for _, store := range []trustmanager.KeyStore{origStore, cleanStore, backup} {
		// new key should appear in all stores
		gottenKey, _, err := store.GetKey(newKey.ID())
		assert.NoError(t, err)
		assert.Equal(t, gottenKey.ID(), newKey.ID())

		listedKeys := store.ListKeys()
		_, ok := listedKeys[newKey.ID()]
		assert.True(t, ok)

		// old key should not be in the non-backup stores
		if store != backup {
			_, _, err := store.GetKey(keyIDToDelete)
			assert.Error(t, err)
			_, ok = listedKeys[keyIDToDelete]
			assert.False(t, ok)
		}
	}
}

// RemoveKey removes a key from the yubikey, but not from the backup store.
func TestYubiRemoveKey(t *testing.T) {
	if !YubikeyAccessible() {
		t.Skip("Must have Yubikey access.")
	}
	clearAllKeys(t)

	SetYubikeyKeyMode(KeymodeNone)
	defer func() {
		SetYubikeyKeyMode(KeymodeTouch | KeymodePinOnce)
	}()

	backup := trustmanager.NewKeyMemoryStore(ret)
	origStore, err := NewYubiKeyStore(backup, ret)
	assert.NoError(t, err)

	key, err := testAddKey(t, origStore)
	assert.NoError(t, err)
	err = origStore.RemoveKey(key.ID())
	assert.NoError(t, err)

	// key remains in the backup store
	backupKey, role, err := backup.GetKey(key.ID())
	assert.NoError(t, err)
	assert.Equal(t, data.CanonicalRootRole, role)
	assert.Equal(t, key.ID(), backupKey.ID())

	// create a new store, since we want to be sure the original store's cache
	// is not masking any issues
	cleanStore, err := NewYubiKeyStore(trustmanager.NewKeyMemoryStore(ret), ret)
	assert.NoError(t, err)

	// key is not in either the original store or the clean store
	for _, store := range []*YubiKeyStore{origStore, cleanStore} {
		_, _, err := store.GetKey(key.ID())
		assert.Error(t, err)
	}
}

// ImportKey imports a key as root without adding it to the backup store
func TestYubiImportNewKey(t *testing.T) {
	if !YubikeyAccessible() {
		t.Skip("Must have Yubikey access.")
	}
	clearAllKeys(t)

	SetYubikeyKeyMode(KeymodeNone)
	defer func() {
		SetYubikeyKeyMode(KeymodeTouch | KeymodePinOnce)
	}()

	backup := trustmanager.NewKeyMemoryStore(ret)
	origStore, err := NewYubiKeyStore(backup, ret)
	assert.NoError(t, err)

	// generate key and import it
	privKey, err := trustmanager.GenerateECDSAKey(rand.Reader)
	assert.NoError(t, err)

	pemBytes, err := trustmanager.EncryptPrivateKey(privKey, "passphrase")
	assert.NoError(t, err)

	err = origStore.ImportKey(pemBytes, "root")
	assert.NoError(t, err)

	// key is not in backup store
	_, _, err = backup.GetKey(privKey.ID())
	assert.Error(t, err)

	// create a new store, since we want to be sure the original store's cache
	// is not masking any issues
	cleanStore, err := NewYubiKeyStore(trustmanager.NewKeyMemoryStore(ret), ret)
	assert.NoError(t, err)
	for _, store := range []*YubiKeyStore{origStore, cleanStore} {
		gottenKey, role, err := store.GetKey(privKey.ID())
		assert.NoError(t, err)
		assert.Equal(t, data.CanonicalRootRole, role)
		assert.Equal(t, privKey.Public(), gottenKey.Public())
	}
}

// Importing an existing key succeeds, but doesn't actually add the key, nor
// does it write it to backup.
func TestYubiImportExistingKey(t *testing.T) {
	if !YubikeyAccessible() {
		t.Skip("Must have Yubikey access.")
	}
	clearAllKeys(t)

	SetYubikeyKeyMode(KeymodeNone)
	defer func() {
		SetYubikeyKeyMode(KeymodeTouch | KeymodePinOnce)
	}()

	origStore, err := NewYubiKeyStore(trustmanager.NewKeyMemoryStore(ret), ret)
	assert.NoError(t, err)
	key, err := testAddKey(t, origStore)

	backup := trustmanager.NewKeyMemoryStore(ret)
	newStore, err := NewYubiKeyStore(backup, ret)
	assert.NoError(t, err)

	// for sanity, ensure that the key is already in the Yubikey
	k, _, err := newStore.GetKey(key.ID())
	assert.NoError(t, err)
	assert.NotNil(t, k)

	// import the key, which should have already been added to the yubikey
	pemBytes, err := trustmanager.EncryptPrivateKey(key, "passphrase")
	assert.NoError(t, err)
	err = newStore.ImportKey(pemBytes, "root")
	assert.NoError(t, err)

	// key is not in backup store
	_, _, err = backup.GetKey(key.ID())
	assert.Error(t, err)
}

// Importing a key not as root fails, and it is not added to the backup store
func TestYubiImportNonRootKey(t *testing.T) {
	if !YubikeyAccessible() {
		t.Skip("Must have Yubikey access.")
	}
	clearAllKeys(t)

	SetYubikeyKeyMode(KeymodeNone)
	defer func() {
		SetYubikeyKeyMode(KeymodeTouch | KeymodePinOnce)
	}()

	backup := trustmanager.NewKeyMemoryStore(ret)
	origStore, err := NewYubiKeyStore(backup, ret)
	assert.NoError(t, err)

	// generate key and import it
	privKey, err := trustmanager.GenerateECDSAKey(rand.Reader)
	assert.NoError(t, err)

	pemBytes, err := trustmanager.EncryptPrivateKey(privKey, "passphrase")
	assert.NoError(t, err)

	err = origStore.ImportKey(pemBytes, privKey.ID())
	assert.Error(t, err)

	// key is not in backup store
	_, _, err = backup.GetKey(privKey.ID())
	assert.Error(t, err)
}

// One cannot export from hardware - it will not export from the backup
func TestYubiExportKeyFails(t *testing.T) {
	if !YubikeyAccessible() {
		t.Skip("Must have Yubikey access.")
	}
	clearAllKeys(t)

	SetYubikeyKeyMode(KeymodeNone)
	defer func() {
		SetYubikeyKeyMode(KeymodeTouch | KeymodePinOnce)
	}()

	store, err := NewYubiKeyStore(trustmanager.NewKeyMemoryStore(ret), ret)
	assert.NoError(t, err)

	key, err := testAddKey(t, store)
	assert.NoError(t, err)

	_, err = store.ExportKey(key.ID())
	assert.Error(t, err)
}

// If there are keys in the backup store but no keys in the Yubikey,
// listing and getting cannot access the keys in the backup store
func TestYubiListAndGetKeysIgnoresBackup(t *testing.T) {
	if !YubikeyAccessible() {
		t.Skip("Must have Yubikey access.")
	}
	clearAllKeys(t)

	SetYubikeyKeyMode(KeymodeNone)
	defer func() {
		SetYubikeyKeyMode(KeymodeTouch | KeymodePinOnce)
	}()

	backup := trustmanager.NewKeyMemoryStore(ret)
	key, err := testAddKey(t, backup)
	assert.NoError(t, err)

	store, err := NewYubiKeyStore(trustmanager.NewKeyMemoryStore(ret), ret)
	assert.Len(t, store.ListKeys(), 0)
	_, _, err = store.GetKey(key.ID())
	assert.Error(t, err)
}

// Get a YubiPrivateKey.  Check that it has the right algorithm, etc, and
// specifically that you cannot get the private bytes out.  Assume we can
// sign something.
func TestYubiKeyAndSign(t *testing.T) {
	if !YubikeyAccessible() {
		t.Skip("Must have Yubikey access.")
	}
	clearAllKeys(t)

	SetYubikeyKeyMode(KeymodeNone)
	defer func() {
		SetYubikeyKeyMode(KeymodeTouch | KeymodePinOnce)
	}()

	store, err := NewYubiKeyStore(trustmanager.NewKeyMemoryStore(ret), ret)
	assert.NoError(t, err)

	ecdsaPrivateKey, err := testAddKey(t, store)
	assert.NoError(t, err)

	yubiPrivateKey, _, err := store.GetKey(ecdsaPrivateKey.ID())
	assert.NoError(t, err)

	assert.Equal(t, data.ECDSAKey, yubiPrivateKey.Algorithm())
	assert.Equal(t, data.ECDSASignature, yubiPrivateKey.SignatureAlgorithm())
	assert.Equal(t, ecdsaPrivateKey.Public(), yubiPrivateKey.Public())
	assert.Nil(t, yubiPrivateKey.Private())

	// The signature should be verified, but the importing the verifiers causes
	// an import cycle.  A bigger refactor needs to be done to fix it.
	msg := []byte("Hello there")
	_, err = yubiPrivateKey.Sign(rand.Reader, msg, nil)
	assert.NoError(t, err)
}

// ----- Negative tests that use stubbed pkcs11 for error injection -----

type pkcs11Stubbable interface {
	setLibLoader(pkcs11LibLoader)
}

var setupErrors = []string{"Initialize", "GetSlotList", "OpenSession"}

// Create a new store, so that we avoid any cache issues, and list keys
func cleanListKeys(t *testing.T) map[string]string {
	cleanStore, err := NewYubiKeyStore(trustmanager.NewKeyMemoryStore(ret), ret)
	assert.NoError(t, err)
	return cleanStore.ListKeys()
}

// If an error occurs during login, which only some functions do, the function
// under test will clean up after itself
func testYubiFunctionCleansUpOnLoginError(t *testing.T, toStub pkcs11Stubbable,
	functionUnderTest func() error) {

	toStub.setLibLoader(func(string) IPKCS11Ctx {
		return NewStubCtx(map[string]bool{"Login": true})
	})

	err := functionUnderTest()
	assert.Error(t, err)
	// a lot of these functions wrap other errors
	assert.Contains(t, err.Error(), trustmanager.ErrAttemptsExceeded{}.Error())

	// Set Up another time, to ensure we weren't left in a bad state
	// by the previous runs
	ctx, session, err := SetupHSMEnv(pkcs11Lib, defaultLoader)
	assert.NoError(t, err)
	cleanup(ctx, session)
}

// If one of the specified pkcs11 functions errors, the function under test
// will clean up after itself
func testYubiFunctionCleansUpOnSpecifiedErrors(t *testing.T,
	toStub pkcs11Stubbable, functionUnderTest func() error,
	dependentFunctions []string, functionShouldError bool) {

	for _, methodName := range dependentFunctions {

		toStub.setLibLoader(func(string) IPKCS11Ctx {
			return NewStubCtx(
				map[string]bool{methodName: true})
		})

		err := functionUnderTest()
		if functionShouldError {
			assert.Error(t, err,
				fmt.Sprintf("Didn't error when %s errored.", methodName))
			// a lot of these functions wrap other errors
			assert.Contains(t, err.Error(), errInjected{methodName}.Error())
		} else {
			assert.NoError(t, err)
		}
	}

	// Set Up another time, to ensure we weren't left in a bad state
	// by the previous runs
	ctx, session, err := SetupHSMEnv(pkcs11Lib, defaultLoader)
	assert.NoError(t, err)
	cleanup(ctx, session)
}

func TestYubiAddKeyCleansUpOnError(t *testing.T) {
	if !YubikeyAccessible() {
		t.Skip("Must have Yubikey access.")
	}
	clearAllKeys(t)

	SetYubikeyKeyMode(KeymodeNone)
	defer func() {
		SetYubikeyKeyMode(KeymodeTouch | KeymodePinOnce)
	}()

	backup := trustmanager.NewKeyMemoryStore(ret)
	origStore, err := NewYubiKeyStore(backup, ret)
	assert.NoError(t, err)

	var _addkey = func() error {
		_, err := testAddKey(t, origStore)
		return err
	}

	testYubiFunctionCleansUpOnLoginError(t, origStore, _addkey)
	// all the PKCS11 functions AddKey depends on that aren't the login/logout
	testYubiFunctionCleansUpOnSpecifiedErrors(t, origStore, _addkey,
		append(
			setupErrors,
			"FindObjectsInit",
			"FindObjects",
			"FindObjectsFinal",
			"CreateObject",
		), true)

	// given that everything should have errored, there should be no keys on
	// the yubikey and no keys in backup
	assert.Len(t, backup.ListKeys(), 0)
	assert.Len(t, cleanListKeys(t), 0)

	// Logout should not cause a function failure - it s a cleanup failure,
	// which shouldn't break anything, and it should clean up after itself.
	// The key should be added to both stores
	testYubiFunctionCleansUpOnSpecifiedErrors(t, origStore, _addkey,
		[]string{"Logout"}, false)

	listedKeys := cleanListKeys(t)
	assert.Len(t, backup.ListKeys(), 1)
	assert.Len(t, listedKeys, 1)

	// Currently, if GetAttributeValue fails, the function succeeds, because if
	// we can't get the attribute value of an object, we don't know what slot
	// it's in, we assume its occupied slot is free (hence this failure will
	// cause the previous key to be overwritten).  This behavior may need to
	// be revisited.
	testYubiFunctionCleansUpOnSpecifiedErrors(t, origStore, _addkey,
		[]string{"GetAttributeValue"}, false)

	newListedKeys := cleanListKeys(t)
	// because the original key got overwritten
	assert.Len(t, backup.ListKeys(), 2)
	assert.Len(t, newListedKeys, 1)
	for k := range newListedKeys {
		_, ok := listedKeys[k]
		assert.False(t, ok)
	}
}

func TestYubiGetKeyCleansUpOnError(t *testing.T) {
	if !YubikeyAccessible() {
		t.Skip("Must have Yubikey access.")
	}
	clearAllKeys(t)

	SetYubikeyKeyMode(KeymodeNone)
	defer func() {
		SetYubikeyKeyMode(KeymodeTouch | KeymodePinOnce)
	}()

	store, err := NewYubiKeyStore(trustmanager.NewKeyMemoryStore(ret), ret)
	assert.NoError(t, err)
	key, err := testAddKey(t, store)
	assert.NoError(t, err)

	var _getkey = func() error {
		_, _, err := store.GetKey(key.ID())
		return err
	}

	// all the PKCS11 functions GetKey depends on
	testYubiFunctionCleansUpOnSpecifiedErrors(t, store, _getkey,
		append(
			setupErrors,
			"FindObjectsInit",
			"FindObjects",
			"FindObjectsFinal",
			"GetAttributeValue",
		), true)
}

func TestYubiImportKeyCleansUpOnError(t *testing.T) {
	if !YubikeyAccessible() {
		t.Skip("Must have Yubikey access.")
	}
	clearAllKeys(t)

	SetYubikeyKeyMode(KeymodeNone)
	defer func() {
		SetYubikeyKeyMode(KeymodeTouch | KeymodePinOnce)
	}()

	origStore, err := NewYubiKeyStore(trustmanager.NewKeyMemoryStore(ret), ret)
	assert.NoError(t, err)

	privKey, err := trustmanager.GenerateECDSAKey(rand.Reader)
	assert.NoError(t, err)

	pemBytes, err := trustmanager.EncryptPrivateKey(privKey, "passphrase")
	assert.NoError(t, err)

	var _importkey = func() error { return origStore.ImportKey(pemBytes, "root") }

	testYubiFunctionCleansUpOnLoginError(t, origStore, _importkey)
	// all the PKCS11 functions ImportKey depends on that aren't the login/logout
	testYubiFunctionCleansUpOnSpecifiedErrors(t, origStore, _importkey,
		append(
			setupErrors,
			"FindObjectsInit",
			"FindObjects",
			"FindObjectsFinal",
			"CreateObject",
		), true)

	// given that everything should have errored, there should be no keys on
	// the yubikey
	assert.Len(t, cleanListKeys(t), 0)

	// Logout should not cause a function failure - it s a cleanup failure,
	// which shouldn't break anything, and it should clean up after itself.
	// The key should be added to both stores
	testYubiFunctionCleansUpOnSpecifiedErrors(t, origStore, _importkey,
		[]string{"Logout"}, false)

	listedKeys := cleanListKeys(t)
	assert.Len(t, listedKeys, 1)

	// Currently, if GetAttributeValue fails, the function succeeds, because if
	// we can't get the attribute value of an object, we don't know what slot
	// it's in, we assume its occupied slot is free (hence this failure will
	// cause the previous key to be overwritten).  This behavior may need to
	// be revisited.
	for k := range listedKeys {
		err := origStore.RemoveKey(k)
		assert.NoError(t, err)
	}
	testYubiFunctionCleansUpOnSpecifiedErrors(t, origStore, _importkey,
		[]string{"GetAttributeValue"}, false)

	assert.Len(t, cleanListKeys(t), 1)
}

func TestYubiRemoveKeyCleansUpOnError(t *testing.T) {
	if !YubikeyAccessible() {
		t.Skip("Must have Yubikey access.")
	}
	clearAllKeys(t)

	SetYubikeyKeyMode(KeymodeNone)
	defer func() {
		SetYubikeyKeyMode(KeymodeTouch | KeymodePinOnce)
	}()

	origStore, err := NewYubiKeyStore(trustmanager.NewKeyMemoryStore(ret), ret)
	assert.NoError(t, err)
	key, err := testAddKey(t, origStore)
	assert.NoError(t, err)

	var _removekey = func() error { return origStore.RemoveKey(key.ID()) }

	testYubiFunctionCleansUpOnLoginError(t, origStore, _removekey)
	// RemoveKey just succeeds if we can't set up the yubikey
	testYubiFunctionCleansUpOnSpecifiedErrors(t, origStore, _removekey, setupErrors, false)
	// all the PKCS11 functions RemoveKey depends on that aren't the login/logout
	// or setup/cleanup
	testYubiFunctionCleansUpOnSpecifiedErrors(t, origStore, _removekey,
		[]string{
			"FindObjectsInit",
			"FindObjects",
			"FindObjectsFinal",
			"DestroyObject",
		}, true)

	// given that everything should have errored, there should still be 1 key
	// on the yubikey
	assert.Len(t, cleanListKeys(t), 1)

	// this will not fail, but it should clean up after itself, and the key
	// should be added to both stores
	testYubiFunctionCleansUpOnSpecifiedErrors(t, origStore, _removekey,
		[]string{"Logout"}, false)

	assert.Len(t, cleanListKeys(t), 0)
}

func TestYubiListKeyCleansUpOnError(t *testing.T) {
	if !YubikeyAccessible() {
		t.Skip("Must have Yubikey access.")
	}
	clearAllKeys(t)

	SetYubikeyKeyMode(KeymodeNone)
	defer func() {
		SetYubikeyKeyMode(KeymodeTouch | KeymodePinOnce)
	}()

	// Do not call NewYubiKeyStore, because it list keys immediately to
	// build the cache.
	store := &YubiKeyStore{
		passRetriever: ret,
		keys:          make(map[string]yubiSlot),
		backupStore:   trustmanager.NewKeyMemoryStore(ret),
		libLoader:     defaultLoader,
	}

	var _listkeys = func() error {
		// ListKeys never fails
		store.ListKeys()
		return nil
	}

	// all the PKCS11 functions ListKey depends on - list keys never errors
	testYubiFunctionCleansUpOnSpecifiedErrors(t, store, _listkeys,
		append(
			setupErrors,
			"FindObjectsInit",
			"FindObjects",
			"FindObjectsFinal",
			"GetAttributeValue",
		), false)
}

// export key fails anyway, don't bother testing

func TestYubiSignCleansUpOnError(t *testing.T) {
	if !YubikeyAccessible() {
		t.Skip("Must have Yubikey access.")
	}
	clearAllKeys(t)

	SetYubikeyKeyMode(KeymodeNone)
	defer func() {
		SetYubikeyKeyMode(KeymodeTouch | KeymodePinOnce)
	}()

	store, err := NewYubiKeyStore(trustmanager.NewKeyMemoryStore(ret), ret)
	assert.NoError(t, err)

	key, err := testAddKey(t, store)
	assert.NoError(t, err)

	privKey, _, err := store.GetKey(key.ID())
	assert.NoError(t, err)

	yubiPrivateKey, ok := privKey.(*YubiPrivateKey)
	assert.True(t, ok)

	var _sign = func() error {
		_, err = yubiPrivateKey.Sign(rand.Reader, []byte("Hello there"), nil)
		return err
	}

	testYubiFunctionCleansUpOnLoginError(t, yubiPrivateKey, _sign)
	// all the PKCS11 functions SignKey depends on that is not login/logout
	testYubiFunctionCleansUpOnSpecifiedErrors(t, yubiPrivateKey, _sign,
		append(
			setupErrors,
			"FindObjectsInit",
			"FindObjects",
			"FindObjectsFinal",
			"SignInit",
			"Sign",
		), true)
	// this will not fail, but it should clean up after itself, and the key
	// should be added to both stores
	testYubiFunctionCleansUpOnSpecifiedErrors(t, yubiPrivateKey, _sign,
		[]string{"Logout"}, false)
}

// -----  Stubbed pkcs11 for testing error conditions ------
// This is just a passthrough to the underlying pkcs11 library, with optional
// error injection.  This is to ensure that if errors occur during the process
// of interacting with the Yubikey, that everything gets cleaned up sanely.

// Note that this does not actually replicate an actual PKCS11 failure, since
// who knows what the pkcs11 function call may have done to the key before it
// errored. This just tests that we handle an error ok.

type errInjected struct {
	methodName string
}

func (e errInjected) Error() string {
	return fmt.Sprintf("Injected failure in %s", e.methodName)
}

const (
	uninitialized = 0
	initialized   = 1
	sessioned     = 2
	loggedin      = 3
)

type StubCtx struct {
	ctx                IPKCS11Ctx
	functionShouldFail map[string]bool
}

func NewStubCtx(functionShouldFail map[string]bool) *StubCtx {
	realCtx := defaultLoader(pkcs11Lib)
	return &StubCtx{
		ctx:                realCtx,
		functionShouldFail: functionShouldFail,
	}
}

// Returns an error if we're supposed to error for this method
func (s *StubCtx) checkErr(methodName string) error {
	if val, ok := s.functionShouldFail[methodName]; ok && val {
		return errInjected{methodName: methodName}
	}
	return nil
}

func (s *StubCtx) Destroy() {
	// can't error
	s.ctx.Destroy()
}

func (s *StubCtx) Initialize() error {
	err := s.checkErr("Initialize")
	if err != nil {
		return err
	}
	return s.ctx.Initialize()
}

func (s *StubCtx) Finalize() error {
	err := s.checkErr("Finalize")
	if err != nil {
		return err
	}
	return s.ctx.Finalize()
}

func (s *StubCtx) GetSlotList(tokenPresent bool) ([]uint, error) {
	err := s.checkErr("GetSlotList")
	if err != nil {
		return nil, err
	}
	return s.ctx.GetSlotList(tokenPresent)
}

func (s *StubCtx) OpenSession(slotID uint, flags uint) (pkcs11.SessionHandle, error) {
	err := s.checkErr("OpenSession")
	if err != nil {
		return pkcs11.SessionHandle(0), err
	}
	return s.ctx.OpenSession(slotID, flags)
}

func (s *StubCtx) CloseSession(sh pkcs11.SessionHandle) error {
	err := s.checkErr("CloseSession")
	if err != nil {
		return err
	}
	return s.ctx.CloseSession(sh)
}

func (s *StubCtx) Login(sh pkcs11.SessionHandle, userType uint, pin string) error {
	err := s.checkErr("Login")
	if err != nil {
		return err
	}
	return s.ctx.Login(sh, userType, pin)
}

func (s *StubCtx) Logout(sh pkcs11.SessionHandle) error {
	err := s.checkErr("Logout")
	if err != nil {
		return err
	}
	return s.ctx.Logout(sh)
}

func (s *StubCtx) CreateObject(sh pkcs11.SessionHandle, temp []*pkcs11.Attribute) (
	pkcs11.ObjectHandle, error) {
	err := s.checkErr("CreateObject")
	if err != nil {
		return pkcs11.ObjectHandle(0), err
	}
	return s.ctx.CreateObject(sh, temp)
}

func (s *StubCtx) DestroyObject(sh pkcs11.SessionHandle, oh pkcs11.ObjectHandle) error {
	err := s.checkErr("DestroyObject")
	if err != nil {
		return err
	}
	return s.ctx.DestroyObject(sh, oh)
}

func (s *StubCtx) GetAttributeValue(sh pkcs11.SessionHandle, o pkcs11.ObjectHandle,
	a []*pkcs11.Attribute) ([]*pkcs11.Attribute, error) {
	err := s.checkErr("GetAttributeValue")
	if err != nil {
		return nil, err
	}
	return s.ctx.GetAttributeValue(sh, o, a)
}

func (s *StubCtx) FindObjectsInit(sh pkcs11.SessionHandle, temp []*pkcs11.Attribute) error {
	err := s.checkErr("FindObjectsInit")
	if err != nil {
		return err
	}
	return s.ctx.FindObjectsInit(sh, temp)
}

func (s *StubCtx) FindObjects(sh pkcs11.SessionHandle, max int) (
	[]pkcs11.ObjectHandle, bool, error) {
	err := s.checkErr("FindObjects")
	if err != nil {
		return nil, false, err
	}
	return s.ctx.FindObjects(sh, max)
}

func (s *StubCtx) FindObjectsFinal(sh pkcs11.SessionHandle) error {
	err := s.checkErr("FindObjectsFinal")
	if err != nil {
		return err
	}
	return s.ctx.FindObjectsFinal(sh)
}

func (s *StubCtx) SignInit(sh pkcs11.SessionHandle, m []*pkcs11.Mechanism,
	o pkcs11.ObjectHandle) error {
	err := s.checkErr("SignInit")
	if err != nil {
		return err
	}
	return s.ctx.SignInit(sh, m, o)
}

func (s *StubCtx) Sign(sh pkcs11.SessionHandle, message []byte) ([]byte, error) {
	// a call to Sign will clear SignInit whether or not it fails, so
	// replicate that by calling Sign, then optionally returning an error.
	sig, sigErr := s.ctx.Sign(sh, message)
	err := s.checkErr("Sign")
	if err != nil {
		return nil, err
	}
	return sig, sigErr
}
