package timestamp

import (
	"bytes"
	"testing"
	"time"

	"github.com/docker/go/canonical/json"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/signed"
	"github.com/docker/notary/tuf/testutils"
	"github.com/stretchr/testify/require"

	"github.com/docker/notary/server/storage"
)

func TestTimestampExpired(t *testing.T) {
	ts := &data.SignedTimestamp{
		Signatures: nil,
		Signed: data.Timestamp{
			SignedCommon: data.SignedCommon{Expires: time.Now().AddDate(-1, 0, 0)},
		},
	}
	require.True(t, timestampExpired(ts), "Timestamp should have expired")
}

func TestTimestampNotExpired(t *testing.T) {
	ts := &data.SignedTimestamp{
		Signatures: nil,
		Signed: data.Timestamp{
			SignedCommon: data.SignedCommon{Expires: time.Now().AddDate(1, 0, 0)},
		},
	}
	require.False(t, timestampExpired(ts), "Timestamp should NOT have expired")
}

func TestGetTimestampKey(t *testing.T) {
	store := storage.NewMemStorage()
	crypto := signed.NewEd25519()
	k, err := GetOrCreateTimestampKey("gun", store, crypto, data.ED25519Key)
	require.Nil(t, err, "Expected nil error")
	require.NotNil(t, k, "Key should not be nil")

	k2, err := GetOrCreateTimestampKey("gun", store, crypto, data.ED25519Key)

	require.Nil(t, err, "Expected nil error")

	// trying to get the same key again should return the same value
	require.Equal(t, k, k2, "Did not receive same key when attempting to recreate.")
	require.NotNil(t, k2, "Key should not be nil")
}

// If there is no previous timestamp or the previous timestamp is corrupt, then
// even if everything else is in place, getting the timestamp fails
func TestGetTimestampNoPreviousTimestamp(t *testing.T) {
	repo, crypto, err := testutils.EmptyRepo("gun")
	require.NoError(t, err)

	rootJSON, err := json.Marshal(repo.Root)
	require.NoError(t, err)
	snapJSON, err := json.Marshal(repo.Snapshot)
	require.NoError(t, err)

	for _, timestampJSON := range [][]byte{nil, []byte("invalid JSON")} {
		store := storage.NewMemStorage()

		// so we know it's not a failure in getting root or snapshot
		require.NoError(t,
			store.UpdateCurrent("gun", storage.MetaUpdate{Role: data.CanonicalRootRole, Version: 0, Data: rootJSON}))
		require.NoError(t,
			store.UpdateCurrent("gun", storage.MetaUpdate{Role: data.CanonicalSnapshotRole, Version: 0, Data: snapJSON}))

		if timestampJSON != nil {
			require.NoError(t,
				store.UpdateCurrent("gun",
					storage.MetaUpdate{Role: data.CanonicalTimestampRole, Version: 0, Data: timestampJSON}))
		}

		// create a key to be used by GetOrCreateTimestamp
		key, err := crypto.Create(data.CanonicalTimestampRole, "gun", data.ECDSAKey)
		require.NoError(t, err)
		require.NoError(t, store.SetKey("gun", data.CanonicalTimestampRole, key.Algorithm(), key.Public()))

		_, _, err = GetOrCreateTimestamp("gun", store, crypto)
		require.Error(t, err, "GetTimestamp should have failed")
		if timestampJSON == nil {
			require.IsType(t, storage.ErrNotFound{}, err)
		} else {
			require.IsType(t, &json.SyntaxError{}, err)
		}
	}
}

// If there WAS a pre-existing timestamp, and it is not expired, then just return it (it doesn't
// load any other metadata that it doesn't need, like root)
func TestGetTimestampReturnsPreviousTimestampIfUnexpired(t *testing.T) {
	store := storage.NewMemStorage()
	repo, crypto, err := testutils.EmptyRepo("gun")
	require.NoError(t, err)

	snapJSON, err := json.Marshal(repo.Snapshot)
	require.NoError(t, err)
	timestampJSON, err := json.Marshal(repo.Timestamp)
	require.NoError(t, err)

	require.NoError(t, store.UpdateCurrent("gun",
		storage.MetaUpdate{Role: data.CanonicalSnapshotRole, Version: 0, Data: snapJSON}))
	require.NoError(t, store.UpdateCurrent("gun",
		storage.MetaUpdate{Role: data.CanonicalTimestampRole, Version: 0, Data: timestampJSON}))

	_, gottenTimestamp, err := GetOrCreateTimestamp("gun", store, crypto)
	require.NoError(t, err, "GetTimestamp should not have failed")
	require.True(t, bytes.Equal(timestampJSON, gottenTimestamp))
}

