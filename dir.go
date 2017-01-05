// © 2014 Tajti Attila

package basedir

import (
	"os"
	"path/filepath"
	"runtime"
	"syscall"
)

const FileModeDir os.FileMode = 0700 // Permission mode to create missing base directories

var (
	Config Dir // Config stores configuration files.
	Data   Dir // Data stores application data files.
	Cache  Dir // Cache stores non-essential data files.

	Gopath Dir // Gopath is for accessing files within GOROOT and GOPATH.
)

func init() {
	switch runtime.GOOS {
	case "linux":
		Config = newDir(
			"XDG_CONFIG_HOME", "~/.config",
			"XDG_CONFIG_DIRS", "/etc/xdg")
		Data = newDir(
			"XDG_DATA_HOME", "~/.local/share",
			"XDG_DATA_DIRS", "/usr/local/share/", "/usr/share/")
		Cache = newDir("XDG_CACHE_HOME", "~/.cache", "", "")
	case "darwin":
		Config = newDir(
			"XDG_CONFIG_HOME", "~/.config",
			"XDG_CONFIG_DIRS", "/etc/xdg")
		Data = newDir(
			"XDG_DATA_HOME", "~/Library/Application Support",
			"XDG_DATA_DIRS", "/usr/local/share/:/usr/share/")
		Cache = newDir("XDG_CACHE_HOME", "~/Library/Caches", "", "")
	case "windows":
		Config = newDir(
			"XDG_CONFIG_HOME", "~/.config",
			"XDG_CONFIG_DIRS", "")
		Data = newDir(
			"XDG_DATA_HOME", "~/.local/share",
			"XDG_DATA_DIRS", "")
		Cache = newDir("XDG_CACHE_HOME", "~/.cache", "", "")
	}
	Gopath = newDir("GOROOT", "", "GOPATH", "")
}

// Dir represents a base directory.
type Dir []string

// Open the file subpath inside a directory. The search is started in
// the most important directory, which is specified by the HOME variable.
// If no file is found, the error for the first directory is returned.
func (d Dir) Open(subpath string) (f *os.File, err error) {
	for _, path := range d {
		var e error
		if f, e = os.Open(filepath.Join(path, subpath)); e == nil {
			return f, nil
		}
		if err == nil {
			err = e // keep first error
		}
	}
	return
}

// Opens all files matching subpath from all directories.
// If no file is found, the error for the first directory is returned.
func (d Dir) OpenAll(subpath string) (files []*os.File, err error) {
	for _, path := range d {
		var f *os.File
		if f, err = os.Open(filepath.Join(path, subpath)); err == nil {
			files = append(files, f)
		}
	}
	if len(files) != 0 {
		err = nil
	}
	return
}

// Create creates a new file in d under subpath.
// If the parent(s) of subpath do not exist, it is created with FileModeDir.
func (d Dir) Create(subpath string) (*os.File, error) {
	path := filepath.Join(d[0], subpath)
	if err := os.MkdirAll(filepath.Dir(path), FileModeDir); err != nil {
		return nil, err
	}
	return os.Create(path)
}

// Mkdir creates a new directory inside d.
// If the parents of subdir do not exist, they are created with FileModeDir.
// The subdirectory will be created using the permission bits perm.
func (d Dir) Mkdir(subpath string, perm os.FileMode) error {
	if err := os.MkdirAll(d[0], FileModeDir); err != nil {
		return err
	}
	return os.Mkdir(filepath.Join(d[0], subpath), perm)
}

// MkdirAll creates a new directory, along with any necessary parents inside the "home" directory.
// If the "home" directory does not exist, it is created with FileModeDir.
// The subdirectories will be created using the permission bits.
func (d Dir) MkdirAll(subpath string, perm os.FileMode) error {
	if err := os.MkdirAll(d[0], FileModeDir); err != nil {
		return err
	}
	return os.MkdirAll(filepath.Join(d[0], subpath), perm)
}

// Remove removes the named file or directory inside the "home" directory.
func (d Dir) Remove(subpath string) error {
	return os.Remove(filepath.Join(d[0], subpath))
}

// RemoveAll removes the named file or directory and any children it contains inside the "home" directory.
func (d Dir) RemoveAll(subpath string) error {
	return os.RemoveAll(filepath.Join(d[0], subpath))
}

// Dir returns the absolute path if subpath exists in any of the base directories.
// The absolute path returned always contain a path separator at the end if err is nil.
func (d Dir) Dir(subpath string) (path string, err error) {
	for _, path := range d {
		path = filepath.Join(path, subpath)
		fi, err2 := os.Stat(path)
		if err2 == nil {
			if fi.IsDir() {
				return addsep(path), nil
			} else {
				err2 = &os.PathError{"Opendir", path, syscall.ENOTDIR}
			}
		}
		if err == nil {
			err = err2
		}
	}
	return "", err
}

// EnsureDir returns the absolute path if subpath exists in any of the base directories.
// If subpath does not exist in any of the base directories,
// it is created in the "home" directory with FileModeDir.
//
// It is equivalent to calling Dir() and then MkdirAll() if the first call has failed.
func (d Dir) EnsureDir(subpath string, perm os.FileMode) (path string, err error) {
	path, err = d.Dir(subpath)
	if err != nil {
		err = d.MkdirAll(subpath, perm)
		if err == nil {
			path = addsep(filepath.Join(d[0], subpath))
		}
	}
	return
}

// newDir creates a new Dir based on the environment variables
// varHome and varExtraDirs. If one or both environment variable(s)
// are undefined, the arguments home and dirs in their absence.
func newDir(varHome, home, varDirs string, dirs ...string) Dir {
	if envHome := os.Getenv(varHome); envHome != "" {
		home = envHome
	} else {
		home = expandTilde(home)
	}

	if envDirs := filepath.SplitList(os.Getenv(varDirs)); len(envDirs) != 0 {
		dirs = envDirs
	}
	d := append(Dir{home}, dirs...)

	// TODO: check empty slice in public funcs of d instead
	if len(d) == 0 {
		p, err := os.Getwd()
		if err != nil {
			p = "."
		}
		d = append(d, p)
	}
	return d
}
