//nolint
package version

import (
	"fmt"
	"runtime"
	"strings"
)

// Variables set by build flags
var (
	Commit        = ""
	Version       = "" //git tag
	VendorDirHash = ""
	BuildTags     = ""
	ServerName	  = ""
)

type versionInfo struct {
	AppVersion     string `json:"app_version"` //app 版本号
	GitCommit     string `json:"commit"`
	VendorDirHash string `json:"vendor_hash"`
	BuildTags     string `json:"build_tags"`
	GoVersion     string `json:"go"`
}

//func (v versionInfo) String() string {
//	return fmt.Sprintf(`app version: %s
//git commit: %s
//vendor hash: %s
//build tags: %s
//%s`, v.AppVersion, v.GitCommit, v.VendorDirHash, v.BuildTags, v.GoVersion)
//}
func (v versionInfo) String() string {
	return fmt.Sprintf(`app version: %s`, v.AppVersion)
}

func newVersionInfo() versionInfo {
	return versionInfo{
		Version,
		Commit,
		VendorDirHash,
		BuildTags,
		fmt.Sprintf("go version %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)}
}

func GetBigVersion(appVserion string) string {
	return appVserion[0:strings.Index(appVserion,".")]
}