// © 2014 Tajti Attila

package basedir

var (
	Config = newDir("XDG_CONFIG_HOME", "~/.config", "XDG_CONFIG_DIRS", "/etc/xdg")
	Data   = newDir("XDG_DATA_HOME", "~/Library/Application Support", "XDG_DATA_DIRS", "/usr/local/share/:/usr/share/")
	Cache  = newDir("XDG_CACHE_HOME", "~/Library/Caches", "", "")
)
