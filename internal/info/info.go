package info

import (
	"fmt"
	"time"

	"github.com/m-oons/mike/internal/windows/shell32"
)

const (
	AppName    = "Mike"
	Author     = "Moons"
	Repository = "m-oons/mike"
)

var (
	Version = "dev"
	Commit  = "-"
	Date    = "-"
)

func OpenRepository() {
	repo := fmt.Sprintf("https://github.com/%s", Repository)
	shell32.OpenURL(repo)
}

func VersionString() string {
	if Version == "dev" {
		return "dev"
	}

	return fmt.Sprintf("%s (%s)", Version, Commit)
}

func DateString() string {
	if Version == "dev" {
		return time.Now().Format("2006-01-02T15:04:05Z")
	}

	return Date
}
