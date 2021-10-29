package paths

import (
	"os"
	"os/user"
	"runtime"
)

func ProgramDataFolder1() string {
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
	return ProgramDataFolder1()
}

func HomeGazerFolder() string {
	usr, err := user.Current()
	if err == nil {
		return usr.HomeDir + "/gazer"
	}
	return ProgramDataFolder1()
}

func DocumentsFolder() string {
	usr, err := user.Current()
	if err == nil {
		if runtime.GOOS == "windows" {
			return usr.HomeDir + "\\Documents"
		}
		if runtime.GOOS == "darwin" {
			return usr.HomeDir + "/Documents"
		}
		if runtime.GOOS == "linux" {
			return usr.HomeDir + "/Documents"
		}
		return usr.HomeDir
	}
	return ProgramDataFolder1()
}

func DownloadsFolder() string {
	usr, err := user.Current()
	if err == nil {
		if runtime.GOOS == "windows" {
			return usr.HomeDir + "\\Downloads"
		}
		if runtime.GOOS == "darwin" {
			return usr.HomeDir + "/Downloads"
		}
		if runtime.GOOS == "linux" {
			return usr.HomeDir + "/Downloads"
		}
		return usr.HomeDir
	}
	return ProgramDataFolder1()
}

func PicturesFolder() string {
	usr, err := user.Current()
	if err == nil {
		if runtime.GOOS == "windows" {
			return usr.HomeDir + "\\Pictures"
		}
		if runtime.GOOS == "darwin" {
			return usr.HomeDir + "/Pictures"
		}
		if runtime.GOOS == "linux" {
			return usr.HomeDir + "/Pictures"
		}
		return usr.HomeDir
	}
	return ProgramDataFolder1()
}

func DesktopFolder() string {
	usr, err := user.Current()
	if err == nil {
		if runtime.GOOS == "windows" {
			return usr.HomeDir + "\\Desktop"
		}
		if runtime.GOOS == "darwin" {
			return usr.HomeDir + "/Desktop"
		}
		if runtime.GOOS == "linux" {
			return usr.HomeDir + "/Desktop"
		}
		return usr.HomeDir
	}
	return ProgramDataFolder1()
}
