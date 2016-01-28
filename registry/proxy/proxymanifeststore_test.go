package proxy

import (
	"io"
	"testing"

	"github.com/docker/distribution"
	"github.com/docker/distribution/context"
	"github.com/docker/distribution/digest"
	"github.com/docker/distribution/manifest"
	"github.com/docker/distribution/manifest/schema1"
	"github.com/docker/distribution/reference"
	"github.com/docker/distribution/registry/proxy/scheduler"
	"github.com/docker/distribution/registry/storage"
	"github.com/docker/distribution/registry/storage/cache/memory"
	"github.com/docker/distribution/registry/storage/driver/inmemory"
	"github.com/docker/distribution/testutil"
	"github.com/docker/libtrust"
)

type statsManifest struct {
	manifests distribution.ManifestService
	stats     map[string]int
}

type manifestStoreTestEnv struct {
	manifestDigest digest.Digest // digest of the signed manifest in the local storage
	manifests      proxyManifestStore
}

func (te manifestStoreTestEnv) LocalStats() *map[string]int {
	ls := te.manifests.localManifests.(statsManifest).stats
	return &ls
}

func (te manifestStoreTestEnv) RemoteStats() *map[string]int {
	rs := te.manifests.remoteManifests.(statsManifest).stats
	return &rs
}

func (sm statsManifest) Delete(ctx context.Context, dgst digest.Digest) error {
	sm.stats["delete"]++
	return sm.manifests.Delete(ctx, dgst)
}

func (sm statsManifest) Exists(ctx context.Context, dgst digest.Digest) (bool, error) {
	sm.stats["exists"]++
	return sm.manifests.Exists(ctx, dgst)
}

func (sm statsManifest) Get(ctx context.Context, dgst digest.Digest, options ...distribution.ManifestServiceOption) (distribution.Manifest, error) {
	sm.stats["get"]++
	return sm.manifests.Get(ctx, dgst)
}

func (sm statsManifest) Put(ctx context.Context, manifest distribution.Manifest, options ...distribution.ManifestServiceOption) (digest.Digest, error) {
	sm.stats["put"]++
	return sm.manifests.Put(ctx, manifest)
}

/*func (sm statsManifest) Enumerate(ctx context.Context, manifests []distribution.Manifest, last distribution.Manifest) (n int, err error) {
	sm.stats["enumerate"]++
	return sm.manifests.Enumerate(ctx, manifests, last)
}
*/

func newManifestStoreTestEnv(t *testing.T, name, tag string) *manifestStoreTestEnv {
	nameRef, err := reference.ParseNamed(name)
	if err != nil {
		t.Fatalf("unable to parse reference: %s", err)
	}

	ctx := context.Background()
	truthRegistry, err := storage.NewRegistry(ctx, inmemory.New(), storage.BlobDescriptorCacheProvider(memory.NewInMemoryBlobDescriptorCacheProvider()))
	if err != nil {
		t.Fatalf("error creating registry: %v", err)
	}
	truthRepo, err := truthRegistry.Repository(ctx, nameRef)
	if err != nil {
		t.Fatalf("unexpected error getting repo: %v", err)
	}
	tr, err := truthRepo.Manifests(ctx)
	if err != nil {
		t.Fatal(err.Error())
	}
	truthManifests := statsManifest{
		manifests: tr,
		stats:     make(map[string]int),
	}

	manifestDigest, err := populateRepo(t, ctx, truthRepo, name, tag)
	if err != nil {
		t.Fatalf(err.Error())
	}

	localRegistry, err := storage.NewRegistry(ctx, inmemory.New(), storage.BlobDescriptorCacheProvider(memory.NewInMemoryBlobDescriptorCacheProvider()), storage.EnableRedirect, storage.DisableDigestResumption)
	if err != nil {
		t.Fatalf("error creating registry: %v", err)
	}
	localRepo, err := localRegistry.Repository(ctx, nameRef)
	if err != nil {
		t.Fatalf("unexpected error getting repo: %v", err)
	}
	lr, err := localRepo.Manifests(ctx)
	if err != nil {
		t.Fatal(err.Error())
	}

	localManifests := statsManifest{
		manifests: lr,
		stats:     make(map[string]int),
	}

	s := scheduler.New(ctx, inmemory.New(), "/scheduler-state.json")
	return &manifestStoreTestEnv{
		manifestDigest: manifestDigest,
		manifests: proxyManifestStore{
			ctx:             ctx,
			localManifests:  localManifests,
			remoteManifests: truthManifests,
			scheduler:       s,
			repositoryName:  nameRef,
		},
	}
}

