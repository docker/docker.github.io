package snapshot

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/docker/notary/server/storage"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/signed"
)

func TestSnapshotExpired(t *testing.T) {
	sn := &data.SignedSnapshot{
		Signatures: nil,
		Signed: data.Snapshot{
			Expires: time.Now().AddDate(-1, 0, 0),
		},
	}
	assert.True(t, snapshotExpired(sn), "Snapshot should have expired")
}

func TestSnapshotNotExpired(t *testing.T) {
	sn := &data.SignedSnapshot{
		Signatures: nil,
		Signed: data.Snapshot{
			Expires: time.Now().AddDate(1, 0, 0),
		},
	}
	assert.False(t, snapshotExpired(sn), "Snapshot should NOT have expired")
}

func TestGetSnapshotKeyCreate(t *testing.T) {
	store := storage.NewMemStorage()
	crypto := signed.NewEd25519()
	k, err := GetOrCreateSnapshotKey("gun", store, crypto, data.ED25519Key)
	assert.Nil(t, err, "Expected nil error")
	assert.NotNil(t, k, "Key should not be nil")

	k2, err := GetOrCreateSnapshotKey("gun", store, crypto, data.ED25519Key)

	assert.Nil(t, err, "Expected nil error")

	// trying to get the same key again should return the same value
	assert.Equal(t, k, k2, "Did not receive same key when attempting to recreate.")
	assert.NotNil(t, k2, "Key should not be nil")
}

func TestGetSnapshotKeyExisting(t *testing.T) {
	store := storage.NewMemStorage()
	crypto := signed.NewEd25519()
	key, err := crypto.Create(data.CanonicalSnapshotRole, data.ED25519Key)
	assert.NoError(t, err)

	store.SetKey("gun", data.CanonicalSnapshotRole, data.ED25519Key, key.Public())

	k, err := GetOrCreateSnapshotKey("gun", store, crypto, data.ED25519Key)
	assert.Nil(t, err, "Expected nil error")
	assert.NotNil(t, k, "Key should not be nil")
	assert.Equal(t, key, k, "Did not receive same key when attempting to recreate.")
	assert.NotNil(t, k, "Key should not be nil")

	k2, err := GetOrCreateSnapshotKey("gun", store, crypto, data.ED25519Key)

	assert.Nil(t, err, "Expected nil error")

	// trying to get the same key again should return the same value
	assert.Equal(t, k, k2, "Did not receive same key when attempting to recreate.")
	assert.NotNil(t, k2, "Key should not be nil")
}

type keyStore struct {
	getCalled bool
	k         data.PublicKey
}

func (ks *keyStore) GetKey(gun, role string) (string, []byte, error) {
	defer func() { ks.getCalled = true }()
	if ks.getCalled {
		return ks.k.Algorithm(), ks.k.Public(), nil
	}
	return "", nil, &storage.ErrNoKey{}
}

func (ks keyStore) SetKey(gun, role, algorithm string, public []byte) error {
	return &storage.ErrKeyExists{}
}

func TestGetSnapshotKeyExistsOnSet(t *testing.T) {
	crypto := signed.NewEd25519()
	key, err := crypto.Create(data.CanonicalSnapshotRole, data.ED25519Key)
	assert.NoError(t, err)
	store := &keyStore{k: key}

	k, err := GetOrCreateSnapshotKey("gun", store, crypto, data.ED25519Key)
	assert.Nil(t, err, "Expected nil error")
	assert.NotNil(t, k, "Key should not be nil")
	assert.Equal(t, key, k, "Did not receive same key when attempting to recreate.")
	assert.NotNil(t, k, "Key should not be nil")

	k2, err := GetOrCreateSnapshotKey("gun", store, crypto, data.ED25519Key)

	assert.Nil(t, err, "Expected nil error")

	// trying to get the same key again should return the same value
	assert.Equal(t, k, k2, "Did not receive same key when attempting to recreate.")
	assert.NotNil(t, k2, "Key should not be nil")
}

func TestRoleExpired(t *testing.T) {
	meta := data.FileMeta{
		Hashes: data.Hashes{
			"sha256": []byte{1},
		},
	}
	newData := []byte{2}
	res, _ := roleExpired(newData, meta)
	assert.True(t, res)
}

func TestRoleNotExpired(t *testing.T) {
	newData := []byte{2}
	currMeta, err := data.NewFileMeta(bytes.NewReader(newData), "sha256")
	assert.NoError(t, err)

	meta := data.FileMeta{
		Hashes: data.Hashes{
			"sha256": currMeta.Hashes["sha256"],
		},
	}

	res, _ := roleExpired(newData, meta)
	assert.False(t, res)
}

func TestGetSnapshotNotExists(t *testing.T) {
	store := storage.NewMemStorage()
	crypto := signed.NewEd25519()

	_, err := GetOrCreateSnapshot("gun", store, crypto)
	assert.Error(t, err)
}

func TestGetSnapshotCurrValid(t *testing.T) {
	store := storage.NewMemStorage()
	crypto := signed.NewEd25519()

	_, err := GetOrCreateSnapshotKey("gun", store, crypto, data.ED25519Key)

	newData := []byte{2}
	currMeta, err := data.NewFileMeta(bytes.NewReader(newData), "sha256")
	assert.NoError(t, err)

	snapshot := &data.SignedSnapshot{
		Signed: data.Snapshot{
			Expires: data.DefaultExpires(data.CanonicalSnapshotRole),
			Meta: data.Files{
				data.CanonicalRootRole: currMeta,
			},
		},
	}
	snapJSON, _ := json.Marshal(snapshot)

	// test when db is missing the role data
	store.UpdateCurrent("gun", storage.MetaUpdate{Role: "snapshot", Version: 0, Data: snapJSON})
	_, err = GetOrCreateSnapshot("gun", store, crypto)
	assert.NoError(t, err)

	// test when db has the role data
	store.UpdateCurrent("gun", storage.MetaUpdate{Role: "root", Version: 0, Data: newData})
	_, err = GetOrCreateSnapshot("gun", store, crypto)
	assert.NoError(t, err)

	// test when db role data is expired
	store.UpdateCurrent("gun", storage.MetaUpdate{Role: "root", Version: 1, Data: []byte{3}})
	_, err = GetOrCreateSnapshot("gun", store, crypto)
	assert.NoError(t, err)
}

func TestGetSnapshotCurrExpired(t *testing.T) {
	store := storage.NewMemStorage()
	crypto := signed.NewEd25519()

	_, err := GetOrCreateSnapshotKey("gun", store, crypto, data.ED25519Key)

	snapshot := &data.SignedSnapshot{}
	snapJSON, _ := json.Marshal(snapshot)

	store.UpdateCurrent("gun", storage.MetaUpdate{Role: "snapshot", Version: 0, Data: snapJSON})
	_, err = GetOrCreateSnapshot("gun", store, crypto)
	assert.NoError(t, err)
}

func TestGetSnapshotCurrCorrupt(t *testing.T) {
	store := storage.NewMemStorage()
	crypto := signed.NewEd25519()

	_, err := GetOrCreateSnapshotKey("gun", store, crypto, data.ED25519Key)

	snapshot := &data.SignedSnapshot{}
	snapJSON, _ := json.Marshal(snapshot)

	store.UpdateCurrent("gun", storage.MetaUpdate{Role: "snapshot", Version: 0, Data: snapJSON[1:]})
	_, err = GetOrCreateSnapshot("gun", store, crypto)
	assert.Error(t, err)
}
