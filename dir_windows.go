// © 2014 Tajti Attila

package basedir

var (
	Config = newDir("XDG_CONFIG_HOME", "~/.config", "XDG_CONFIG_DIRS", "")
	Data   = newDir("XDG_DATA_HOME", "~/.local/share", "XDG_DATA_DIRS", "")
	Cache  = newDir("XDG_CACHE_HOME", "~/.cache", "", "")
)