func populateRepo(t *testing.T, ctx context.Context, repository distribution.Repository, name, tag string) (digest.Digest, error) {
	m := schema1.Manifest{
		Versioned: manifest.Versioned{
			SchemaVersion: 1,
		},
		Name: name,
		Tag:  tag,
	}

	for i := 0; i < 2; i++ {
		wr, err := repository.Blobs(ctx).Create(ctx)
		if err != nil {
			t.Fatalf("unexpected error creating test upload: %v", err)
		}

		rs, ts, err := testutil.CreateRandomTarFile()
		if err != nil {
			t.Fatalf("unexpected error generating test layer file")
		}
		dgst := digest.Digest(ts)
		if _, err := io.Copy(wr, rs); err != nil {
			t.Fatalf("unexpected error copying to upload: %v", err)
		}

		if _, err := wr.Commit(ctx, distribution.Descriptor{Digest: dgst}); err != nil {
			t.Fatalf("unexpected error finishing upload: %v", err)
		}
	}

	pk, err := libtrust.GenerateECP256PrivateKey()
	if err != nil {
		t.Fatalf("unexpected error generating private key: %v", err)
	}

	sm, err := schema1.Sign(&m, pk)
	if err != nil {
		t.Fatalf("error signing manifest: %v", err)
	}

	ms, err := repository.Manifests(ctx)
	if err != nil {
		t.Fatalf(err.Error())
	}
	dgst, err := ms.Put(ctx, sm)
	if err != nil {
		t.Fatalf("unexpected errors putting manifest: %v", err)
	}

	return dgst, nil
}

// TestProxyManifests contains basic acceptance tests
// for the pull-through behavior
func TestProxyManifests(t *testing.T) {
	name := "foo/bar"
	env := newManifestStoreTestEnv(t, name, "latest")

	localStats := env.LocalStats()
	remoteStats := env.RemoteStats()

	ctx := context.Background()
	// Stat - must check local and remote
	exists, err := env.manifests.Exists(ctx, env.manifestDigest)
	if err != nil {
		t.Fatalf("Error checking existance")
	}
	if !exists {
		t.Errorf("Unexpected non-existant manifest")
	}

	if (*localStats)["exists"] != 1 && (*remoteStats)["exists"] != 1 {
		t.Errorf("Unexpected exists count : \n%v \n%v", localStats, remoteStats)
	}

	// Get - should succeed and pull manifest into local
	_, err = env.manifests.Get(ctx, env.manifestDigest)
	if err != nil {
		t.Fatal(err)
	}

	if (*localStats)["get"] != 1 && (*remoteStats)["get"] != 1 {
		t.Errorf("Unexpected get count")
	}

	if (*localStats)["put"] != 1 {
		t.Errorf("Expected local put")
	}

	// Stat - should only go to local
	exists, err = env.manifests.Exists(ctx, env.manifestDigest)
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Errorf("Unexpected non-existant manifest")
	}

	if (*localStats)["exists"] != 2 && (*remoteStats)["exists"] != 1 {
		t.Errorf("Unexpected exists count")
	}

	// Get - should get from remote, to test freshness
	_, err = env.manifests.Get(ctx, env.manifestDigest)
	if err != nil {
		t.Fatal(err)
	}

	if (*remoteStats)["get"] != 2 && (*remoteStats)["exists"] != 1 && (*localStats)["put"] != 1 {
		t.Errorf("Unexpected get count")
	}
}

func TestProxyTagService(t *testing.T) {

}
