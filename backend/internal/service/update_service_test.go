//go:build unit

package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type updateServiceCacheStub struct {
	data string
}

func (s *updateServiceCacheStub) GetUpdateInfo(context.Context) (string, error) {
	if s.data == "" {
		return "", errors.New("cache miss")
	}
	return s.data, nil
}

func (s *updateServiceCacheStub) SetUpdateInfo(_ context.Context, data string, _ time.Duration) error {
	s.data = data
	return nil
}

type updateServiceGitHubClientStub struct {
	release *GitHubRelease
	calls   int
}

func (s *updateServiceGitHubClientStub) FetchLatestRelease(context.Context, string) (*GitHubRelease, error) {
	s.calls++
	if s.release == nil {
		return nil, errors.New("missing release")
	}
	return s.release, nil
}

func (s *updateServiceGitHubClientStub) DownloadFile(context.Context, string, string, int64) error {
	panic("DownloadFile should not be called when no update is available")
}

func (s *updateServiceGitHubClientStub) FetchChecksumFile(context.Context, string) ([]byte, error) {
	panic("FetchChecksumFile should not be called when no update is available")
}

func TestUpdateServicePerformUpdateNoUpdateReturnsSentinel(t *testing.T) {
	svc := NewUpdateService(
		&updateServiceCacheStub{},
		&updateServiceGitHubClientStub{
			release: &GitHubRelease{
				TagName: "v0.1.132",
				Name:    "v0.1.132",
			},
		},
		"0.1.132",
		"release",
	)

	err := svc.PerformUpdate(context.Background())

	require.Error(t, err)
	require.True(t, errors.Is(err, ErrNoUpdateAvailable))
	require.ErrorIs(t, err, ErrNoUpdateAvailable)
}

func TestUpdateServiceCheckUpdateRenderManagedDeploySkipsReleaseCheck(t *testing.T) {
	t.Setenv(renderEnvName, "true")

	cache := &updateServiceCacheStub{
		data: `{"latest":"9.9.9","timestamp":4102444800}`,
	}
	client := &updateServiceGitHubClientStub{
		release: &GitHubRelease{
			TagName: "v9.9.9",
			Name:    "v9.9.9",
		},
	}
	svc := NewUpdateService(cache, client, "0.1.133", "release")

	info, err := svc.CheckUpdate(context.Background(), false)

	require.NoError(t, err)
	require.Equal(t, "0.1.133", info.CurrentVersion)
	require.Equal(t, "0.1.133", info.LatestVersion)
	require.False(t, info.HasUpdate)
	require.Equal(t, "release", info.BuildType)
	require.Contains(t, info.Warning, "Render auto deploy")
	require.Zero(t, client.calls)
}
