// +build !windows,!darwin

package paths

import "os"

var hasVendorName = true
var systemSettingFolders string
var globalSettingFolder string
var cacheFolder string

func init() {
	systemSettingFolders = os.Getenv("HOME")
}
