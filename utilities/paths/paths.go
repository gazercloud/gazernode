package paths

import (
	"os"
	"os/user"
	"runtime"
)

func ProgramDataFolder() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("PROGRAMDATA")
	}
	if runtime.GOOS == "linux" {
		return "/var"
	}
	if runtime.GOOS == "darwin" {
		return "/var"
	}
	return ""
}

func HomeFolder() string {
	usr, err := user.Current()
	if err == nil {
		return usr.HomeDir
	}
	return ProgramDataFolder()
}

func DocumentsFolder() string {
	usr, err := user.Current()
	if err == nil {
		if runtime.GOOS == "windows" {
			return usr.HomeDir + "/Documents"
		}
		if runtime.GOOS == "darwin" {
			return usr.HomeDir + "/Documents"
		}
		if runtime.GOOS == "linux" {
			return usr.HomeDir + "/Documents"
		}
		return usr.HomeDir
	}
	return ProgramDataFolder()
}

func DownloadsFolder() string {
	usr, err := user.Current()
	if err == nil {
		if runtime.GOOS == "windows" {
			return usr.HomeDir + "/Downloads"
		}
		if runtime.GOOS == "darwin" {
			return usr.HomeDir + "/Downloads"
		}
		if runtime.GOOS == "linux" {
			return usr.HomeDir + "/Downloads"
		}
		return usr.HomeDir
	}
	return ProgramDataFolder()
}

func PicturesFolder() string {
	usr, err := user.Current()
	if err == nil {
		if runtime.GOOS == "windows" {
			return usr.HomeDir + "/Pictures"
		}
		if runtime.GOOS == "darwin" {
			return usr.HomeDir + "/Pictures"
		}
		if runtime.GOOS == "linux" {
			return usr.HomeDir + "/Pictures"
		}
		return usr.HomeDir
	}
	return ProgramDataFolder()
}

func DesktopFolder() string {
	usr, err := user.Current()
	if err == nil {
		if runtime.GOOS == "windows" {
			return usr.HomeDir + "/Desktop"
		}
		if runtime.GOOS == "darwin" {
			return usr.HomeDir + "/Desktop"
		}
		if runtime.GOOS == "linux" {
			return usr.HomeDir + "/Desktop"
		}
		return usr.HomeDir
	}
	return ProgramDataFolder()
}
