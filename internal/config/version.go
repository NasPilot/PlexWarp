package config

import (
	"PlexWarp/constants"
	"runtime"
	"time"
)

var (
	appVersion string = "v0.1.1"
	commitHash string = "Unknown"
	buildDate  string = "Unknown"
)

func parseBuildTime(s string) string {
	if t, err := time.Parse(time.RFC3339, s); err != nil {
		return "Unknown"
	} else {
		return t.Local().Format(constants.FORMATE_TIME + " -07:00")
	}
}

// Version 返回版本信息
func Version() VersionInfo {
	return VersionInfo{
		AppVersion: appVersion,
		CommitHash: commitHash,
		BuildData:  parseBuildTime(buildDate),
		GoVersion:  runtime.Version(),
		OS:         runtime.GOOS,
		Arch:       runtime.GOARCH,
	}
}