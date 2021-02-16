package paths

import "os"

var hasVendorName = true
var systemSettingFolders = os.Getenv("PROGRAMDATA")
var globalSettingFolder = os.Getenv("APPDATA")
var cacheFolder = os.Getenv("LOCALAPPDATA")
