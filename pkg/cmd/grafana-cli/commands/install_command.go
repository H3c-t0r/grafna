package commands

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/grafana/grafana/pkg/cmd/grafana-cli/utils"
	"github.com/grafana/grafana/pkg/util/errutil"
	"golang.org/x/xerrors"

	"github.com/grafana/grafana/pkg/cmd/grafana-cli/logger"
	"github.com/grafana/grafana/pkg/cmd/grafana-cli/models"
	"github.com/grafana/grafana/pkg/cmd/grafana-cli/services"
	"github.com/grafana/grafana/pkg/setting"
)

func validateInput(c utils.CommandLine, pluginsDir string) error {
	arg := c.Args().First()
	if arg == "" {
		return errors.New("please specify plugin to install")
	}

	if pluginsDir == "" {
		return errors.New("missing pluginsDir flag")
	}

	fileInfo, err := os.Stat(pluginsDir)
	if err != nil {
		if err = os.MkdirAll(pluginsDir, os.ModePerm); err != nil {
			return fmt.Errorf("pluginsDir (%s) is not a writable directory", pluginsDir)
		}
		return nil
	}

	if !fileInfo.IsDir() {
		return errors.New("path is not a directory")
	}

	return nil
}

func (cmd Command) installCommand(c utils.CommandLine) error {
	pluginsDir := setting.PluginsPath
	if err := validateInput(c, pluginsDir); err != nil {
		return err
	}

	pluginToInstall := c.Args().First()
	version := c.Args().Get(1)

	return InstallPlugin(pluginToInstall, version, c, cmd.Client)
}

// InstallPlugin downloads the plugin code as a zip file from the Grafana.com API
// and then extracts the zip into the plugins directory.
func InstallPlugin(pluginName, version string, c utils.CommandLine, client utils.ApiClient) error {
	pluginsDir := setting.PluginsPath
	downloadURL := c.PluginURL()
	isInternal := false

	var checksum string
	if downloadURL == "" {
		if strings.HasPrefix(pluginName, "grafana-") {
			// At this point the plugin download is going through grafana.com API and thus the name is validated.
			// Checking for grafana prefix is how it is done there so no 3rd party plugin should have that prefix.
			// You can supply custom plugin name and then set custom download url to 3rd party plugin but then that
			// is up to the user to know what she is doing.
			isInternal = true
		}
		plugin, err := client.GetPlugin(pluginName, c.RepoDirectory())
		if err != nil {
			return err
		}

		v, err := SelectVersion(&plugin, version)
		if err != nil {
			return err
		}

		if version == "" {
			version = v.Version
		}
		downloadURL = fmt.Sprintf("%s/%s/versions/%s/download",
			c.String("repo"),
			pluginName,
			version,
		)

		// Plugins which are downloaded just as sourcecode zipball from github do not have checksum
		if v.Arch != nil {
			checksum = v.Arch[osAndArchString()].Md5
		}
	}

	logger.Infof("installing %v @ %v\n", pluginName, version)
	logger.Infof("from: %v\n", downloadURL)
	logger.Infof("into: %v\n", pluginsDir)
	logger.Info("\n")

	// Create temp file for downloading zip file
	tmpFile, err := ioutil.TempFile("", "*.zip")
	if err != nil {
		return errutil.Wrap("failed to create temporary file", err)
	}
	defer os.Remove(tmpFile.Name())

	err = client.DownloadFile(pluginName, tmpFile, downloadURL, checksum)
	if err != nil {
		tmpFile.Close()
		return errutil.Wrap("failed to download plugin archive", err)
	}
	err = tmpFile.Close()
	if err != nil {
		return errutil.Wrap("failed to close tmp file", err)
	}

	err = extractFiles(tmpFile.Name(), pluginName, pluginsDir, isInternal)
	if err != nil {
		return errutil.Wrap("failed to extract plugin archive", err)
	}

	logger.Infof("%s Installed %s successfully \n", color.GreenString("✔"), pluginName)

	res, _ := services.ReadPlugin(pluginsDir, pluginName)
	for _, v := range res.Dependencies.Plugins {
		if err := InstallPlugin(v.Id, "", c, client); err != nil {
			return errutil.Wrapf(err, "failed to install plugin '%s'", v.Id)
		}

		logger.Infof("Installed dependency: %v ✔\n", v.Id)
	}

	return err
}

func osAndArchString() string {
	osString := strings.ToLower(runtime.GOOS)
	arch := runtime.GOARCH
	return osString + "-" + arch
}

func supportsCurrentArch(version *models.Version) bool {
	if version.Arch == nil {
		return true
	}
	for arch := range version.Arch {
		if arch == osAndArchString() || arch == "any" {
			return true
		}
	}
	return false
}

func latestSupportedVersion(plugin *models.Plugin) *models.Version {
	for _, ver := range plugin.Versions {
		if supportsCurrentArch(&ver) {
			return &ver
		}
	}
	return nil
}

