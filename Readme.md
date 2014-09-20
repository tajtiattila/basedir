
basedir
=======

Package basedir provides functionality to access application-specific
config, data and cache directories for not
completely unlike the XDG base directory specification:
 http://standards.freedesktop.org/basedir-spec/basedir-spec-0.6.html

This package uses os.PathListSeparator instead of as a path separator,
that corresponds the colon ':' specified by XDG on Linux and Max OS X.
It is ';' however on Windows.

The default directories for config, data and cache match that
of XDG on Linux only. On Mac OS X, File System Programming Guide of
the Mac Developer Library takes precedence. On Windows, the HOME
directories are the same as on Linux, but the defaults for DIRS
are not used.

The defaults are:

Linux

   XDG_CONFIG_HOME=~/.config
   XDG_CONFIG_DIRS=/etc/xdg
   XDG_DATA_HOME=~/.local/share
   XDG_DATA_DIRS=/usr/local/share/:/usr/share/
   XDG_CACHE_HOME=~/.cache

Mac OS X

The default for XDG_CONFIG_HOME ~/.config, because ~/Library/Preferences is
only intended for files created indirectly through NSUserDefaults. The
following defaults are different from the defaults for Linux:

   XDG_CACHE_HOME=/Library/Caches
   XDG_CONFIG_HOME=/Library/Application Support

Windows

The defaults for XDG_DATA_DIRS and XDG_CONFIG_DIRS are empty, the other
defaults are the same as on Linux.
