// © 2014 Tajti Attila

package basedir

import (
	"os"
	"path/filepath"
	"syscall"
)

const FileModeDir os.FileMode = 0700 // Permission mode to create missing base directories

// Dir represents a base directory
type Dir struct {
	dirs []string
}

// Gopath is a Dir for accessing files within GOROOT and GOPATH.
var Gopath = newDir("GOROOT", "", "GOPATH", "")

// Open the file subpath inside a directory. The search is started in
// the most important directory, which is specified by the HOME variable.
// If no file is found, the error for the first directory is returned.
func (d *Dir) Open(subpath string) (f *os.File, err error) {
	for _, path := range d.dirs {
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
func (d *Dir) OpenAll(subpath string) (files []*os.File, err error) {
	for _, path := range d.dirs {
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

// Create a new file in the "home" directory.
// If the "home" directory does not exist, it is created with FileModeDir.
func (d *Dir) Create(subpath string) (*os.File, error) {
	path := filepath.Join(d.dirs[0], subpath)
	if err := os.MkdirAll(filepath.Dir(path), FileModeDir); err != nil {
		return nil, err
	}
	return os.Create(path)
}

// Mkdir creates a new directory inside the "home" directory.
// If the "home" directory does not exist, it is created with FileModeDir.
// The subdirectory will be created using the permission bits.
func (d *Dir) Mkdir(subpath string, perm os.FileMode) error {
	if err := os.MkdirAll(d.dirs[0], FileModeDir); err != nil {
		return err
	}
	return os.Mkdir(filepath.Join(d.dirs[0], subpath), perm)
}

// MkdirAll creates a new directory, along with any necessary parents inside the "home" directory.
// If the "home" directory does not exist, it is created with FileModeDir.
// The subdirectories will be created using the permission bits.
func (d *Dir) MkdirAll(subpath string, perm os.FileMode) error {
	if err := os.MkdirAll(d.dirs[0], FileModeDir); err != nil {
		return err
	}
	return os.MkdirAll(filepath.Join(d.dirs[0], subpath), perm)
}

// Remove removes the named file or directory inside the "home" directory.
func (d *Dir) Remove(subpath string) error {
	return os.Remove(filepath.Join(d.dirs[0], subpath))
}

// RemoveAll removes the named file or directory and any children it contains inside the "home" directory.
func (d *Dir) RemoveAll(subpath string) error {
	return os.RemoveAll(filepath.Join(d.dirs[0], subpath))
}

// Dir returns the absolute path if subpath exists in any of the base directories.
// The absolute path returned always contain a path separator at the end if err is nil.
func (d *Dir) Dir(subpath string) (path string, err error) {
	for _, path := range d.dirs {
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
func (d *Dir) EnsureDir(subpath string, perm os.FileMode) (path string, err error) {
	path, err = d.Dir(subpath)
	if err != nil {
		err = d.MkdirAll(subpath, perm)
		if err == nil {
			path = addsep(filepath.Join(d.dirs[0], subpath))
		}
	}
	return
}

// envHome: environment variable name for base directory
// defHome: environment variable name for list of additional directories list
// envDirs: default base directory inside $HOME, if not provided in environment variable
// defDirs: default optional
func newDir(envHome, defHome, envDirs string, defDirs ...string) *Dir {
	base := os.Getenv(envHome)
	if base == "" {
		base = expandTilde(defHome)
	}
	d := new(Dir)
	d.add(base)
	if dirs := os.Getenv(envDirs); dirs != "" {
		s := 0
		for p, ch := range dirs {
			if ch == os.PathListSeparator {
				d.add(dirs[s:p])
				s = p + 1
			}
		}
		d.add(dirs[s:])
	} else {
		d.dirs = append(d.dirs, defDirs...)
	}
	if len(d.dirs) == 0 {
		p, err := os.Getwd()
		if err != nil {
			p = "."
		}
		d.add(p)
	}
	return d
}

func (d *Dir) add(path string) {
	if path != "" {
		d.dirs = append(d.dirs, path)
	}
}