func TestGetTimestampOldTimestampExpired(t *testing.T) {
	store := storage.NewMemStorage()
	repo, crypto, err := testutils.EmptyRepo("gun")
	require.NoError(t, err)

	rootJSON, err := json.Marshal(repo.Root)
	require.NoError(t, err)
	snapJSON, err := json.Marshal(repo.Snapshot)
	require.NoError(t, err)

	// create an expired timestamp
	_, err = repo.SignTimestamp(time.Now().AddDate(-1, -1, -1))
	require.True(t, repo.Timestamp.Signed.Expires.Before(time.Now()))
	require.NoError(t, err)
	timestampJSON, err := json.Marshal(repo.Timestamp)
	require.NoError(t, err)

	// set all the metadata
	require.NoError(t, store.UpdateCurrent("gun",
		storage.MetaUpdate{Role: data.CanonicalRootRole, Version: 0, Data: rootJSON}))
	require.NoError(t, store.UpdateCurrent("gun",
		storage.MetaUpdate{Role: data.CanonicalSnapshotRole, Version: 0, Data: snapJSON}))
	require.NoError(t, store.UpdateCurrent("gun",
		storage.MetaUpdate{Role: data.CanonicalTimestampRole, Version: 1, Data: timestampJSON}))

	_, gottenTimestamp, err := GetOrCreateTimestamp("gun", store, crypto)
	require.NoError(t, err, "GetTimestamp errored")

	require.False(t, bytes.Equal(timestampJSON, gottenTimestamp),
		"Timestamp was not regenerated when old one was expired")

	signedMeta := &data.SignedMeta{}
	require.NoError(t, json.Unmarshal(gottenTimestamp, signedMeta))
	// the new metadata is not expired
	require.True(t, signedMeta.Signed.Expires.After(time.Now()))
}

// If the root or snapshot is missing or corrupt, no timestamp can be generated
func TestCannotMakeNewTimestampIfNoRootOrSnapshot(t *testing.T) {
	repo, crypto, err := testutils.EmptyRepo("gun")
	require.NoError(t, err)

	rootJSON, err := json.Marshal(repo.Root)
	require.NoError(t, err)
	snapJSON, err := json.Marshal(repo.Snapshot)
	require.NoError(t, err)

	// create an expired timestamp
	_, err = repo.SignTimestamp(time.Now().AddDate(-1, -1, -1))
	require.True(t, repo.Timestamp.Signed.Expires.Before(time.Now()))
	require.NoError(t, err)
	timestampJSON, err := json.Marshal(repo.Timestamp)
	require.NoError(t, err)

	invalids := []struct {
		test map[string][]byte
		err  error
	}{
		{
			test: map[string][]byte{data.CanonicalRootRole: rootJSON, data.CanonicalSnapshotRole: []byte("invalid JSON")},
			err:  storage.ErrNotFound{},
		},
		{
			test: map[string][]byte{data.CanonicalRootRole: []byte("invalid JSON"), data.CanonicalSnapshotRole: snapJSON},
			err:  &json.SyntaxError{},
		},
		{
			test: map[string][]byte{data.CanonicalRootRole: rootJSON},
			err:  storage.ErrNotFound{},
		},
		{
			test: map[string][]byte{data.CanonicalSnapshotRole: snapJSON},
			err:  storage.ErrNotFound{},
		},
	}

	for _, test := range invalids {
		dataToSet := test.test
		store := storage.NewMemStorage()
		for roleName, jsonBytes := range dataToSet {
			require.NoError(t, store.UpdateCurrent("gun",
				storage.MetaUpdate{Role: roleName, Version: 0, Data: jsonBytes}))
		}
		require.NoError(t, store.UpdateCurrent("gun",
			storage.MetaUpdate{Role: data.CanonicalTimestampRole, Version: 1, Data: timestampJSON}))

		_, _, err := GetOrCreateTimestamp("gun", store, crypto)
		require.Error(t, err, "GetTimestamp errored")
		require.IsType(t, test.err, err)
	}
}

func TestCreateTimestampNoKeyInCrypto(t *testing.T) {
	store := storage.NewMemStorage()
	repo, _, err := testutils.EmptyRepo("gun")
	require.NoError(t, err)

	rootJSON, err := json.Marshal(repo.Root)
	require.NoError(t, err)
	snapJSON, err := json.Marshal(repo.Snapshot)
	require.NoError(t, err)

	// create an expired timestamp
	_, err = repo.SignTimestamp(time.Now().AddDate(-1, -1, -1))
	require.True(t, repo.Timestamp.Signed.Expires.Before(time.Now()))
	require.NoError(t, err)
	timestampJSON, err := json.Marshal(repo.Timestamp)
	require.NoError(t, err)

	// set all the metadata so we know the failure to sign is just because of the key
	require.NoError(t, store.UpdateCurrent("gun",
		storage.MetaUpdate{Role: data.CanonicalRootRole, Version: 0, Data: rootJSON}))
	require.NoError(t, store.UpdateCurrent("gun",
		storage.MetaUpdate{Role: data.CanonicalSnapshotRole, Version: 0, Data: snapJSON}))
	require.NoError(t, store.UpdateCurrent("gun",
		storage.MetaUpdate{Role: data.CanonicalTimestampRole, Version: 1, Data: timestampJSON}))

	// pass it a new cryptoservice without the key
	_, _, err = GetOrCreateTimestamp("gun", store, signed.NewEd25519())
	require.Error(t, err)
	require.IsType(t, signed.ErrInsufficientSignatures{}, err)
}
