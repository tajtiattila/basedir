// © 2014 Tajti Attila

package basedir

import (
	"os"
	"os/user"
)

func addsep(s string) string {
	if s != "" {
		if s[len(s)-1] != os.PathSeparator {
			s += "/"
		}
		return s
	}
	panic("should not happen")
}

func expandTilde(s string) string {
	if len(s) > 2 && s[0] == '~' && os.IsPathSeparator(s[1]) {
		home := os.Getenv("HOME")
		if home == "" {
			if user, err := user.Current(); err != nil {
				home = user.HomeDir
			} else {
				os.Stderr.Write([]byte("basedir: unable to determine home directory"))
				home = "."
			}
		}
		return home + s[1:]
	}
	return s
}