// SelectVersion returns latest version if none is specified or the specified version. If the version string is not
// matched to existing version it errors out. It also errors out if version that is matched is not available for current
// os and platform. It expects plugin.Versions to be sorted so the newest version is first.
func SelectVersion(plugin *models.Plugin, version string) (*models.Version, error) {
	var ver models.Version

	latestForArch := latestSupportedVersion(plugin)
	if latestForArch == nil {
		return nil, xerrors.New("plugin is not supported on your architecture and OS.")
	}

	if version == "" {
		return latestForArch, nil
	}
	for _, v := range plugin.Versions {
		if v.Version == version {
			ver = v
			break
		}
	}

	if len(ver.Version) == 0 {
		return nil, xerrors.New("could not find the version you're looking for")
	}

	if !supportsCurrentArch(&ver) {
		return nil, xerrors.Errorf("the version you want is not supported on your architecture and OS. Latest suitable version is %s", latestForArch.Version)
	}

	return &ver, nil
}

func RemoveGitBuildFromName(pluginName, filename string) string {
	r := regexp.MustCompile("^[a-zA-Z0-9_.-]*/")
	return r.ReplaceAllString(filename, pluginName+"/")
}

var permissionsDeniedMessage = "Could not create %s. Permission denied. Make sure you have write access to plugindir"

func extractFiles(archiveFile string, pluginName string, filePath string, allowSymlinks bool) error {
	logger.Debugf("Extracting archive %v to %v...\n", archiveFile, filePath)

	r, err := zip.OpenReader(archiveFile)
	if err != nil {
		return err
	}
	for _, zf := range r.File {
		newFileName := RemoveGitBuildFromName(pluginName, zf.Name)
		if !isPathSafe(newFileName, filepath.Join(filePath, pluginName)) {
			return xerrors.Errorf("filepath: %q tries to write outside of plugin directory: %q. This can be a security risk.", zf.Name, path.Join(filePath, pluginName))
		}
		newFile := path.Join(filePath, newFileName)

		if zf.FileInfo().IsDir() {
			err := os.Mkdir(newFile, 0755)
			if os.IsPermission(err) {
				return fmt.Errorf(permissionsDeniedMessage, newFile)
			}
		} else {
			if isSymlink(zf) {
				if !allowSymlinks {
					logger.Errorf("%v: plugin archive contains symlink which is not allowed. Skipping \n", zf.Name)
					continue
				}
				err = extractSymlink(zf, newFile)
				if err != nil {
					logger.Errorf("Failed to extract symlink: %v \n", err)
					continue
				}
			} else {
				err = extractFile(zf, newFile)
				if err != nil {
					return errutil.Wrap("failed to extract file", err)
				}
			}
		}
	}

	return nil
}

func isSymlink(file *zip.File) bool {
	return file.Mode()&os.ModeSymlink == os.ModeSymlink
}

func extractSymlink(file *zip.File, filePath string) error {
	// symlink target is the contents of the file
	src, err := file.Open()
	if err != nil {
		return errutil.Wrap("failed to extract file", err)
	}
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, src)
	if err != nil {
		return errutil.Wrap("failed to copy symlink contents", err)
	}
	err = os.Symlink(strings.TrimSpace(buf.String()), filePath)
	if err != nil {
		return errutil.Wrapf(err, "failed to make symbolic link for %v", filePath)
	}
	return nil
}

func extractFile(file *zip.File, filePath string) (err error) {
	fileMode := file.Mode()
	// This is entry point for backend plugins so we want to make them executable
	if strings.HasSuffix(filePath, "_linux_amd64") || strings.HasSuffix(filePath, "_darwin_amd64") {
		fileMode = os.FileMode(0755)
	}

	dst, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fileMode)
	if err != nil {
		if os.IsPermission(err) {
			return xerrors.Errorf(permissionsDeniedMessage, filePath)
		}

		unwrappedError := xerrors.Unwrap(err)
		if unwrappedError != nil && strings.EqualFold(unwrappedError.Error(), "text file busy") {
			return fmt.Errorf("file %s is in use. Please stop Grafana, install the plugin and restart Grafana", filePath)
		}

		return errutil.Wrap("failed to open file", err)
	}
	defer func() {
		err = dst.Close()
	}()

	src, err := file.Open()
	if err != nil {
		return errutil.Wrap("failed to extract file", err)
	}
	defer func() {
		err = src.Close()
	}()

	_, err = io.Copy(dst, src)
	return
}

// isPathSafe checks if the filePath does not resolve outside of destination. This is used to prevent
// https://snyk.io/research/zip-slip-vulnerability
// Based on https://github.com/mholt/archiver/pull/65/files#diff-635e4219ee55ef011b2b32bba065606bR109
func isPathSafe(filePath string, destination string) bool {
	destpath := filepath.Join(destination, filePath)
	return strings.HasPrefix(destpath, destination)
}
