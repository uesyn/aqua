package download

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/aquaproj/aqua/pkg/config"
	"github.com/aquaproj/aqua/pkg/domain"
	"github.com/aquaproj/aqua/pkg/runtime"
	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

var errUnknownChecksumFileType = errors.New("unknown checksum type")

type ChecksumDownloader struct {
	github    domain.RepositoriesService
	runtime   *runtime.Runtime
	http      HTTPDownloader
	ghRelease domain.GitHubReleaseDownloader
}

func NewChecksumDownloader(gh domain.RepositoriesService, rt *runtime.Runtime, httpDownloader HTTPDownloader) *ChecksumDownloader {
	return &ChecksumDownloader{
		github:    gh,
		runtime:   rt,
		http:      httpDownloader,
		ghRelease: NewGitHubReleaseDownloader(gh, httpDownloader),
	}
}

func (dl *ChecksumDownloader) DownloadChecksum(ctx context.Context, logE *logrus.Entry, rt *runtime.Runtime, pkg *config.Package) (io.ReadCloser, int64, error) {
	pkgInfo := pkg.PackageInfo
	switch pkg.PackageInfo.Checksum.Type {
	case config.PkgInfoTypeGitHubRelease:
		asset, err := pkg.RenderChecksumFileName(rt)
		if err != nil {
			return nil, 0, fmt.Errorf("render a checksum file name: %w", err)
		}
		return dl.ghRelease.DownloadGitHubRelease(ctx, logE, &domain.DownloadGitHubReleaseParam{ //nolint:wrapcheck
			RepoOwner: pkgInfo.RepoOwner,
			RepoName:  pkgInfo.RepoName,
			Version:   pkg.Package.Version,
			Asset:     asset,
		})
	case config.PkgInfoTypeHTTP:
		u, err := pkg.RenderChecksumURL(rt)
		if err != nil {
			return nil, 0, fmt.Errorf("render a checksum file name: %w", err)
		}
		rc, code, err := dl.http.Download(ctx, u)
		if err != nil {
			return rc, code, fmt.Errorf("download a checksum file: %w", logerr.WithFields(err, logrus.Fields{
				"download_url": u,
			}))
		}
		return rc, code, nil
	default:
		return nil, 0, logerr.WithFields(errUnknownChecksumFileType, logrus.Fields{ //nolint:wrapcheck
			"package_type": pkgInfo.GetType(),
		})
	}
}
