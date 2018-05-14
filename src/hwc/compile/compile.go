package compile

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudfoundry/libbuildpack"

	"fmt"
)

type Manifest interface {
	DefaultVersion(string) (libbuildpack.Dependency, error)
	InstallDependency(libbuildpack.Dependency, string) error
	RootDir() string
}

type Compiler struct {
	BuildDir string
	Manifest Manifest
	Log      *libbuildpack.Logger
}

func (c *Compiler) Compile() error {
	err := c.CheckWebConfig()
	if err != nil {
		c.Log.Error("Unable to locate web.config: %s", err.Error())
		return err
	}

	err = c.InstallHWC()
	if err != nil {
		c.Log.Error("Unable to install HWC: %s", err.Error())
		return err
	}

	return nil
}

var (
	errInvalidBuildDir  = errors.New("Invalid build directory provided")
	errMissingWebConfig = errors.New("Missing Web.config")
)

func (c *Compiler) CheckWebConfig() error {
	_, err := os.Stat(c.BuildDir)
	if err != nil {
		return errInvalidBuildDir
	}

	files, err := ioutil.ReadDir(c.BuildDir)
	if err != nil {
		return errInvalidBuildDir
	}

	var webConfigExists bool
	for _, f := range files {
		if strings.ToLower(f.Name()) == "web.config" {
			webConfigExists = true
			break
		}
	}

	if !webConfigExists {
		return errMissingWebConfig
	}
	return nil
}

// Copy directory, skipping files that already exist
func CopyFilesNoOverwrite(srcDir, destDir string) error {
	destExists, _ := libbuildpack.FileExists(destDir)
	if !destExists {
		return errors.New("destination dir must exist")
	}

	files, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, f := range files {
		src := filepath.Join(srcDir, f.Name())
		dest := filepath.Join(destDir, f.Name())

		if f.IsDir() {
			err = os.MkdirAll(dest, f.Mode())
			if err != nil {
				return err
			}
			if err := CopyFilesNoOverwrite(src, dest); err != nil {
				return err
			}
		} else {
			if exists, err := libbuildpack.FileExists(dest); exists {
				if err != nil {return err}
				// fmt.Println("Using file from app directory instead of from buildpack: %s",
				// 	filepath.Base(dest))
				continue
			}
			if err = libbuildpack.CopyFile(src, dest); err != nil {
				return err
			}
		}
	}

	return nil
}

// Hacky solution until PCF supports profile.d and multi buildpack for HWC
func (c *Compiler) InstallAppdynamics() error {
	c.Log.BeginStep("Installing Appdynamics")

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
    	return err
    }

	// Copy files from buildpack's /profile to the server's /.cloudfoundry dir
	// We can run these scripts using the buildpack's start command
	profileSrcDir := filepath.Join(dir, "..", "profile")
	profileDstDir := filepath.Join(c.BuildDir, ".cloudfoundry")

	// Must create source directory first
	if err = os.MkdirAll(profileDstDir, 0755); err != nil { return err }

	err = libbuildpack.CopyDirectory(profileSrcDir, profileDstDir)
	if err != nil {
		fmt.Println("Error copying /profile files")
	} else {
		fmt.Println("Copied files from /profile to /.cloudfoundry")
	}

	// Copy AppDynamics dlls to build directory so we can access them after the app starts because
	// the app will be run in a different container than from where we are now (staging)
	appdynamicsDir := filepath.Join(dir, "..", "appdynamics")
	appdynamicsBuildDir := filepath.Join(c.BuildDir, ".appdynamics")
	if err = os.MkdirAll(appdynamicsBuildDir, 0755); err != nil { return err }

	if err = CopyFilesNoOverwrite(appdynamicsDir, appdynamicsBuildDir); err != nil {
		fmt.Println("Error copying AppDynamics files")
	} else {
		fmt.Println("Copied AppDynamics files to application directory")
	}

	return nil
}

func (c *Compiler) InstallHWC() error {
	if err := c.InstallAppdynamics(); err != nil {
		return err
	}

	c.Log.BeginStep("Installing HWC")
	
	defaultHWC, err := c.Manifest.DefaultVersion("hwc")
	if err != nil {
		return err
	}

	c.Log.Info("HWC version %s", defaultHWC.Version)
	hwcDir := filepath.Join(c.BuildDir, ".cloudfoundry")
	return c.Manifest.InstallDependency(defaultHWC, hwcDir)
}
